package repository

import (
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	pbuserroles "github.com/MamangRust/monolith-point-of-sale-pb/user_roles"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	roleadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/role"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	UserCommand     UserCommandRepository
	UserQuery       UserQueryRepository
	RoleQuery       roleadapter.QueryRepository
	UserRoleAdapter *roleadapter.UserRoleAdapter
}

type GuardOptions struct {
	Role []adapter.GuardOption
}

func NewRepositories(db *mongo.Database,
	roleQueryClient pbroles.RoleQueryServiceClient,
	userRoleClient pbuserroles.UserRoleServiceClient,
	guards ...GuardOptions) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	roleAdapter := roleadapter.NewRoleAdapter(roleQueryClient, g.Role...)
	userRoleAdapter := roleadapter.NewUserRoleAdapter(userRoleClient, g.Role...)

	return &Repositories{
		UserCommand:     NewUserCommandRepository(db.Collection("users")),
		UserQuery:       NewUserQueryRepository(db.Collection("users")),
		RoleQuery:       roleAdapter,
		UserRoleAdapter: userRoleAdapter,
	}
}
