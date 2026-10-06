package service

import (
	"context"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-product/repository"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type ProductQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllProducts) ([]*repository.ProductResult, *int, error)
	FindByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest) ([]*repository.ProductByMerchantResult, *int, error)
	FindByCategory(ctx context.Context, req *requests.ProductByCategoryRequest) ([]*repository.ProductByCategoryResult, *int, error)
	FindById(ctx context.Context, productID int) (*salemongo.Product, error)
	FindByActive(ctx context.Context, req *requests.FindAllProducts) ([]*repository.ProductResultDeleteAt, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllProducts) ([]*repository.ProductResultDeleteAt, *int, error)
}

type ProductCommandService interface {
	CreateProduct(ctx context.Context, req *requests.CreateProductRequest) (*salemongo.Product, error)
	UpdateProduct(ctx context.Context, req *requests.UpdateProductRequest) (*salemongo.Product, error)
	DecrementProductCountStock(ctx context.Context, productID int, quantity int) (*salemongo.Product, error)
	IncrementProductCountStock(ctx context.Context, productID int, quantity int) (*salemongo.Product, error)
	TrashProduct(ctx context.Context, productID int) (*salemongo.Product, error)
	RestoreProduct(ctx context.Context, productID int) (*salemongo.Product, error)
	DeleteProductPermanent(ctx context.Context, productID int) (bool, error)
	RestoreAllProducts(ctx context.Context) (bool, error)
	DeleteAllProductsPermanent(ctx context.Context) (bool, error)
}
