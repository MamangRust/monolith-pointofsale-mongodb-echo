package mencache

import (
	"context"
	"fmt"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

const (
	transactionMonthAmountSuccessKey = "transaction:month:amount:success:month:%d:year:%d"
	transactionMonthAmountFailedKey  = "transaction:month:amount:failed:month:%d:year:%d"

	transactionYearAmountSuccessKey = "transaction:year:amount:success:year:%d"
	transactionYearAmountFailedKey  = "transaction:year:amount:failed:year:%d"

	transactionMonthMethodSuccessKey = "transaction:month:method:success:month:%d:year:%d"
	transactionMonthMethodFailedKey  = "transaction:month:method:failed:month:%d:year:%d"

	transactionYearMethodSuccessKey = "transaction:year:method:success:year:%d"
	transactionYearMethodFailedKey  = "transaction:year:method:failed:year:%d"
)

type transactionStatsCache struct {
	store *cache.CacheStore
}

func NewTransactionStatsCache(store *cache.CacheStore) TransactionStatsCache {
	return &transactionStatsCache{store: store}
}

func (t *transactionStatsCache) GetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransaction) ([]*salemongo.TransactionMonthlyAmountSuccessRow, bool) {
	key := fmt.Sprintf(transactionMonthAmountSuccessKey, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*salemongo.TransactionMonthlyAmountSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransaction, res []*salemongo.TransactionMonthlyAmountSuccessRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionMonthAmountSuccessKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearAmountSuccessCached(ctx context.Context, year int) ([]*salemongo.TransactionYearlyAmountSuccessRow, bool) {
	key := fmt.Sprintf(transactionYearAmountSuccessKey, year)

	result, found := cache.GetFromCache[[]*salemongo.TransactionYearlyAmountSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedYearAmountSuccessCached(ctx context.Context, year int, res []*salemongo.TransactionYearlyAmountSuccessRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionYearAmountSuccessKey, year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransaction) ([]*salemongo.TransactionMonthlyAmountFailedRow, bool) {
	key := fmt.Sprintf(transactionMonthAmountFailedKey, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*salemongo.TransactionMonthlyAmountFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransaction, res []*salemongo.TransactionMonthlyAmountFailedRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionMonthAmountFailedKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearAmountFailedCached(ctx context.Context, year int) ([]*salemongo.TransactionYearlyAmountFailedRow, bool) {
	key := fmt.Sprintf(transactionYearAmountFailedKey, year)

	result, found := cache.GetFromCache[[]*salemongo.TransactionYearlyAmountFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedYearAmountFailedCached(ctx context.Context, year int, res []*salemongo.TransactionYearlyAmountFailedRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionYearAmountFailedKey, year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransaction) ([]*salemongo.TransactionMonthlyMethodSuccessRow, bool) {
	key := fmt.Sprintf(transactionMonthMethodSuccessKey, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*salemongo.TransactionMonthlyMethodSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransaction, res []*salemongo.TransactionMonthlyMethodSuccessRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionMonthMethodSuccessKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearMethodSuccessCached(ctx context.Context, year int) ([]*salemongo.TransactionYearlyMethodSuccessRow, bool) {
	key := fmt.Sprintf(transactionYearMethodSuccessKey, year)

	result, found := cache.GetFromCache[[]*salemongo.TransactionYearlyMethodSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedYearMethodSuccessCached(ctx context.Context, year int, res []*salemongo.TransactionYearlyMethodSuccessRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionYearMethodSuccessKey, year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransaction) ([]*salemongo.TransactionMonthlyMethodFailedRow, bool) {
	key := fmt.Sprintf(transactionMonthMethodFailedKey, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*salemongo.TransactionMonthlyMethodFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransaction, res []*salemongo.TransactionMonthlyMethodFailedRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionMonthMethodFailedKey, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsCache) GetCachedYearMethodFailedCached(ctx context.Context, year int) ([]*salemongo.TransactionYearlyMethodFailedRow, bool) {
	key := fmt.Sprintf(transactionYearMethodFailedKey, year)

	result, found := cache.GetFromCache[[]*salemongo.TransactionYearlyMethodFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsCache) SetCachedYearMethodFailedCached(ctx context.Context, year int, res []*salemongo.TransactionYearlyMethodFailedRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactionYearMethodFailedKey, year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}
