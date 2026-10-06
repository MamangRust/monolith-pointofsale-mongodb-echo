package mencache

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-category/repository"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type CategoryQueryCache interface {
	GetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResult, *int, bool)
	SetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory, data []*repository.CategoryResult, total *int)

	GetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResultDeleteAt, *int, bool)
	SetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory, data []*repository.CategoryResultDeleteAt, total *int)

	GetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResultDeleteAt, *int, bool)
	SetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory, data []*repository.CategoryResultDeleteAt, total *int)

	GetCachedCategoryCache(ctx context.Context, id int) (*salemongo.Category, bool)
	SetCachedCategoryCache(ctx context.Context, data *salemongo.Category)
}

type CategoryCommandCache interface {
	DeleteCachedCategoryCache(ctx context.Context, id int)
	DeleteCachedCategoryAllCache(ctx context.Context)
}

// CategoryStatsCache caches aggregate category price statistics.
type CategoryStatsCache interface {
	GetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice) ([]*salemongo.CategoryMonthlyTotalPriceRow, bool)
	SetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice, res []*salemongo.CategoryMonthlyTotalPriceRow)
	GetCachedYearTotalPriceCache(ctx context.Context, year int) ([]*salemongo.CategoryYearlyTotalPriceRow, bool)
	SetCachedYearTotalPriceCache(ctx context.Context, year int, res []*salemongo.CategoryYearlyTotalPriceRow)
	GetCachedMonthPriceCache(ctx context.Context, year int) ([]*salemongo.CategoryMonthlyRow, bool)
	SetCachedMonthPriceCache(ctx context.Context, year int, res []*salemongo.CategoryMonthlyRow)
	GetCachedYearPriceCache(ctx context.Context, year int) ([]*salemongo.CategoryYearlyRow, bool)
	SetCachedYearPriceCache(ctx context.Context, year int, res []*salemongo.CategoryYearlyRow)
}

// CategoryStatsByIdCache caches per-category price statistics.
type CategoryStatsByIdCache interface {
	GetCachedMonthTotalPriceByIdCache(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*salemongo.CategoryMonthlyTotalPriceByIdRow, bool)
	SetCachedMonthTotalPriceByIdCache(ctx context.Context, req *requests.MonthTotalPriceCategory, res []*salemongo.CategoryMonthlyTotalPriceByIdRow)
	GetCachedYearTotalPriceByIdCache(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*salemongo.CategoryYearlyTotalPriceByIdRow, bool)
	SetCachedYearTotalPriceByIdCache(ctx context.Context, req *requests.YearTotalPriceCategory, res []*salemongo.CategoryYearlyTotalPriceByIdRow)
	GetCachedMonthPriceByIdCache(ctx context.Context, req *requests.MonthPriceId) ([]*salemongo.CategoryMonthlyByIdRow, bool)
	SetCachedMonthPriceByIdCache(ctx context.Context, req *requests.MonthPriceId, res []*salemongo.CategoryMonthlyByIdRow)
	GetCachedYearPriceByIdCache(ctx context.Context, req *requests.YearPriceId) ([]*salemongo.CategoryYearlyByIdRow, bool)
	SetCachedYearPriceByIdCache(ctx context.Context, req *requests.YearPriceId, res []*salemongo.CategoryYearlyByIdRow)
}

// CategoryStatsByMerchantCache caches per-merchant category price statistics.
type CategoryStatsByMerchantCache interface {
	GetCachedMonthTotalPriceByMerchantCache(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*salemongo.CategoryMonthlyTotalPriceByMerchantRow, bool)
	SetCachedMonthTotalPriceByMerchantCache(ctx context.Context, req *requests.MonthTotalPriceMerchant, res []*salemongo.CategoryMonthlyTotalPriceByMerchantRow)
	GetCachedYearTotalPriceByMerchantCache(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*salemongo.CategoryYearlyTotalPriceByMerchantRow, bool)
	SetCachedYearTotalPriceByMerchantCache(ctx context.Context, req *requests.YearTotalPriceMerchant, res []*salemongo.CategoryYearlyTotalPriceByMerchantRow)
	GetCachedMonthPriceByMerchantCache(ctx context.Context, req *requests.MonthPriceMerchant) ([]*salemongo.CategoryMonthlyByMerchantRow, bool)
	SetCachedMonthPriceByMerchantCache(ctx context.Context, req *requests.MonthPriceMerchant, res []*salemongo.CategoryMonthlyByMerchantRow)
	GetCachedYearPriceByMerchantCache(ctx context.Context, req *requests.YearPriceMerchant) ([]*salemongo.CategoryYearlyByMerchantRow, bool)
	SetCachedYearPriceByMerchantCache(ctx context.Context, req *requests.YearPriceMerchant, res []*salemongo.CategoryYearlyByMerchantRow)
}
