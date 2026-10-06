package apps

import (
	"context"
	"os"
	"time"

	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	"github.com/MamangRust/monolith-point-of-sale-pkg/resilience"

	mencache "github.com/MamangRust/monolith-point-of-sale-order/cache"
	"github.com/MamangRust/monolith-point-of-sale-order/handler"
	"github.com/MamangRust/monolith-point-of-sale-order/repository"
	"github.com/MamangRust/monolith-point-of-sale-order/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/server"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
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

	// Cashier, merchant, product and order-item live in other services, so the
	// order flow reaches them over gRPC under dependency guards.
	cashierAddr := getEnv("GRPC_CASHIER_ADDR", "localhost:50055")
	merchantAddr := getEnv("GRPC_MERCHANT_ADDR", "localhost:50056")
	orderItemAddr := getEnv("GRPC_ORDERITEM_ADDR", "localhost:50057")
	productAddr := getEnv("GRPC_PRODUCT_ADDR", "localhost:50059")

	srv.Logger.Info("Connecting to remote gRPC services from Order service",
		zap.String("cashier", cashierAddr),
		zap.String("merchant", merchantAddr),
		zap.String("order_item", orderItemAddr),
		zap.String("product", productAddr),
	)

	cashierConn, err := server.NewGRPCClient(cashierAddr)
	if err != nil {
		return nil, err
	}

	merchantConn, err := server.NewGRPCClient(merchantAddr)
	if err != nil {
		cashierConn.Close()
		return nil, err
	}

	orderItemConn, err := server.NewGRPCClient(orderItemAddr)
	if err != nil {
		cashierConn.Close()
		merchantConn.Close()
		return nil, err
	}

	productConn, err := server.NewGRPCClient(productAddr)
	if err != nil {
		cashierConn.Close()
		merchantConn.Close()
		orderItemConn.Close()
		return nil, err
	}

	go func() {
		<-srv.Ctx.Done()
		srv.Logger.Info("Closing order service remote gRPC connections")
		cashierConn.Close()
		merchantConn.Close()
		orderItemConn.Close()
		productConn.Close()
	}()

	guardCashier := resilience.NewDependencyGuard("cashier", 5, 30, 100, 3*time.Second, srv.Logger)
	guardMerchant := resilience.NewDependencyGuard("merchant", 5, 30, 100, 3*time.Second, srv.Logger)
	guardProduct := resilience.NewDependencyGuard("product", 5, 30, 100, 3*time.Second, srv.Logger)
	guardOrderItem := resilience.NewDependencyGuard("order_item", 5, 30, 100, 3*time.Second, srv.Logger)

	repos := repository.NewRepositories(
		srv.MongoDB,
		pbcashiers.NewCashierQueryServiceClient(cashierConn),
		pbmerchants.NewMerchantQueryServiceClient(merchantConn),
		pbproducts.NewProductQueryServiceClient(productConn),
		pbproducts.NewProductCommandServiceClient(productConn),
		pborder_items.NewOrderItemQueryServiceClient(orderItemConn),
		pborder_items.NewOrderItemCommandServiceClient(orderItemConn),
		repository.GuardOptions{
			Cashier:   []adapter.GuardOption{adapter.WithDependencyGuard(guardCashier)},
			Merchant:  []adapter.GuardOption{adapter.WithDependencyGuard(guardMerchant)},
			Product:   []adapter.GuardOption{adapter.WithDependencyGuard(guardProduct)},
			OrderItem: []adapter.GuardOption{adapter.WithDependencyGuard(guardOrderItem)},
		},
	)

	traceLoggerObservability := observability.NewTraceLoggerObservability(srv.Logger)

	cacheMetrics, err := observability.NewCacheMetrics("order")
	if err != nil {
		return nil, err
	}

	cacheStore := cache.NewCacheStore(srv.Redis, srv.Logger, cacheMetrics)
	mencacheObj := mencache.NewMencache(cacheStore)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Ctx:           context.Background(),
		Repositories:  repos,
		Logger:        srv.Logger,
		Observability: traceLoggerObservability,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pborders.RegisterOrderQueryServiceServer(gs, handlers.Order)
	}

	return srv, nil
}
