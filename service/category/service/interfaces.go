package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-category/repository"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type CategoryQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResult, *int, error)
	FindById(ctx context.Context, category_id int) (*salemongo.Category, error)
	FindByName(ctx context.Context, name string) (*salemongo.Category, error)
	FindByIds(ctx context.Context, ids []int) ([]*salemongo.Category, error)
	FindByActive(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResultDeleteAt, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*repository.CategoryResultDeleteAt, *int, error)
}

type CategoryCommandService interface {
	CreateCategory(ctx context.Context, req *requests.CreateCategoryRequest) (*salemongo.Category, error)
	UpdateCategory(ctx context.Context, req *requests.UpdateCategoryRequest) (*salemongo.Category, error)
	TrashedCategory(ctx context.Context, category_id int) (*salemongo.Category, error)
	RestoreCategory(ctx context.Context, categoryID int) (*salemongo.Category, error)
	DeleteCategoryPermanent(ctx context.Context, categoryID int) (bool, error)
	RestoreAllCategories(ctx context.Context) (bool, error)
	DeleteAllCategoriesPermanent(ctx context.Context) (bool, error)
}

// CategoryStatsService exposes aggregate category price statistics.
type CategoryStatsService interface {
	FindMonthlyTotalPrice(ctx context.Context, req *requests.MonthTotalPrice) ([]*salemongo.CategoryMonthlyTotalPriceRow, error)
	FindYearlyTotalPrice(ctx context.Context, year int) ([]*salemongo.CategoryYearlyTotalPriceRow, error)
	FindMonthPrice(ctx context.Context, year int) ([]*salemongo.CategoryMonthlyRow, error)
	FindYearPrice(ctx context.Context, year int) ([]*salemongo.CategoryYearlyRow, error)
}

// CategoryStatsByIdService exposes per-category price statistics.
type CategoryStatsByIdService interface {
	FindMonthlyTotalPriceById(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*salemongo.CategoryMonthlyTotalPriceByIdRow, error)
	FindYearlyTotalPriceById(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*salemongo.CategoryYearlyTotalPriceByIdRow, error)
	FindMonthPriceById(ctx context.Context, req *requests.MonthPriceId) ([]*salemongo.CategoryMonthlyByIdRow, error)
	FindYearPriceById(ctx context.Context, req *requests.YearPriceId) ([]*salemongo.CategoryYearlyByIdRow, error)
}

// CategoryStatsByMerchantQueryService exposes per-merchant category price statistics.
type CategoryStatsByMerchantQueryService interface {
	FindMonthlyTotalPriceByMerchant(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*salemongo.CategoryMonthlyTotalPriceByMerchantRow, error)
	FindYearlyTotalPriceByMerchant(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*salemongo.CategoryYearlyTotalPriceByMerchantRow, error)
	FindMonthPriceByMerchant(ctx context.Context, req *requests.MonthPriceMerchant) ([]*salemongo.CategoryMonthlyByMerchantRow, error)
	FindYearPriceByMerchant(ctx context.Context, req *requests.YearPriceMerchant) ([]*salemongo.CategoryYearlyByMerchantRow, error)
}
