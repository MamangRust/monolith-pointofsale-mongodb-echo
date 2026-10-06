package apps

import (
	"os"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pb"
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	pbuserroles "github.com/MamangRust/monolith-point-of-sale-pb/user_roles"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"

	mencache "github.com/MamangRust/monolith-point-of-sale-auth/cache"
	"github.com/MamangRust/monolith-point-of-sale-auth/handler"
	"github.com/MamangRust/monolith-point-of-sale-auth/repository"
	"github.com/MamangRust/monolith-point-of-sale-auth/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	"github.com/MamangRust/monolith-point-of-sale-pkg/auth"
	"github.com/MamangRust/monolith-point-of-sale-pkg/hash"
	"github.com/MamangRust/monolith-point-of-sale-pkg/kafka"
	"github.com/MamangRust/monolith-point-of-sale-pkg/outbox"
	"github.com/MamangRust/monolith-point-of-sale-pkg/resilience"
	"github.com/MamangRust/monolith-point-of-sale-pkg/server"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"github.com/spf13/viper"
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

	tokenManager, err := auth.NewManager(viper.GetString("SECRET_KEY"))
	if err != nil {
		return nil, err
	}

	// users and roles/user_roles now live in the User and Role services, so
	// auth reaches them over gRPC instead of reading their collections directly.
	userAddr := getEnv("GRPC_USER_ADDR", "localhost:50053")
	roleAddr := getEnv("GRPC_ROLE_ADDR", "localhost:50052")

	srv.Logger.Info("Connecting to remote gRPC services from Auth service",
		zap.String("user", userAddr), zap.String("role", roleAddr))

	userConn, err := server.NewGRPCClient(userAddr)
	if err != nil {
		return nil, err
	}

	roleConn, err := server.NewGRPCClient(roleAddr)
	if err != nil {
		userConn.Close()
		return nil, err
	}

	go func() {
		<-srv.Ctx.Done()
		srv.Logger.Info("Closing auth service remote gRPC connections")
		userConn.Close()
		roleConn.Close()
	}()

	// Per-call deadline + circuit breaker + bulkhead for each identity
	// dependency (user, role, user_role).
	guardUser := resilience.NewDependencyGuard("user", 5, 30, 100, 3*time.Second, srv.Logger)
	guardRole := resilience.NewDependencyGuard("role", 5, 30, 100, 3*time.Second, srv.Logger)
	guardUserRole := resilience.NewDependencyGuard("user_role", 5, 30, 100, 3*time.Second, srv.Logger)

	repos := repository.NewRepositories(
		srv.MongoDB,
		pbusers.NewUserQueryServiceClient(userConn),
		pbusers.NewUserCommandServiceClient(userConn),
		pbroles.NewRoleQueryServiceClient(roleConn),
		pbuserroles.NewUserRoleServiceClient(roleConn),
		repository.GuardOptions{
			User:     []adapter.GuardOption{adapter.WithDependencyGuard(guardUser)},
			Role:     []adapter.GuardOption{adapter.WithDependencyGuard(guardRole)},
			UserRole: []adapter.GuardOption{adapter.WithDependencyGuard(guardUserRole)},
		},
	)
	hash := hash.NewHashingPassword()

	var myKafka *kafka.Kafka
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		myKafka = kafka.NewKafka(srv.Logger, []string{brokers})
	}

	mencacheObj := mencache.NewMencache(srv.CacheStore)

	outboxService := outbox.NewOutboxService(srv.MongoDB, myKafka, srv.Logger)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Repositories:  repos,
		Hash:          hash,
		Token:         tokenManager,
		Logger:        srv.Logger,
		Kafka:         myKafka,
		Observability: observability.NewTraceLoggerObservability(srv.Logger),
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterAuthServiceServer(gs, handlers.Auth)
	}

	go outboxService.Start(srv.Ctx, outbox.OutboxRelayInterval, outbox.OutboxRelayBatchSize)

	return srv, nil
}
