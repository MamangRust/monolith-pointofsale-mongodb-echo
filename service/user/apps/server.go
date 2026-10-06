package apps

import (
	"os"
	"time"

	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	pbuserroles "github.com/MamangRust/monolith-point-of-sale-pb/user_roles"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	"github.com/MamangRust/monolith-point-of-sale-pkg/hash"
	"github.com/MamangRust/monolith-point-of-sale-pkg/resilience"
	"github.com/MamangRust/monolith-point-of-sale-pkg/server"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	mencache "github.com/MamangRust/monolith-point-of-sale-user/cache"
	"github.com/MamangRust/monolith-point-of-sale-user/handler"
	"github.com/MamangRust/monolith-point-of-sale-user/repository"
	"github.com/MamangRust/monolith-point-of-sale-user/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	roleAddr := os.Getenv("GRPC_ROLE_ADDR")
	if roleAddr == "" {
		roleAddr = "localhost:50052"
	}

	srv.Logger.Info("Connecting to Role service via gRPC", zap.String("addr", roleAddr))
	roleConn, err := server.NewGRPCClient(roleAddr)
	if err != nil {
		return nil, err
	}

	go func() {
		<-srv.Ctx.Done()
		srv.Logger.Info("Closing user service remote gRPC connections")
		roleConn.Close()
	}()

	roleQueryClient := pbroles.NewRoleQueryServiceClient(roleConn)
	userRoleClient := pbuserroles.NewUserRoleServiceClient(roleConn)

	guardRole := resilience.NewDependencyGuard("role", 5, 30, 100, 3*time.Second, srv.Logger)

	repos := repository.NewRepositories(srv.MongoDB, roleQueryClient, userRoleClient,
		repository.GuardOptions{
			Role: []adapter.GuardOption{adapter.WithDependencyGuard(guardRole)},
		})
	hash := hash.NewHashingPassword()
	traceLoggerObservability := observability.NewTraceLoggerObservability(srv.Logger)

	cacheMetrics, err := observability.NewCacheMetrics("user")
	if err != nil {
		return nil, err
	}

	cacheStore := cache.NewCacheStore(srv.Redis, srv.Logger, cacheMetrics)
	mencacheObj := mencache.NewMencache(cacheStore)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Repositories:  repos,
		Hash:          hash,
		Logger:        srv.Logger,
		Observability: traceLoggerObservability,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbusers.RegisterUserQueryServiceServer(gs, handlers.User)
	}

	return srv, nil
}
