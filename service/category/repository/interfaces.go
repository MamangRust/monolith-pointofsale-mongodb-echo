package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

// Local result types for stats queries

type CategoryResult struct {
	CategoryID   int32
	Name         string
	Description  *string
	SlugCategory *string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	TotalCount   int64
}

type CategoryResultDeleteAt struct {
	CategoryID   int32
	Name         string
	Description  *string
	SlugCategory *string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	DeletedAt    *time.Time
	TotalCount   int64
}

type CategoryQueryRepository interface {
	FindAllCategory(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResult, *int, error)
	FindById(ctx context.Context, category_id int) (*salemongo.Category, error)
	FindByNameAndId(ctx context.Context, req *requests.CategoryNameAndId) (*salemongo.Category, error)
	FindByName(ctx context.Context, name string) (*salemongo.Category, error)
	FindByIds(ctx context.Context, ids []int) ([]*salemongo.Category, error)
	FindByIdTrashed(ctx context.Context, category_id int) (*salemongo.Category, error)
	FindByActive(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResultDeleteAt, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResultDeleteAt, *int, error)
}

type CategoryCommandRepository interface {
	CreateCategory(ctx context.Context, request *requests.CreateCategoryRequest) (*salemongo.Category, error)
	UpdateCategory(ctx context.Context, request *requests.UpdateCategoryRequest) (*salemongo.Category, error)
	TrashedCategory(ctx context.Context, category_id int) (*salemongo.Category, error)
	RestoreCategory(ctx context.Context, category_id int) (*salemongo.Category, error)
	DeleteCategoryPermanently(ctx context.Context, category_id int) (bool, error)
	RestoreAllCategories(ctx context.Context) (bool, error)
	DeleteAllPermanentCategories(ctx context.Context) (bool, error)
}

// sqlc row type aliases for handler compatibility
type CategoryRow = CategoryResult
type CategoryActiveRow = CategoryResultDeleteAt
type CategoryTrashedRow = CategoryResultDeleteAt
type CategoryMonthlyRow = salemongo.CategoryMonthlyRow
type CategoryMonthlyByIdRow = salemongo.CategoryMonthlyByIdRow
type CategoryMonthlyByMerchantRow = salemongo.CategoryMonthlyByMerchantRow
type CategoryMonthlyTotalPriceRow = salemongo.CategoryMonthlyTotalPriceRow
type CategoryMonthlyTotalPriceByIdRow = salemongo.CategoryMonthlyTotalPriceByIdRow
type CategoryMonthlyTotalPriceByMerchantRow = salemongo.CategoryMonthlyTotalPriceByMerchantRow
type CategoryYearlyRow = salemongo.CategoryYearlyRow
type CategoryYearlyByIdRow = salemongo.CategoryYearlyByIdRow
type CategoryYearlyByMerchantRow = salemongo.CategoryYearlyByMerchantRow
type CategoryYearlyTotalPriceRow = salemongo.CategoryYearlyTotalPriceRow
type CategoryYearlyTotalPriceByIdRow = salemongo.CategoryYearlyTotalPriceByIdRow
type CategoryYearlyTotalPriceByMerchantRow = salemongo.CategoryYearlyTotalPriceByMerchantRow

// CategoryStatsByIdRepository provides price statistics for a single category.
type CategoryStatsByIdRepository interface {
	GetMonthlyTotalPriceById(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*salemongo.CategoryMonthlyTotalPriceByIdRow, error)
	GetYearlyTotalPricesById(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*salemongo.CategoryYearlyTotalPriceByIdRow, error)
	GetMonthPriceById(ctx context.Context, req *requests.MonthPriceId) ([]*salemongo.CategoryMonthlyByIdRow, error)
	GetYearPriceById(ctx context.Context, req *requests.YearPriceId) ([]*salemongo.CategoryYearlyByIdRow, error)
}

// CategoryStatsRepository provides aggregate category price statistics.
type CategoryStatsRepository interface {
	GetMonthlyTotalPrice(ctx context.Context, req *requests.MonthTotalPrice) ([]*salemongo.CategoryMonthlyTotalPriceRow, error)
	GetYearlyTotalPrices(ctx context.Context, year int) ([]*salemongo.CategoryYearlyTotalPriceRow, error)
	GetMonthPrice(ctx context.Context, year int) ([]*salemongo.CategoryMonthlyRow, error)
	GetYearPrice(ctx context.Context, year int) ([]*salemongo.CategoryYearlyRow, error)
}

// CategoryStatsByMerchantRepository provides category price statistics grouped by merchant.
type CategoryStatsByMerchantRepository interface {
	GetMonthlyTotalPriceByMerchant(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*salemongo.CategoryMonthlyTotalPriceByMerchantRow, error)
	GetYearlyTotalPricesByMerchant(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*salemongo.CategoryYearlyTotalPriceByMerchantRow, error)
	GetMonthPriceByMerchant(ctx context.Context, req *requests.MonthPriceMerchant) ([]*salemongo.CategoryMonthlyByMerchantRow, error)
	GetYearPriceByMerchant(ctx context.Context, req *requests.YearPriceMerchant) ([]*salemongo.CategoryYearlyByMerchantRow, error)
}
