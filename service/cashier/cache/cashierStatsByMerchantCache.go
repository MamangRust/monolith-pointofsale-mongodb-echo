package mencache

import (
	"context"
	"fmt"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

const (
	cashierStatsMonthTotalSalesByMerchantCacheKey = "cashier:stats:month:%d:year:%d:id:%d"
	cashierStatsYearTotalSalesByMerchantCacheKey  = "cashier:stats:year:%d:merchant:%d"

	cashierStatsMonthSalesByMerchantCacheKey = "cashier:stats:month:%d:merchant:%d"
	cashierStatsYearSalesByMerchantCacheKey  = "cashier:stats:year:%d:merchant:%d"
)

type cashierStatsByMerchantCache struct {
	store *cache.CacheStore
}

func NewCashierStatsByMerchantCache(store *cache.CacheStore) CashierStatsByMerchantCache {
	return &cashierStatsByMerchantCache{store: store}
}

func (s *cashierStatsByMerchantCache) GetMonthlyTotalSalesByMerchantCache(ctx context.Context, req *requests.MonthTotalSalesMerchant) ([]*salemongo.CashierMonthlyTotalSalesRow, bool) {
	key := fmt.Sprintf(cashierStatsMonthTotalSalesByMerchantCacheKey, req.Month, req.Year, req.MerchantID)
	result, found := cache.GetFromCache[[]*salemongo.CashierMonthlyTotalSalesRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByMerchantCache) SetMonthlyTotalSalesByMerchantCache(ctx context.Context, req *requests.MonthTotalSalesMerchant, res []*salemongo.CashierMonthlyTotalSalesRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthTotalSalesByMerchantCacheKey, req.Month, req.Year, req.MerchantID)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsByMerchantCache) GetYearlyTotalSalesByMerchantCache(ctx context.Context, req *requests.YearTotalSalesMerchant) ([]*salemongo.CashierYearlyTotalSalesRow, bool) {
	key := fmt.Sprintf(cashierStatsYearTotalSalesByMerchantCacheKey, req.Year, req.MerchantID)
	result, found := cache.GetFromCache[[]*salemongo.CashierYearlyTotalSalesRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByMerchantCache) SetYearlyTotalSalesByMerchantCache(ctx context.Context, req *requests.YearTotalSalesMerchant, res []*salemongo.CashierYearlyTotalSalesRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearTotalSalesByMerchantCacheKey, req.Year, req.MerchantID)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsByMerchantCache) GetMonthlyCashierByMerchantCache(ctx context.Context, req *requests.MonthCashierMerchant) ([]*salemongo.CashierMonthlySaleByMerchantRow, bool) {
	key := fmt.Sprintf(cashierStatsMonthSalesByMerchantCacheKey, req.Year, req.MerchantID)
	result, found := cache.GetFromCache[[]*salemongo.CashierMonthlySaleByMerchantRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByMerchantCache) SetMonthlyCashierByMerchantCache(ctx context.Context, req *requests.MonthCashierMerchant, res []*salemongo.CashierMonthlySaleByMerchantRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthSalesByMerchantCacheKey, req.Year, req.MerchantID)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsByMerchantCache) GetYearlyCashierByMerchantCache(ctx context.Context, req *requests.YearCashierMerchant) ([]*salemongo.CashierYearlySaleRow, bool) {
	key := fmt.Sprintf(cashierStatsYearSalesByMerchantCacheKey, req.Year, req.MerchantID)
	result, found := cache.GetFromCache[[]*salemongo.CashierYearlySaleRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByMerchantCache) SetYearlyCashierByMerchantCache(ctx context.Context, req *requests.YearCashierMerchant, res []*salemongo.CashierYearlySaleRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearSalesByMerchantCacheKey, req.Year, req.MerchantID)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}
