package service

import (
	"context"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-user/repository"
)

type UserQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, error)
	FindByID(ctx context.Context, id int) (*salemongo.User, error)
	FindByEmail(ctx context.Context, email string) (*salemongo.User, error)
	FindByEmailAndVerify(ctx context.Context, email string) (*salemongo.User, error)
	FindByVerificationCode(ctx context.Context, verificationCode string) (*salemongo.User, error)
	FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, error)
}

type UserCommandService interface {
	CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*salemongo.User, error)
	CreateUserRecord(ctx context.Context, request *requests.RegisterRequest) (*salemongo.User, error)
	UpdateUserIsVerified(ctx context.Context, user_id int, is_verified bool) (*salemongo.User, error)
	UpdateUserPassword(ctx context.Context, user_id int, password string) (*salemongo.User, error)
	UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*salemongo.User, error)
	TrashedUser(ctx context.Context, user_id int) (*salemongo.User, error)
	RestoreUser(ctx context.Context, user_id int) (*salemongo.User, error)
	DeleteUserPermanent(ctx context.Context, user_id int) (bool, error)
	RestoreAllUser(ctx context.Context) (bool, error)
	DeleteAllUserPermanent(ctx context.Context) (bool, error)
}
