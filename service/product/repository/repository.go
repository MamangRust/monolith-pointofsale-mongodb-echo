package repository

import (
	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	categoryadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/category"
	merchantadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/merchant"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	ProductCommand ProductCommandRepository
	ProductQuery   ProductQueryRepository
	CategoryQuery  categoryadapter.Repository
	MerchantQuery  merchantadapter.Repository
}

// GuardOptions collects the dependency guards for each remote dependency. The
// options are applied to the adapters built inside NewRepositories, so every
// outbound gRPC call runs under a per-call deadline, bulkhead and circuit
// breaker.
type GuardOptions struct {
	Category []adapter.GuardOption
	Merchant []adapter.GuardOption
}

// NewRepositories assembles the product-owned collection with the shared
// category/merchant gRPC adapters for the references a product carries.
func NewRepositories(
	db *mongo.Database,
	categoryClient pbcategories.CategoryQueryServiceClient,
	merchantClient pbmerchants.MerchantQueryServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		ProductCommand: NewProductCommandRepository(db.Collection("products")),
		ProductQuery:   NewProductQueryRepository(db.Collection("products")),
		CategoryQuery:  categoryadapter.New(categoryClient, g.Category...),
		MerchantQuery:  merchantadapter.New(merchantClient, g.Merchant...),
	}
}
