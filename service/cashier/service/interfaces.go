package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-cashier/repository"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type CashierQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllCashiers) ([]*repository.CashierResult, *int, error)
	FindById(ctx context.Context, cashierID int) (*salemongo.Cashier, error)
	FindByActive(ctx context.Context, req *requests.FindAllCashiers) ([]*repository.CashierResultDeleteAt, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllCashiers) ([]*repository.CashierResultDeleteAt, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant) ([]*repository.CashierResult, *int, error)
}

type CashierCommandService interface {
	CreateCashier(ctx context.Context, req *requests.CreateCashierRequest) (*salemongo.Cashier, error)
	UpdateCashier(ctx context.Context, req *requests.UpdateCashierRequest) (*salemongo.Cashier, error)
	TrashedCashier(ctx context.Context, cashierID int) (*salemongo.Cashier, error)
	RestoreCashier(ctx context.Context, cashierID int) (*salemongo.Cashier, error)
	DeleteCashierPermanent(ctx context.Context, cashierID int) (bool, error)
	RestoreAllCashier(ctx context.Context) (bool, error)
	DeleteAllCashierPermanent(ctx context.Context) (bool, error)
}

// CashierStatsService exposes aggregate cashier sales statistics.
type CashierStatsService interface {
	FindMonthlyTotalSales(ctx context.Context, req *requests.MonthTotalSales) ([]*salemongo.CashierMonthlyTotalSalesRow, error)
	FindYearlyTotalSales(ctx context.Context, year int) ([]*salemongo.CashierYearlyTotalSalesRow, error)
	FindMonthlySales(ctx context.Context, year int) ([]*salemongo.CashierMonthlySaleRow, error)
	FindYearlySales(ctx context.Context, year int) ([]*salemongo.CashierYearlySaleRow, error)
}

// CashierStatsByIdService exposes per-cashier sales statistics.
type CashierStatsByIdService interface {
	FindMonthlyTotalSalesById(ctx context.Context, req *requests.MonthTotalSalesCashier) ([]*salemongo.CashierMonthlyTotalSalesRow, error)
	FindYearlyTotalSalesById(ctx context.Context, req *requests.YearTotalSalesCashier) ([]*salemongo.CashierYearlyTotalSalesRow, error)
	FindMonthlyCashierById(ctx context.Context, req *requests.MonthCashierId) ([]*salemongo.CashierMonthlySaleByIdRow, error)
	FindYearlyCashierById(ctx context.Context, req *requests.YearCashierId) ([]*salemongo.CashierYearlySaleRow, error)
}

// CashierStatsByMerchant exposes per-merchant cashier sales statistics.
type CashierStatsByMerchant interface {
	FindMonthlyTotalSalesByMerchant(ctx context.Context, req *requests.MonthTotalSalesMerchant) ([]*salemongo.CashierMonthlyTotalSalesRow, error)
	FindYearlyTotalSalesByMerchant(ctx context.Context, req *requests.YearTotalSalesMerchant) ([]*salemongo.CashierYearlyTotalSalesRow, error)
	FindMonthlyCashierByMerchant(ctx context.Context, req *requests.MonthCashierMerchant) ([]*salemongo.CashierMonthlySaleByMerchantRow, error)
	FindYearlyCashierByMerchant(ctx context.Context, req *requests.YearCashierMerchant) ([]*salemongo.CashierYearlySaleRow, error)
}
