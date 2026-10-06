package mencache

import (
	"context"
	"fmt"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

const (
	cashierStatsMonthTotalSalesCacheKey = "cashier:stats:month:%d:year:%d"
	cashierStatsYearTotalSalesCacheKey  = "cashier:stats:year:%d"

	cashierStatsMonthSalesCacheKey = "cashier:stats:month:%d"
	cashierStatsYearSalesCacheKey  = "cashier:stats:year:%d"
)

type cashierStatsCache struct {
	store *cache.CacheStore
}

func NewCashierStatsCache(store *cache.CacheStore) CashierStatsCache {
	return &cashierStatsCache{store: store}
}

func (s *cashierStatsCache) GetMonthlyTotalSalesCache(ctx context.Context, req *requests.MonthTotalSales) ([]*salemongo.CashierMonthlyTotalSalesRow, bool) {
	key := fmt.Sprintf(cashierStatsMonthTotalSalesCacheKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*salemongo.CashierMonthlyTotalSalesRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsCache) SetMonthlyTotalSalesCache(ctx context.Context, req *requests.MonthTotalSales, res []*salemongo.CashierMonthlyTotalSalesRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthTotalSalesCacheKey, req.Month, req.Year)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsCache) GetYearlyTotalSalesCache(ctx context.Context, year int) ([]*salemongo.CashierYearlyTotalSalesRow, bool) {
	key := fmt.Sprintf(cashierStatsYearTotalSalesCacheKey, year)
	result, found := cache.GetFromCache[[]*salemongo.CashierYearlyTotalSalesRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsCache) SetYearlyTotalSalesCache(ctx context.Context, year int, res []*salemongo.CashierYearlyTotalSalesRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearTotalSalesCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsCache) GetMonthlySalesCache(ctx context.Context, year int) ([]*salemongo.CashierMonthlySaleRow, bool) {
	key := fmt.Sprintf(cashierStatsMonthSalesCacheKey, year)
	result, found := cache.GetFromCache[[]*salemongo.CashierMonthlySaleRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsCache) SetMonthlySalesCache(ctx context.Context, year int, res []*salemongo.CashierMonthlySaleRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthSalesCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsCache) GetYearlySalesCache(ctx context.Context, year int) ([]*salemongo.CashierYearlySaleRow, bool) {
	key := fmt.Sprintf(cashierStatsYearSalesCacheKey, year)
	result, found := cache.GetFromCache[[]*salemongo.CashierYearlySaleRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsCache) SetYearlySalesCache(ctx context.Context, year int, res []*salemongo.CashierYearlySaleRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearSalesCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}
