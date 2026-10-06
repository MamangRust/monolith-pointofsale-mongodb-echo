package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type UserResult struct {
	UserID     int32
	Firstname  string
	Lastname   string
	Email      string
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DeletedAt  *time.Time
	TotalCount int64
}

type UserQueryRepository interface {
	FindAllUsers(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error)
	FindById(ctx context.Context, user_id int) (*salemongo.User, error)
	FindByEmail(ctx context.Context, email string) (*salemongo.User, error)
	FindByEmailAndVerify(ctx context.Context, email string) (*salemongo.User, error)
	FindByVerificationCode(ctx context.Context, verification_code string) (*salemongo.User, error)
	FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error)
}

type UserCommandRepository interface {
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

// sqlc row type aliases for handler compatibility
type UserRow = UserResult
type UserActiveRow = UserResult
type UserTrashedRow = UserResult
