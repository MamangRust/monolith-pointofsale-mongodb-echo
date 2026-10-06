package handler

import (
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	pbuserroles "github.com/MamangRust/monolith-point-of-sale-pb/user_roles"
)

type RoleQueryHandleGrpc interface {
	pbroles.RoleQueryServiceServer
}

type RoleCommandHandleGrpc interface {
	pbroles.RoleCommandServiceServer
}

type UserRoleHandleGrpc interface {
	pbuserroles.UserRoleServiceServer
}
