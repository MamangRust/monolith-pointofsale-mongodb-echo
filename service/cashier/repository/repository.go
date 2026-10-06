package repository

import (
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	merchantadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/merchant"
	useradapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/user"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	CashierQuery           CashierQueryRepository
	CashierCommand         CashierCommandRepository
	MerchantQuery          merchantadapter.Repository
	UserQuery              useradapter.QueryRepository
	CashierStats           CashierStatsRepository
	CashierStatsById       CashierStatByIdRepository
	CashierStatsByMerchant CashierStatByMerchantRepository
}

// GuardOptions collects the dependency guards for each remote dependency. The
// options are applied to the adapters built inside NewRepositories, so every
// outbound gRPC call runs under a per-call deadline, bulkhead and circuit
// breaker.
type GuardOptions struct {
	Merchant []adapter.GuardOption
	User     []adapter.GuardOption
}

// NewRepositories assembles the cashier-owned collections with the shared
// merchant/user gRPC adapters used while validating a new cashier.
func NewRepositories(
	db *mongo.Database,
	merchantClient pbmerchants.MerchantQueryServiceClient,
	userClient pbusers.UserQueryServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		CashierQuery:           NewCashierQueryRepository(db.Collection("cashiers")),
		CashierCommand:         NewCashierCommandRepository(db.Collection("cashiers")),
		MerchantQuery:          merchantadapter.New(merchantClient, g.Merchant...),
		UserQuery:              useradapter.New(userClient, nil, g.User...),
		CashierStats:           NewCashierStatsRepository(db),
		CashierStatsById:       NewCashierStatsByIdRepository(db),
		CashierStatsByMerchant: NewCashierStatsByMerchantRepository(db),
	}
}
