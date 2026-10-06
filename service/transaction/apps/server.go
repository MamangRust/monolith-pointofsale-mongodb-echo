package apps

import (
	"os"
	"time"

	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	pbtransactions "github.com/MamangRust/monolith-point-of-sale-pb/transactions"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	"github.com/MamangRust/monolith-point-of-sale-pkg/kafka"
	"github.com/MamangRust/monolith-point-of-sale-pkg/outbox"
	"github.com/MamangRust/monolith-point-of-sale-pkg/resilience"
	"github.com/MamangRust/monolith-point-of-sale-pkg/server"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	mencache "github.com/MamangRust/monolith-point-of-sale-transacton/cache"
	"github.com/MamangRust/monolith-point-of-sale-transacton/handler"
	"github.com/MamangRust/monolith-point-of-sale-transacton/repository"
	"github.com/MamangRust/monolith-point-of-sale-transacton/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	// Connect to merchant service
	merchantAddr := os.Getenv("GRPC_MERCHANT_ADDR")
	if merchantAddr == "" {
		merchantAddr = "localhost:50052"
	}
	srv.Logger.Info("Connecting to Merchant service via gRPC", zap.String("addr", merchantAddr))
	merchantConn, err := server.NewGRPCClient(merchantAddr)
	if err != nil {
		return nil, err
	}
	go func() {
		<-srv.Ctx.Done()
		merchantConn.Close()
	}()
	merchantClient := pbmerchants.NewMerchantQueryServiceClient(merchantConn)

	// Connect to order service
	orderAddr := os.Getenv("GRPC_ORDER_ADDR")
	if orderAddr == "" {
		orderAddr = "localhost:50056"
	}
	srv.Logger.Info("Connecting to Order service via gRPC", zap.String("addr", orderAddr))
	orderConn, err := server.NewGRPCClient(orderAddr)
	if err != nil {
		return nil, err
	}
	go func() {
		<-srv.Ctx.Done()
		orderConn.Close()
	}()
	orderClient := pborders.NewOrderQueryServiceClient(orderConn)

	// Connect to cashier service
	cashierAddr := os.Getenv("GRPC_CASHIER_ADDR")
	if cashierAddr == "" {
		cashierAddr = "localhost:50055"
	}
	srv.Logger.Info("Connecting to Cashier service via gRPC", zap.String("addr", cashierAddr))
	cashierConn, err := server.NewGRPCClient(cashierAddr)
	if err != nil {
		return nil, err
	}
	go func() {
		<-srv.Ctx.Done()
		cashierConn.Close()
	}()
	cashierClient := pbcashiers.NewCashierQueryServiceClient(cashierConn)

	// Connect to order item service
	orderItemAddr := os.Getenv("GRPC_ORDER_ITEM_ADDR")
	if orderItemAddr == "" {
		orderItemAddr = "localhost:50060"
	}
	srv.Logger.Info("Connecting to Order Item service via gRPC", zap.String("addr", orderItemAddr))
	orderItemConn, err := server.NewGRPCClient(orderItemAddr)
	if err != nil {
		return nil, err
	}
	go func() {
		<-srv.Ctx.Done()
		orderItemConn.Close()
	}()
	orderItemClient := pborder_items.NewOrderItemQueryServiceClient(orderItemConn)

	guardCashier := resilience.NewDependencyGuard("cashier", 5, 30, 100, 3*time.Second, srv.Logger)
	guardMerchant := resilience.NewDependencyGuard("merchant", 5, 30, 100, 3*time.Second, srv.Logger)
	guardOrder := resilience.NewDependencyGuard("order", 5, 30, 100, 3*time.Second, srv.Logger)
	guardOrderItem := resilience.NewDependencyGuard("order_item", 5, 30, 100, 3*time.Second, srv.Logger)

	repos := repository.NewRepositories(srv.MongoDB, cashierClient, merchantClient, orderClient, orderItemClient,
		repository.GuardOptions{
			Cashier:   []adapter.GuardOption{adapter.WithDependencyGuard(guardCashier)},
			Merchant:  []adapter.GuardOption{adapter.WithDependencyGuard(guardMerchant)},
			Order:     []adapter.GuardOption{adapter.WithDependencyGuard(guardOrder)},
			OrderItem: []adapter.GuardOption{adapter.WithDependencyGuard(guardOrderItem)},
		})

	// Kafka bersifat opsional: tanpa KAFKA_BROKERS (mis. E2E lokal tanpa kafka)
	// service tetap jalan dan event email di-skip (guard s.kafka != nil).
	var myKafka *kafka.Kafka
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		myKafka = kafka.NewKafka(srv.Logger, []string{brokers})
	}

	mencacheObj := mencache.NewMencache(srv.CacheStore)

	traceLoggerObservability := observability.NewTraceLoggerObservability(srv.Logger)

	outboxService := outbox.NewOutboxService(srv.MongoDB, myKafka, srv.Logger)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Repositories:  repos,
		Logger:        srv.Logger,
		Kafka:         myKafka,
		Observability: traceLoggerObservability,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbtransactions.RegisterTransactionQueryServiceServer(gs, handlers.Transaction)
	}

	// Start the outbox relay so events committed with the business writes are
	// published to Kafka with durable retry and dead-letter semantics.
	go outboxService.Start(srv.Ctx, outbox.OutboxRelayInterval, outbox.OutboxRelayBatchSize)

	return srv, nil
}
