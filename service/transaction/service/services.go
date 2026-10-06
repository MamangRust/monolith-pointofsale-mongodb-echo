package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/kafka"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	mencache "github.com/MamangRust/monolith-point-of-sale-transacton/cache"
	"github.com/MamangRust/monolith-point-of-sale-transacton/repository"
)

type Service struct {
	TransactionQuery           TransactionQueryService
	TransactionCommand         TransactionCommandService
	TransactionStats           TransactionStatsService
	TransactionStatsByMerchant TransactionStatsByMerchantQueryService
}

type Deps struct {
	Ctx           context.Context
	Kafka         *kafka.Kafka
	Mencache      mencache.Mencache
	CacheStore    *cache.CacheStore
	Repositories  *repository.Repositories
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	var kafkaPublisher EmailEventPublisher
	if deps.Kafka != nil {
		kafkaPublisher = deps.Kafka
	}

	return &Service{
		TransactionQuery:           NewTransactionQueryService(deps.Mencache, deps.Repositories.TransactionQuery, deps.Logger, deps.Observability),
		TransactionCommand:         NewTransactionCommandService(kafkaPublisher, deps.Mencache, deps.Repositories.CashierQuery, deps.Repositories.MerchantQuery, deps.Repositories.TransactionQuery, deps.Repositories.TransactionCommand, deps.Repositories.OrderQuery, deps.Repositories.OrderItemQuery, deps.Logger, deps.Observability),
		TransactionStats:           NewTransactionStatsService(mencache.NewTransactionStatsCache(deps.CacheStore), deps.Repositories.TransactionStats, deps.Logger, deps.Observability),
		TransactionStatsByMerchant: NewTransactionStatsByMerchantQueryService(mencache.NewTransactionStatsByMerchantCache(deps.CacheStore), deps.Repositories.TransactionStatsByMerchant, deps.Logger, deps.Observability),
	}
}
