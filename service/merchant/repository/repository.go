package repository

import (
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	useradapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/user"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	MerchantCommand         MerchantCommandRepository
	MerchantQuery           MerchantQueryRepository
	MerchantDocumentCommand MerchantDocumentCommandRepository
	MerchantDocumentQuery   MerchantDocumentQueryRepository
	UserQuery               UserQueryRepository
}

// GuardOptions collects the dependency guards for each remote dependency. The
// options are applied to the adapter built inside NewRepositories, so every
// outbound gRPC call runs under a per-call deadline, bulkhead and circuit
// breaker.
type GuardOptions struct {
	User []adapter.GuardOption
}

// NewRepositories assembles the merchant repositories: Mongo-backed merchant
// and merchant-document repos plus the user gRPC adapter, constructed here
// under its guard.
func NewRepositories(db *mongo.Database, userClient pbusers.UserQueryServiceClient, guards ...GuardOptions) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		MerchantCommand:         NewMerchantCommandRepository(db.Collection("merchants")),
		MerchantQuery:           NewMerchantQueryRepository(db.Collection("merchants")),
		MerchantDocumentCommand: NewMerchantDocumentCommandRepository(db.Collection("merchant_documents")),
		MerchantDocumentQuery:   NewMerchantDocumentQueryRepository(db.Collection("merchant_documents")),
		// Merchant only reads users (FindById), so the command client stays nil.
		UserQuery: useradapter.New(userClient, nil, g.User...),
	}
}
