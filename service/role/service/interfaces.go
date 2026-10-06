package service

import (
	"context"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-role/repository"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type RoleQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllRoles) ([]*repository.RoleResult, *int, error)
	FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*repository.RoleResult, *int, error)
	FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*repository.RoleResult, *int, error)
	FindById(ctx context.Context, roleID int) (*salemongo.Role, error)
	FindByName(ctx context.Context, name string) (*salemongo.Role, error)
	FindByUserId(ctx context.Context, id int) ([]*salemongo.Role, error)
}

type RoleCommandService interface {
	CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*salemongo.Role, error)
	UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*salemongo.Role, error)
	TrashedRole(ctx context.Context, roleID int) (*salemongo.Role, error)
	RestoreRole(ctx context.Context, roleID int) (*salemongo.Role, error)
	DeleteRolePermanent(ctx context.Context, roleID int) (bool, error)
	RestoreAllRole(ctx context.Context) (bool, error)
	DeleteAllRolePermanent(ctx context.Context) (bool, error)
}
