package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type RoleResult struct {
	RoleID     int32
	RoleName   string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DeletedAt  *time.Time
	TotalCount int64
}

type RoleQueryRepository interface {
	FindAllRoles(ctx context.Context, req *requests.FindAllRoles) ([]*RoleResult, error)
	FindById(ctx context.Context, roleID int) (*salemongo.Role, error)
	FindByName(ctx context.Context, name string) (*salemongo.Role, error)
	FindByUserId(ctx context.Context, userID int) ([]*salemongo.Role, error)
	FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*RoleResult, error)
	FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*RoleResult, error)
}

type RoleCommandRepository interface {
	CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*salemongo.Role, error)
	UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*salemongo.Role, error)
	TrashedRole(ctx context.Context, roleID int) (*salemongo.Role, error)
	RestoreRole(ctx context.Context, roleID int) (*salemongo.Role, error)
	DeleteRolePermanent(ctx context.Context, roleID int) (bool, error)
	RestoreAllRole(ctx context.Context) (bool, error)
	DeleteAllRolePermanent(ctx context.Context) (bool, error)
}

type UserRoleCommandRepository interface {
	AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*salemongo.UserRole, error)
	RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error
}

// sqlc row type aliases for handler compatibility
type RoleRow = RoleResult
type RoleActiveRow = RoleResult
type RoleTrashedRow = RoleResult
