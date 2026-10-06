package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-point-of-sale-category/cache"
	"github.com/MamangRust/monolith-point-of-sale-category/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
)

type Service struct {
	CategoryQuery           CategoryQueryService
	CategoryCommand         CategoryCommandService
	CategoryStats           CategoryStatsService
	CategoryStatsById       CategoryStatsByIdService
	CategoryStatsByMerchant CategoryStatsByMerchantQueryService
}

type Deps struct {
	Ctx           context.Context
	Mencache      mencache.Mencache
	CacheStore    *cache.CacheStore
	Repositories  *repository.Repositories
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		CategoryQuery: NewCategoryQueryService(&categoryQueryDeps{
			Cache:         deps.Mencache,
			CategoryQuery: deps.Repositories.CategoryQuery,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
		CategoryCommand: NewCategoryCommandService(&categoryCommandDeps{
			Cache:           deps.Mencache,
			CategoryQuery:   deps.Repositories.CategoryQuery,
			CategoryCommand: deps.Repositories.CategoryCommand,
			Logger:          deps.Logger,
			Observability:   deps.Observability,
		}),
		CategoryStats: NewCategoryStatsService(&categoryStatsDeps{
			Cache:                   mencache.NewCategoryStatsCache(deps.CacheStore),
			CategoryStatsRepository: deps.Repositories.CategoryStats,
			Logger:                  deps.Logger,
			Observability:           deps.Observability,
		}),
		CategoryStatsById: NewCategoryStatsByIdService(&categoryStatsByIdDeps{
			Cache:                       mencache.NewCategoryStatsByIdCache(deps.CacheStore),
			CategoryStatsByIdRepository: deps.Repositories.CategoryStatsById,
			Logger:                      deps.Logger,
			Observability:               deps.Observability,
		}),
		CategoryStatsByMerchant: NewCategoryStatsByMerchantQueryService(&categoryStatsByMerchantDeps{
			Cache:                             mencache.NewCategoryStatsByMerchantCache(deps.CacheStore),
			CategoryStatsByMerchantRepository: deps.Repositories.CategoryStatsByMerchant,
			Logger:                            deps.Logger,
			Observability:                     deps.Observability,
		}),
	}
}
