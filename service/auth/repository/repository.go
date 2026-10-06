package repository

import (
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	pbuserroles "github.com/MamangRust/monolith-point-of-sale-pb/user_roles"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	roleadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/role"
	useradapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/user"
	"go.mongodb.org/mongo-driver/mongo"
)

// GuardOptions collects the dependency guards applied to the identity adapters
// built inside NewRepositories. The role query and user_role command sides get
// separate breakers so one failure mode cannot starve the other.
type GuardOptions struct {
	User     []adapter.GuardOption
	Role     []adapter.GuardOption
	UserRole []adapter.GuardOption
}

// Compile-time assertions: the shared adapters must satisfy auth's own
// repository contracts — auth keeps the interfaces, pkg/adapter keeps the
// implementations.
var (
	_ UserRepository     = (*useradapter.Repository)(nil)
	_ RoleRepository     = (*roleadapter.RoleAdapter)(nil)
	_ UserRoleRepository = (*roleadapter.UserRoleAdapter)(nil)
)

// NewRepositories assembles auth's own repositories (refresh/reset tokens,
// outbox) and builds the shared user/role gRPC adapters from the raw clients
// under their guards. Users, roles and user_roles live in other services, so
// auth never touches those collections directly.
func NewRepositories(
	db *mongo.Database,
	userQueryClient pbusers.UserQueryServiceClient,
	userCommandClient pbusers.UserCommandServiceClient,
	roleQueryClient pbroles.RoleQueryServiceClient,
	userRoleClient pbuserroles.UserRoleServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		User:         useradapter.New(userQueryClient, userCommandClient, g.User...),
		RefreshToken: NewRefreshTokenRepository(db.Collection("refresh_tokens")),
		UserRole:     roleadapter.NewUserRoleAdapter(userRoleClient, g.UserRole...),
		Role:         roleadapter.NewRoleAdapter(roleQueryClient, g.Role...),
		ResetToken:   NewResetTokenRepository(db.Collection("reset_tokens")),
		Outbox:       NewOutboxRepository(),
	}
}
