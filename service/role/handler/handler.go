package handler

import (
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	pbuserroles "github.com/MamangRust/monolith-point-of-sale-pb/user_roles"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-role/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Role        pbroles.RoleQueryServiceServer
	RoleCommand pbroles.RoleCommandServiceServer
	UserRole    pbuserroles.UserRoleServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Role:        NewRoleQueryHandleGrpc(deps.Service, deps.Logger),
		RoleCommand: NewRoleCommandHandleGrpc(deps.Service, deps.Logger),
		UserRole:    NewUserRoleHandleGrpc(deps.Service, deps.Logger),
	}
}
