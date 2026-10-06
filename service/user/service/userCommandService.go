package service

import (
	"context"

	roleadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/role"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-pkg/hash"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-point-of-sale-shared/errorhandler"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/role_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/user_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	mencache "github.com/MamangRust/monolith-point-of-sale-user/cache"
	"github.com/MamangRust/monolith-point-of-sale-user/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type userCommandDeps struct {
	Cache           mencache.UserCommandCache
	UserQuery       repository.UserQueryRepository
	UserCommand     repository.UserCommandRepository
	RoleQuery       roleadapter.QueryRepository
	UserRoleAdapter  roleadapter.CommandRepository
	Logger          logger.LoggerInterface
	Hashing         hash.HashPassword
	Observability   observability.TraceLoggerObservability
}

type userCommandService struct {
	mencache              mencache.UserCommandCache
	userQueryRepository   repository.UserQueryRepository
	userCommandRepository repository.UserCommandRepository
	roleAdapter           roleadapter.QueryRepository
	userRoleAdapter       roleadapter.CommandRepository
	logger                logger.LoggerInterface
	hashing               hash.HashPassword
	observability         observability.TraceLoggerObservability
}

func NewUserCommandService(params *userCommandDeps) UserCommandService {
	return &userCommandService{
		mencache:              params.Cache,
		userQueryRepository:   params.UserQuery,
		userCommandRepository: params.UserCommand,
		roleAdapter:           params.RoleQuery,
		userRoleAdapter:       params.UserRoleAdapter,
		logger:                params.Logger,
		hashing:               params.Hashing,
		observability:         params.Observability,
	}
}

func (s *userCommandService) CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*salemongo.User, error) {
	const method = "CreateUser"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() { end(status) }()

	existingUser, err := s.userQueryRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*salemongo.User](s.logger, user_errors.ErrInternalServerError.WithInternal(err), method, span, zap.String("email", request.Email))
	}
	if existingUser != nil {
		status = "error"
		return sharederrorhandler.HandleError[*salemongo.User](s.logger, user_errors.ErrUserEmailAlready, method, span, zap.String("email", request.Email))
	}

	hashedPassword, err := s.hashing.HashPassword(request.Password)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*salemongo.User](s.logger, user_errors.ErrUserPassword.WithInternal(err), method, span, zap.String("email", request.Email))
	}
	request.Password = hashedPassword

	const defaultRoleName = "ROLE_ADMIN"
	role, err := s.roleAdapter.FindByName(ctx, defaultRoleName)
	if err != nil || role == nil {
		status = "error"
		if err == nil {
			err = role_errors.ErrRoleNotFoundRes
		}
		return sharederrorhandler.HandleError[*salemongo.User](s.logger, role_errors.ErrRoleNotFoundRes.WithInternal(err), method, span, zap.String("name", defaultRoleName))
	}

	res, err := s.userCommandRepository.CreateUser(ctx, request)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*salemongo.User](s.logger, user_errors.ErrFailedCreateUser.WithInternal(err), method, span, zap.String("email", request.Email))
	}

	_, err = s.userRoleAdapter.AssignRoleToUser(ctx, &requests.CreateUserRoleRequest{
		UserId: int(res.UserID),
		RoleId: int(role.RoleID),
	})
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*salemongo.User](s.logger, role_errors.ErrRoleNotFoundRes.WithInternal(err), method, span, zap.Int32("user.id", res.UserID))
	}

	logSuccess("Successfully created user", zap.Int32("user.id", res.UserID))
	return res, nil
}

func (s *userCommandService) UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*salemongo.User, error) {
	const method = "UpdateUser"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() { end(status) }()

	existingUser, err := s.userQueryRepository.FindById(ctx, *request.UserID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*salemongo.User](s.logger, err, method, span, zap.Int("user.id", *request.UserID))
	}

	if request.Email != "" && request.Email != existingUser.Email {
		duplicateUser, _ := s.userQueryRepository.FindByEmail(ctx, request.Email)
		if duplicateUser != nil {
			status = "error"
			return sharederrorhandler.HandleError[*salemongo.User](s.logger, user_errors.ErrUserEmailAlready, method, span, zap.String("email", request.Email))
		}
	}

	if request.Password != "" {
		hashedPassword, err := s.hashing.HashPassword(request.Password)
		if err != nil {
			status = "error"
			return sharederrorhandler.HandleError[*salemongo.User](s.logger, user_errors.ErrUserPassword.WithInternal(err), method, span)
		}
		request.Password = hashedPassword
	}

	const defaultRoleName = "ROLE_ADMIN"
	role, err := s.roleAdapter.FindByName(ctx, defaultRoleName)
	if err != nil || role == nil {
		status = "error"
		if err == nil {
			err = role_errors.ErrRoleNotFoundRes
		}
		return sharederrorhandler.HandleError[*salemongo.User](s.logger, role_errors.ErrRoleNotFoundRes.WithInternal(err), method, span, zap.String("name", defaultRoleName))
	}

	res, err := s.userCommandRepository.UpdateUser(ctx, request)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*salemongo.User](s.logger, user_errors.ErrFailedUpdateUser.WithInternal(err), method, span)
	}

	s.mencache.DeleteUserCache(ctx, int(res.UserID))
	logSuccess("Successfully updated user", zap.Int32("user.id", res.UserID))
	return res, nil
}

func (s *userCommandService) TrashedUser(ctx context.Context, userID int) (*salemongo.User, error) {
	const method = "TrashedUser"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() { end(status) }()

	res, err := s.userCommandRepository.TrashedUser(ctx, userID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*salemongo.User](s.logger, user_errors.ErrFailedTrashedUser.WithInternal(err), method, span)
	}
	s.mencache.DeleteUserCache(ctx, userID)
	logSuccess("Successfully trashed user", zap.Int32("user.id", res.UserID))
	return res, nil
}

func (s *userCommandService) RestoreUser(ctx context.Context, userID int) (*salemongo.User, error) {
	const method = "RestoreUser"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() { end(status) }()

	res, err := s.userCommandRepository.RestoreUser(ctx, userID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*salemongo.User](s.logger, user_errors.ErrFailedRestoreUser.WithInternal(err), method, span)
	}
	s.mencache.DeleteUserCache(ctx, userID)
	logSuccess("Successfully restored user", zap.Int32("user.id", res.UserID))
	return res, nil
}

func (s *userCommandService) DeleteUserPermanent(ctx context.Context, userID int) (bool, error) {
	const method = "DeleteUserPermanent"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() { end(status) }()

	success, err := s.userCommandRepository.DeleteUserPermanent(ctx, userID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](s.logger, user_errors.ErrFailedDeletePermanent.WithInternal(err), method, span)
	}
	s.mencache.DeleteUserCache(ctx, userID)
	logSuccess("Successfully deleted user permanently", zap.Int("user.id", userID))
	return success, nil
}

func (s *userCommandService) RestoreAllUser(ctx context.Context) (bool, error) {
	const method = "RestoreAllUser"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() { end(status) }()

	success, err := s.userCommandRepository.RestoreAllUser(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](s.logger, user_errors.ErrFailedRestoreAll.WithInternal(err), method, span)
	}
	s.mencache.DeleteUserAllCache(ctx)
	logSuccess("Successfully restored all users")
	return success, nil
}

func (s *userCommandService) DeleteAllUserPermanent(ctx context.Context) (bool, error) {
	const method = "DeleteAllUserPermanent"
	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)
	defer func() { end(status) }()

	success, err := s.userCommandRepository.DeleteAllUserPermanent(ctx)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[bool](s.logger, user_errors.ErrFailedDeleteAll.WithInternal(err), method, span)
	}
	s.mencache.DeleteUserAllCache(ctx)
	logSuccess("Successfully deleted all users permanently")
	return success, nil
}

// CreateUserRecord stores a registration sent by auth-service. The password
// arrives already hashed (auth owns hashing), so unlike CreateUser this path
// must not hash again — doing so would make the stored hash unusable at login.
func (s *userCommandService) CreateUserRecord(ctx context.Context, request *requests.RegisterRequest) (*salemongo.User, error) {
	const method = "CreateUserRecord"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("email", request.Email))
	defer func() { end(status) }()

	existingUser, err := s.userQueryRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*salemongo.User](s.logger, user_errors.ErrInternalServerError.WithInternal(err), method, span, zap.String("email", request.Email), zap.Error(err))
	}
	if existingUser != nil {
		status = "error"
		return sharederrorhandler.HandleError[*salemongo.User](s.logger, user_errors.ErrUserEmailAlready, method, span, zap.String("email", request.Email))
	}

	user, err := s.userCommandRepository.CreateUserRecord(ctx, request)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*salemongo.User](s.logger, user_errors.ErrFailedCreateUser.WithInternal(err), method, span, zap.String("email", request.Email), zap.Error(err))
	}

	s.mencache.DeleteUserAllCache(ctx)
	logSuccess("Successfully created user record", zap.Int32("user.id", user.UserID))
	return user, nil
}

func (s *userCommandService) UpdateUserIsVerified(ctx context.Context, user_id int, is_verified bool) (*salemongo.User, error) {
	const method = "UpdateUserIsVerified"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user.id", user_id))
	defer func() { end(status) }()

	user, err := s.userCommandRepository.UpdateUserIsVerified(ctx, user_id, is_verified)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*salemongo.User](s.logger, user_errors.ErrFailedUpdateUser.WithInternal(err), method, span, zap.Int("user.id", user_id), zap.Error(err))
	}

	s.mencache.DeleteUserCache(ctx, user_id)
	logSuccess("Successfully updated user verification", zap.Int("user.id", user_id))
	return user, nil
}

// UpdateUserPassword stores a password hash supplied by auth-service; it does
// not hash the value itself.
func (s *userCommandService) UpdateUserPassword(ctx context.Context, user_id int, password string) (*salemongo.User, error) {
	const method = "UpdateUserPassword"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("user.id", user_id))
	defer func() { end(status) }()

	user, err := s.userCommandRepository.UpdateUserPassword(ctx, user_id, password)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*salemongo.User](s.logger, user_errors.ErrUpdateUserPassword.WithInternal(err), method, span, zap.Int("user.id", user_id), zap.Error(err))
	}

	s.mencache.DeleteUserCache(ctx, user_id)
	logSuccess("Successfully updated user password", zap.Int("user.id", user_id))
	return user, nil
}
