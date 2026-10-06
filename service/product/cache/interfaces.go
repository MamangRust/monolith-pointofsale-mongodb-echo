package mencache

import (
	"context"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-product/repository"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type ProductQueryCache interface {
	GetCachedProducts(ctx context.Context, req *requests.FindAllProducts) ([]*repository.ProductResult, *int, bool)
	SetCachedProducts(ctx context.Context, req *requests.FindAllProducts, data []*repository.ProductResult, total *int)

	GetCachedProductsByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest) ([]*repository.ProductByMerchantResult, *int, bool)
	SetCachedProductsByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest, data []*repository.ProductByMerchantResult, total *int)

	GetCachedProductsByCategory(ctx context.Context, req *requests.ProductByCategoryRequest) ([]*repository.ProductByCategoryResult, *int, bool)
	SetCachedProductsByCategory(ctx context.Context, req *requests.ProductByCategoryRequest, data []*repository.ProductByCategoryResult, total *int)

	GetCachedProductActive(ctx context.Context, req *requests.FindAllProducts) ([]*repository.ProductResultDeleteAt, *int, bool)
	SetCachedProductActive(ctx context.Context, req *requests.FindAllProducts, data []*repository.ProductResultDeleteAt, total *int)

	GetCachedProductTrashed(ctx context.Context, req *requests.FindAllProducts) ([]*repository.ProductResultDeleteAt, *int, bool)
	SetCachedProductTrashed(ctx context.Context, req *requests.FindAllProducts, data []*repository.ProductResultDeleteAt, total *int)

	GetCachedProduct(ctx context.Context, productID int) (*salemongo.Product, bool)
	SetCachedProduct(ctx context.Context, data *salemongo.Product)
}

type ProductCommandCache interface {
	DeleteCachedProduct(ctx context.Context, productID int)
	DeleteCachedProductAllCache(ctx context.Context)
}
