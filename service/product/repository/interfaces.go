package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

// ProductResult is the result type for paginated product queries.
type ProductResult struct {
	ProductID    int32
	MerchantID   int32
	CategoryID   int32
	Name         string
	Description  *string
	Price        int32
	CountInStock int32
	Brand        *string
	Weight       *int32
	SlugProduct  *string
	ImageProduct *string
	Barcode      *string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	TotalCount   int64
}

// ProductResultDeleteAt is the result type for paginated product queries with deleted_at.
type ProductResultDeleteAt struct {
	ProductID    int32
	MerchantID   int32
	CategoryID   int32
	Name         string
	Description  *string
	Price        int32
	CountInStock int32
	Brand        *string
	Weight       *int32
	SlugProduct  *string
	ImageProduct *string
	Barcode      *string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	DeletedAt    *time.Time
	TotalCount   int64
}

// ProductByMerchantResult is the result type for product-by-merchant queries.
type ProductByMerchantResult struct {
	TotalCount   int64
	ProductID    int32
	Name         string
	Description  *string
	Price        int32
	CountInStock int32
	Brand        *string
	ImageProduct *string
	Barcode      *string
	CreatedAt    *time.Time
	CategoryName string
}

// ProductByCategoryResult is the result type for product-by-category queries.
type ProductByCategoryResult struct {
	TotalCount   int64
	ProductID    int32
	MerchantID   int32
	CategoryID   int32
	SlugProduct  *string
	Weight       *int32
	Name         string
	Description  *string
	Price        int32
	CountInStock int32
	Brand        *string
	ImageProduct *string
	Barcode      *string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

// The category/merchant references a product carries live in other bounded
// contexts and are reached through the shared gRPC adapters in
// pkg/adapter/{category,merchant} (see repository.go).

type ProductQueryRepository interface {
	FindAllProducts(ctx context.Context, req *requests.FindAllProducts) ([]*ProductResult, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllProducts) ([]*ProductResultDeleteAt, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllProducts) ([]*ProductResultDeleteAt, *int, error)
	FindByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest) ([]*ProductByMerchantResult, *int, error)
	// FindByCategory takes the already-resolved category id: the name→id
	// lookup happens in the service layer through the category adapter.
	FindByCategory(ctx context.Context, categoryID int32, req *requests.ProductByCategoryRequest) ([]*ProductByCategoryResult, *int, error)
	FindById(ctx context.Context, product_id int) (*salemongo.Product, error)
	FindByIdTrashed(ctx context.Context, id int) (*salemongo.Product, error)
}

type ProductCommandRepository interface {
	CreateProduct(ctx context.Context, request *requests.CreateProductRequest) (*salemongo.Product, error)
	UpdateProduct(ctx context.Context, request *requests.UpdateProductRequest) (*salemongo.Product, error)
	UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*salemongo.Product, error)
	DecrementProductCountStock(ctx context.Context, productID int, quantity int) (*salemongo.Product, error)
	IncrementProductCountStock(ctx context.Context, productID int, quantity int) (*salemongo.Product, error)
	TrashedProduct(ctx context.Context, product_id int) (*salemongo.Product, error)
	RestoreProduct(ctx context.Context, product_id int) (*salemongo.Product, error)
	DeleteProductPermanent(ctx context.Context, product_id int) (bool, error)
	RestoreAllProducts(ctx context.Context) (bool, error)
	DeleteAllProductPermanent(ctx context.Context) (bool, error)
}

// sqlc row type aliases for handler compatibility
type ProductRow = ProductResult
type ProductActiveRow = ProductResultDeleteAt
type ProductTrashedRow = ProductResultDeleteAt
type ProductCategoryNameRow = ProductByCategoryResult
type ProductMerchantRow = ProductByMerchantResult
