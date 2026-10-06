package handler

import (
	"context"

	pbuserroles "github.com/MamangRust/monolith-point-of-sale-pb/user_roles"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-role/service"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/role_errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

type userRoleHandleGrpc struct {
	pbuserroles.UnimplementedUserRoleServiceServer
	userRole service.UserRoleService
	logger   logger.LoggerInterface
}

func NewUserRoleHandleGrpc(service *service.Service, logger logger.LoggerInterface) UserRoleHandleGrpc {
	return &userRoleHandleGrpc{
		userRole: service.UserRole,
		logger:   logger,
	}
}

func (s *userRoleHandleGrpc) AssignRoleToUser(ctx context.Context, reqPb *pbuserroles.AssignRoleToUserRequest) (*pbuserroles.ApiResponseUserRole, error) {
	userID := int(reqPb.GetUserId())
	roleID := int(reqPb.GetRoleId())
	if userID <= 0 || roleID <= 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	userRole, err := s.userRole.AssignRoleToUser(ctx, &requests.CreateUserRoleRequest{
		UserId: userID,
		RoleId: roleID,
	})
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pbuserroles.ApiResponseUserRole{
		Status:  "success",
		Message: "Successfully assigned role to user",
		Data:    mapResponseUserRole(userRole),
	}, nil
}

func (s *userRoleHandleGrpc) RemoveRoleFromUser(ctx context.Context, reqPb *pbuserroles.RemoveRoleFromUserRequest) (*emptypb.Empty, error) {
	userID := int(reqPb.GetUserId())
	roleID := int(reqPb.GetRoleId())
	if userID <= 0 || roleID <= 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	err := s.userRole.RemoveRoleFromUser(ctx, &requests.RemoveUserRoleRequest{
		UserId: userID,
		RoleId: roleID,
	})
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &emptypb.Empty{}, nil
}

func mapResponseUserRole(userRole *salemongo.UserRole) *pbuserroles.UserRoleResponse {
	if userRole == nil {
		return nil
	}
	return &pbuserroles.UserRoleResponse{
		UserRoleId: userRole.UserRoleID,
		UserId:     userRole.UserID,
		RoleId:     userRole.RoleID,
	}
}
