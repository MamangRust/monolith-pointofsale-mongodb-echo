package role

import (
	"context"

	pbuserroles "github.com/MamangRust/monolith-point-of-sale-pb/user_roles"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-pkg/resilience"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	userrole_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/user_role_errors"
)

// CommandRepository is the contract consumers depend on for user_roles writes.
type CommandRepository interface {
	AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*salemongo.UserRole, error)
	RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error
}

// UserRoleAdapter implements CommandRepository on top of the generated
// user_role service client. The user_roles collection stays owned by the role
// service; consumers reach it only through this adapter.
type UserRoleAdapter struct {
	client pbuserroles.UserRoleServiceClient
	guard  *resilience.DependencyGuard
}

// NewUserRoleAdapter builds a user_role adapter. Passing zero options leaves
// the guard nil, which makes DependencyGuard.Call a plain passthrough.
func NewUserRoleAdapter(client pbuserroles.UserRoleServiceClient, opts ...adapter.GuardOption) *UserRoleAdapter {
	r := &UserRoleAdapter{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *UserRoleAdapter) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *UserRoleAdapter) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*salemongo.UserRole, error) {
	var resp *pbuserroles.ApiResponseUserRole
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.client.AssignRoleToUser(ctx, &pbuserroles.AssignRoleToUserRequest{
			UserId: int32(req.UserId),
			RoleId: int32(req.RoleId),
		})
		return callErr
	})
	if err != nil || resp == nil || resp.Data == nil {
		return nil, userrole_errors.ErrAssignRoleToUser.WithInternal(err)
	}
	return &salemongo.UserRole{
		UserRoleID: resp.Data.UserRoleId,
		UserID:     resp.Data.UserId,
		RoleID:     resp.Data.RoleId,
	}, nil
}

func (r *UserRoleAdapter) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		_, callErr := r.client.RemoveRoleFromUser(ctx, &pbuserroles.RemoveRoleFromUserRequest{
			UserId: int32(req.UserId),
			RoleId: int32(req.RoleId),
		})
		return callErr
	})
	if err != nil {
		return userrole_errors.ErrRemoveRole.WithInternal(err)
	}
	return nil
}
