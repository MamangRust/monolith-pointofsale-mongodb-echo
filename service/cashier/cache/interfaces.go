package mencache

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-cashier/repository"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type CashierQueryCache interface {
	GetCachedCashiersCache(ctx context.Context, req *requests.FindAllCashiers) ([]*repository.CashierResult, *int, bool)
	SetCachedCashiersCache(ctx context.Context, req *requests.FindAllCashiers, res []*repository.CashierResult, total *int)

	GetCachedCashier(ctx context.Context, cashierID int) (*salemongo.Cashier, bool)
	SetCachedCashier(ctx context.Context, res *salemongo.Cashier)

	GetCachedCashiersActive(ctx context.Context, req *requests.FindAllCashiers) ([]*repository.CashierResultDeleteAt, *int, bool)
	SetCachedCashiersActive(ctx context.Context, req *requests.FindAllCashiers, res []*repository.CashierResultDeleteAt, total *int)

	GetCachedCashiersTrashed(ctx context.Context, req *requests.FindAllCashiers) ([]*repository.CashierResultDeleteAt, *int, bool)
	SetCachedCashiersTrashed(ctx context.Context, req *requests.FindAllCashiers, res []*repository.CashierResultDeleteAt, total *int)

	GetCachedCashiersByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant) ([]*repository.CashierResult, *int, bool)
	SetCachedCashiersByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant, res []*repository.CashierResult, total *int)
}

type CashierCommandCache interface {
	DeleteCashierCache(ctx context.Context, id int)
	DeleteCashierListCache(ctx context.Context)
}

// CashierStatsCache caches aggregate cashier sales statistics.
type CashierStatsCache interface {
	GetMonthlyTotalSalesCache(ctx context.Context, req *requests.MonthTotalSales) ([]*salemongo.CashierMonthlyTotalSalesRow, bool)
	SetMonthlyTotalSalesCache(ctx context.Context, req *requests.MonthTotalSales, res []*salemongo.CashierMonthlyTotalSalesRow)
	GetYearlyTotalSalesCache(ctx context.Context, year int) ([]*salemongo.CashierYearlyTotalSalesRow, bool)
	SetYearlyTotalSalesCache(ctx context.Context, year int, res []*salemongo.CashierYearlyTotalSalesRow)
	GetMonthlySalesCache(ctx context.Context, year int) ([]*salemongo.CashierMonthlySaleRow, bool)
	SetMonthlySalesCache(ctx context.Context, year int, res []*salemongo.CashierMonthlySaleRow)
	GetYearlySalesCache(ctx context.Context, year int) ([]*salemongo.CashierYearlySaleRow, bool)
	SetYearlySalesCache(ctx context.Context, year int, res []*salemongo.CashierYearlySaleRow)
}

// CashierStatsByIdCache caches per-cashier sales statistics.
type CashierStatsByIdCache interface {
	GetMonthlyTotalSalesByIdCache(ctx context.Context, req *requests.MonthTotalSalesCashier) ([]*salemongo.CashierMonthlyTotalSalesRow, bool)
	SetMonthlyTotalSalesByIdCache(ctx context.Context, req *requests.MonthTotalSalesCashier, res []*salemongo.CashierMonthlyTotalSalesRow)
	GetYearlyTotalSalesByIdCache(ctx context.Context, req *requests.YearTotalSalesCashier) ([]*salemongo.CashierYearlyTotalSalesRow, bool)
	SetYearlyTotalSalesByIdCache(ctx context.Context, req *requests.YearTotalSalesCashier, res []*salemongo.CashierYearlyTotalSalesRow)
	GetMonthlyCashierByIdCache(ctx context.Context, req *requests.MonthCashierId) ([]*salemongo.CashierMonthlySaleByIdRow, bool)
	SetMonthlyCashierByIdCache(ctx context.Context, req *requests.MonthCashierId, res []*salemongo.CashierMonthlySaleByIdRow)
	GetYearlyCashierByIdCache(ctx context.Context, req *requests.YearCashierId) ([]*salemongo.CashierYearlySaleRow, bool)
	SetYearlyCashierByIdCache(ctx context.Context, req *requests.YearCashierId, res []*salemongo.CashierYearlySaleRow)
}

// CashierStatsByMerchantCache caches per-merchant cashier sales statistics.
type CashierStatsByMerchantCache interface {
	GetMonthlyTotalSalesByMerchantCache(ctx context.Context, req *requests.MonthTotalSalesMerchant) ([]*salemongo.CashierMonthlyTotalSalesRow, bool)
	SetMonthlyTotalSalesByMerchantCache(ctx context.Context, req *requests.MonthTotalSalesMerchant, res []*salemongo.CashierMonthlyTotalSalesRow)
	GetYearlyTotalSalesByMerchantCache(ctx context.Context, req *requests.YearTotalSalesMerchant) ([]*salemongo.CashierYearlyTotalSalesRow, bool)
	SetYearlyTotalSalesByMerchantCache(ctx context.Context, req *requests.YearTotalSalesMerchant, res []*salemongo.CashierYearlyTotalSalesRow)
	GetMonthlyCashierByMerchantCache(ctx context.Context, req *requests.MonthCashierMerchant) ([]*salemongo.CashierMonthlySaleByMerchantRow, bool)
	SetMonthlyCashierByMerchantCache(ctx context.Context, req *requests.MonthCashierMerchant, res []*salemongo.CashierMonthlySaleByMerchantRow)
	GetYearlyCashierByMerchantCache(ctx context.Context, req *requests.YearCashierMerchant) ([]*salemongo.CashierYearlySaleRow, bool)
	SetYearlyCashierByMerchantCache(ctx context.Context, req *requests.YearCashierMerchant, res []*salemongo.CashierYearlySaleRow)
}
