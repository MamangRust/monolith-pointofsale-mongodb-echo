// Package role adapts the Role service gRPC API into the shared domain model.
// It exposes two narrow adapters over the generated role clients: RoleAdapter
// for role reads and UserRoleAdapter for user_roles writes.
package role

import (
	"context"

	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter/convert"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-pkg/resilience"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/role_errors"
)

// QueryRepository is the contract consumers depend on for role reads.
type QueryRepository interface {
	FindById(ctx context.Context, roleID int) (*salemongo.Role, error)
	FindByName(ctx context.Context, name string) (*salemongo.Role, error)
}

// RoleAdapter implements QueryRepository on top of the generated role query
// client.
type RoleAdapter struct {
	client pbroles.RoleQueryServiceClient
	guard  *resilience.DependencyGuard
}

// NewRoleAdapter builds a role adapter. Passing zero options leaves the guard
// nil, which makes DependencyGuard.Call a plain passthrough.
func NewRoleAdapter(client pbroles.RoleQueryServiceClient, opts ...adapter.GuardOption) *RoleAdapter {
	r := &RoleAdapter{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *RoleAdapter) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *RoleAdapter) FindById(ctx context.Context, roleID int) (*salemongo.Role, error) {
	var resp *pbroles.ApiResponseRole
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.client.FindByIdRole(ctx, &pbroles.FindByIdRoleRequest{RoleId: int32(roleID)})
		return callErr
	})
	if err != nil || resp == nil || resp.Data == nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	return roleToModel(resp.Data), nil
}

func (r *RoleAdapter) FindByName(ctx context.Context, name string) (*salemongo.Role, error) {
	var resp *pbroles.ApiResponseRole
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.client.FindByNameRole(ctx, &pbroles.FindByNameRoleRequest{Name: name})
		return callErr
	})
	if err != nil || resp == nil || resp.Data == nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	return roleToModel(resp.Data), nil
}

func roleToModel(role *pbroles.RoleResponse) *salemongo.Role {
	if role == nil {
		return nil
	}
	return &salemongo.Role{
		RoleID:    role.Id,
		RoleName:  role.Name,
		CreatedAt: convert.TimePtr(role.CreatedAt),
		UpdatedAt: convert.TimePtr(role.UpdatedAt),
	}
}
