package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

// Local result types for cashier queries

type CashierResult struct {
	CashierID  int32
	MerchantID int32
	UserID     int32
	Name       string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	TotalCount int64
}

type CashierResultDeleteAt struct {
	CashierID  int32
	MerchantID int32
	UserID     int32
	Name       string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DeletedAt  *time.Time
	TotalCount int64
}

type MonthlyTotalSalesResult struct {
	Year       string
	Month      string
	TotalSales int64
}

type YearlyTotalSalesResult struct {
	Year       string
	TotalSales int64
}

type MonthlyCashierResult struct {
	Month       string
	CashierID   int32
	CashierName string
	OrderCount  int64
	TotalSales  int64
}

type YearlyCashierResult struct {
	Year        string
	CashierID   int32
	CashierName string
	OrderCount  int64
	TotalSales  int64
}

// The merchant/user lookups cashier validates against live in other bounded
// contexts and are reached through the shared gRPC adapters in
// pkg/adapter/{merchant,user} (see repository.go).

type CashierQueryRepository interface {
	FindAllCashiers(ctx context.Context, req *requests.FindAllCashiers) ([]*CashierResult, *int, error)
	FindById(ctx context.Context, cashier_id int) (*salemongo.Cashier, error)
	FindByActive(ctx context.Context, req *requests.FindAllCashiers) ([]*CashierResultDeleteAt, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllCashiers) ([]*CashierResultDeleteAt, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant) ([]*CashierResult, *int, error)
}

type CashierCommandRepository interface {
	CreateCashier(ctx context.Context, request *requests.CreateCashierRequest) (*salemongo.Cashier, error)
	UpdateCashier(ctx context.Context, request *requests.UpdateCashierRequest) (*salemongo.Cashier, error)
	TrashedCashier(ctx context.Context, cashier_id int) (*salemongo.Cashier, error)
	RestoreCashier(ctx context.Context, cashier_id int) (*salemongo.Cashier, error)
	DeleteCashierPermanent(ctx context.Context, cashier_id int) (bool, error)
	RestoreAllCashier(ctx context.Context) (bool, error)
	DeleteAllCashierPermanent(ctx context.Context) (bool, error)
}

// sqlc row type aliases for handler compatibility
type CashierRow = CashierResult
type CashierActiveRow = CashierResultDeleteAt
type CashierTrashedRow = CashierResultDeleteAt
type CashierByMerchantRow = CashierResult
type CashierMonthlySaleRow = salemongo.CashierMonthlySaleRow
type CashierMonthlySaleByIdRow = salemongo.CashierMonthlySaleByIdRow
type CashierMonthlySaleByMerchantRow = salemongo.CashierMonthlySaleByMerchantRow
type CashierMonthlyTotalSalesRow = salemongo.CashierMonthlyTotalSalesRow
type CashierYearlySaleRow = salemongo.CashierYearlySaleRow
type CashierYearlyTotalSalesRow = salemongo.CashierYearlyTotalSalesRow

// CashierStatByIdRepository provides sales statistics for a single cashier.
type CashierStatByIdRepository interface {
	GetMonthlyTotalSalesById(ctx context.Context, req *requests.MonthTotalSalesCashier) ([]*salemongo.CashierMonthlyTotalSalesRow, error)
	GetYearlyTotalSalesById(ctx context.Context, req *requests.YearTotalSalesCashier) ([]*salemongo.CashierYearlyTotalSalesRow, error)
	GetMonthlyCashierById(ctx context.Context, req *requests.MonthCashierId) ([]*salemongo.CashierMonthlySaleByIdRow, error)
	GetYearlyCashierById(ctx context.Context, req *requests.YearCashierId) ([]*salemongo.CashierYearlySaleRow, error)
}

// CashierStatsRepository provides aggregate sales statistics across cashiers.
type CashierStatsRepository interface {
	GetMonthlyTotalSales(ctx context.Context, req *requests.MonthTotalSales) ([]*salemongo.CashierMonthlyTotalSalesRow, error)
	GetYearlyTotalSales(ctx context.Context, year int) ([]*salemongo.CashierYearlyTotalSalesRow, error)
	GetMonthyCashier(ctx context.Context, year int) ([]*salemongo.CashierMonthlySaleRow, error)
	GetYearlyCashier(ctx context.Context, year int) ([]*salemongo.CashierYearlySaleRow, error)
}

// CashierStatByMerchantRepository provides sales statistics grouped by merchant.
type CashierStatByMerchantRepository interface {
	GetMonthlyTotalSalesByMerchant(ctx context.Context, req *requests.MonthTotalSalesMerchant) ([]*salemongo.CashierMonthlyTotalSalesRow, error)
	GetYearlyTotalSalesByMerchant(ctx context.Context, req *requests.YearTotalSalesMerchant) ([]*salemongo.CashierYearlyTotalSalesRow, error)
	GetMonthlyCashierByMerchant(ctx context.Context, req *requests.MonthCashierMerchant) ([]*salemongo.CashierMonthlySaleByMerchantRow, error)
	GetYearlyCashierByMerchant(ctx context.Context, req *requests.YearCashierMerchant) ([]*salemongo.CashierYearlySaleRow, error)
}
