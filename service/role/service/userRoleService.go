package service

import (
	"context"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	mencache "github.com/MamangRust/monolith-point-of-sale-role/cache"
	"github.com/MamangRust/monolith-point-of-sale-role/repository"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-point-of-sale-shared/errorhandler"
	userrole_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/user_role_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type UserRoleService interface {
	AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*salemongo.UserRole, error)
	RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error
}

type userRoleService struct {
	mencache        mencache.RoleCommandCache
	userRoleCommand repository.UserRoleCommandRepository
	logger          logger.LoggerInterface
	observability   observability.TraceLoggerObservability
}

func NewUserRoleService(
	mencache mencache.RoleCommandCache,
	userRoleCommand repository.UserRoleCommandRepository,
	logger logger.LoggerInterface,
	obs observability.TraceLoggerObservability,
) UserRoleService {
	return &userRoleService{
		mencache:        mencache,
		userRoleCommand: userRoleCommand,
		logger:          logger,
		observability:   obs,
	}
}

func (s *userRoleService) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*salemongo.UserRole, error) {
	const method = "AssignRoleToUser"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user.id", req.UserId), attribute.Int("role.id", req.RoleId))
	defer func() { end(status) }()

	userRole, err := s.userRoleCommand.AssignRoleToUser(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*salemongo.UserRole](
			s.logger, userrole_errors.ErrFailedAssignRoleToUser.WithInternal(err), method, span, zap.Error(err))
	}
	s.mencache.DeleteCachedRoleAllCache(ctx)
	logSuccess("Successfully assigned role to user",
		zap.Int32("user.id", userRole.UserID), zap.Int32("role.id", userRole.RoleID), zap.Bool("success", true))
	return userRole, nil
}

func (s *userRoleService) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	const method = "RemoveRoleFromUser"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user.id", req.UserId), attribute.Int("role.id", req.RoleId))
	defer func() { end(status) }()

	err := s.userRoleCommand.RemoveRoleFromUser(ctx, req)
	if err != nil {
		status = "error"
		_, retErr := sharederrorhandler.HandleError[error](
			s.logger, userrole_errors.ErrFailedRemoveRole.WithInternal(err), method, span, zap.Error(err))
		return retErr
	}
	s.mencache.DeleteCachedRoleAllCache(ctx)
	logSuccess("Successfully removed role from user",
		zap.Int("user.id", req.UserId), zap.Int("role.id", req.RoleId), zap.Bool("success", true))
	return nil
}
