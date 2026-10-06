package repository

import (
	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	cashieradapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/cashier"
	merchantadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/merchant"
	orderadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/order"
	orderitemadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/orderitem"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	TransactionCommand         TransactionCommandRepository
	TransactionQuery           TransactionQueryRepository
	CashierQuery               CashierQueryRepository
	MerchantQuery              MerchantQueryRepository
	OrderQuery                 OrderQueryRepository
	OrderItemQuery             OrderItemQueryRepository
	TransactionStats           TransactionStatsRepository
	TransactionStatsByMerchant TransactionStatsByMerchantRepository
}

// GuardOptions collects the dependency guards for each remote dependency. The
// options are applied to the adapters built inside NewRepositories, so every
// outbound gRPC call runs under a per-call deadline, bulkhead and circuit
// breaker.
type GuardOptions struct {
	Cashier   []adapter.GuardOption
	Merchant  []adapter.GuardOption
	Order     []adapter.GuardOption
	OrderItem []adapter.GuardOption
}

// NewRepositories assembles the transaction repositories: Mongo-backed
// command/query/stats on the transactions collection plus the cashier, merchant,
// order and order-item gRPC adapters, each constructed here under its guard.
func NewRepositories(
	db *mongo.Database,
	cashierClient pbcashiers.CashierQueryServiceClient,
	merchantClient pbmerchants.MerchantQueryServiceClient,
	orderClient pborders.OrderQueryServiceClient,
	orderItemClient pborder_items.OrderItemQueryServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		TransactionCommand: NewTransactionCommandRepository(db.Collection("transactions")),
		TransactionQuery:   NewTransactionQueryRepository(db.Collection("transactions")),
		CashierQuery:       cashieradapter.New(cashierClient, g.Cashier...),
		MerchantQuery:      merchantadapter.New(merchantClient, g.Merchant...),
		OrderQuery:         orderadapter.New(orderClient, g.Order...),
		// Transaction only reads order items, so the command client stays nil.
		OrderItemQuery:             orderitemadapter.New(orderItemClient, nil, g.OrderItem...),
		TransactionStats:           NewTransactionStatsRepository(db),
		TransactionStatsByMerchant: NewTransactionStatsByMerchantRepository(db),
	}
}
