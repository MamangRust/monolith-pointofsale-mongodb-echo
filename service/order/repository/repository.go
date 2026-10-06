package repository

import (
	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	cashieradapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/cashier"
	merchantadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/merchant"
	orderitemadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/orderitem"
	productadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/product"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	OrderCommand         OrderCommandRepository
	OrderQuery           OrderQueryRepository
	CashierQuery         cashieradapter.QueryRepository
	MerchantQuery        merchantadapter.Repository
	ProductQuery         productadapter.QueryRepository
	ProductCommand       productadapter.CommandRepository
	OrderItemQuery       orderitemadapter.QueryRepository
	OrderItemCommand     orderitemadapter.CommandRepository
	OrderStats           OrderStatsRepository
	OrderStatsByMerchant OrderStatByMerchantRepository
}

// GuardOptions collects the dependency guards for each remote dependency. The
// options are applied to the adapters built inside NewRepositories, so every
// outbound gRPC call runs under a per-call deadline, bulkhead and circuit
// breaker.
type GuardOptions struct {
	Cashier   []adapter.GuardOption
	Merchant  []adapter.GuardOption
	Product   []adapter.GuardOption
	OrderItem []adapter.GuardOption
}

// NewRepositories assembles the order service's own Mongo repositories with
// the shared gRPC adapters for every remote domain it depends on. Only the
// orders writes that belong to this service stay local; everything else is
// reached through pkg/adapter so there is exactly one implementation per
// remote domain.
func NewRepositories(
	db *mongo.Database,
	cashierClient pbcashiers.CashierQueryServiceClient,
	merchantClient pbmerchants.MerchantQueryServiceClient,
	productClient pbproducts.ProductQueryServiceClient,
	productCommandClient pbproducts.ProductCommandServiceClient,
	orderItemClient pborder_items.OrderItemQueryServiceClient,
	orderItemCommandClient pborder_items.OrderItemCommandServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	// One adapter instance serves both the query and the command role, so the
	// product/order_item clients are shared instead of guarded twice.
	productRepo := productadapter.New(productClient, productCommandClient, g.Product...)
	orderItemRepo := orderitemadapter.New(orderItemClient, orderItemCommandClient, g.OrderItem...)

	return &Repositories{
		OrderCommand:         NewOrderCommandRepository(db.Collection("orders")),
		OrderQuery:           NewOrderQueryRepository(db.Collection("orders")),
		CashierQuery:         cashieradapter.New(cashierClient, g.Cashier...),
		MerchantQuery:        merchantadapter.New(merchantClient, g.Merchant...),
		ProductQuery:         productRepo,
		ProductCommand:       productRepo,
		OrderItemQuery:       orderItemRepo,
		OrderItemCommand:     orderItemRepo,
		OrderStats:           NewOrderStatsRepository(db),
		OrderStatsByMerchant: NewOrderStatsByMerchantRepository(db),
	}
}
