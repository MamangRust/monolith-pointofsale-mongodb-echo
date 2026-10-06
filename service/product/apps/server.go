package apps

import (
	"context"
	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"
	"os"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	"github.com/MamangRust/monolith-point-of-sale-pkg/resilience"
	"github.com/MamangRust/monolith-point-of-sale-pkg/server"
	mencache "github.com/MamangRust/monolith-point-of-sale-product/cache"
	"github.com/MamangRust/monolith-point-of-sale-product/handler"
	"github.com/MamangRust/monolith-point-of-sale-product/repository"
	"github.com/MamangRust/monolith-point-of-sale-product/service"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	// A product references a category and a merchant, both owned by other
	// services, so validation and name resolution now go over gRPC.
	categoryAddr := getEnv("GRPC_CATEGORY_ADDR", "localhost:50054")
	merchantAddr := getEnv("GRPC_MERCHANT_ADDR", "localhost:50056")

	srv.Logger.Info("Connecting to remote gRPC services from Product service",
		zap.String("category", categoryAddr), zap.String("merchant", merchantAddr))

	categoryConn, err := server.NewGRPCClient(categoryAddr)
	if err != nil {
		return nil, err
	}

	merchantConn, err := server.NewGRPCClient(merchantAddr)
	if err != nil {
		categoryConn.Close()
		return nil, err
	}

	go func() {
		<-srv.Ctx.Done()
		srv.Logger.Info("Closing product service remote gRPC connections")
		categoryConn.Close()
		merchantConn.Close()
	}()

	guardCategory := resilience.NewDependencyGuard("category", 5, 30, 100, 3*time.Second, srv.Logger)
	guardMerchant := resilience.NewDependencyGuard("merchant", 5, 30, 100, 3*time.Second, srv.Logger)

	repos := repository.NewRepositories(
		srv.MongoDB,
		pbcategories.NewCategoryQueryServiceClient(categoryConn),
		pbmerchants.NewMerchantQueryServiceClient(merchantConn),
		repository.GuardOptions{
			Category: []adapter.GuardOption{adapter.WithDependencyGuard(guardCategory)},
			Merchant: []adapter.GuardOption{adapter.WithDependencyGuard(guardMerchant)},
		},
	)

	mencacheObj := mencache.NewMencache(srv.CacheStore)

	obs := observability.NewTraceLoggerObservability(srv.Logger)

	services := service.NewService(&service.Deps{
		Ctx:           context.Background(),
		Mencache:      mencacheObj,
		Repositories:  repos,
		Logger:        srv.Logger,
		Observability: obs,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbproducts.RegisterProductQueryServiceServer(gs, handlers.Product)
	}

	return srv, nil
}
