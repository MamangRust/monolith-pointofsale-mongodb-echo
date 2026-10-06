// Package user adapts the User service gRPC API into the shared domain model.
// Query and command are exposed as one concrete Repository so a consumer that
// needs both (auth) uses a single instance, while query-only consumers
// (merchant, cashier) simply never reach the command methods.
package user

import (
	"context"

	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter/convert"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-pkg/resilience"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/user_errors"
)

// QueryRepository is the contract consumers depend on for user reads.
type QueryRepository interface {
	FindById(ctx context.Context, id int) (*salemongo.User, error)
	FindByEmail(ctx context.Context, email string) (*salemongo.User, error)
	FindByEmailAndVerify(ctx context.Context, email string) (*salemongo.User, error)
	FindByVerificationCode(ctx context.Context, code string) (*salemongo.User, error)
}

// CommandRepository is the contract consumers depend on for user writes on
// behalf of the auth service.
type CommandRepository interface {
	CreateUser(ctx context.Context, request *requests.RegisterRequest) (*salemongo.User, error)
	UpdateUserIsVerified(ctx context.Context, userID int, isVerified bool) (*salemongo.User, error)
	UpdateUserPassword(ctx context.Context, userID int, password string) (*salemongo.User, error)
}

// Repository implements QueryRepository and CommandRepository over the
// generated user query/command clients.
type Repository struct {
	query   pbusers.UserQueryServiceClient
	command pbusers.UserCommandServiceClient
	guard   *resilience.DependencyGuard
}

// New builds a user adapter. Passing zero options leaves the guard nil, which
// makes DependencyGuard.Call a plain passthrough. Query-only consumers may
// pass nil for the command client; command methods then fail cleanly instead
// of panicking.
func New(query pbusers.UserQueryServiceClient, command pbusers.UserCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{query: query, command: command}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *Repository) FindById(ctx context.Context, id int) (*salemongo.User, error) {
	var resp *pbusers.ApiResponseUser
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.query.FindById(ctx, &pbusers.FindByIdUserRequest{Id: int32(id)})
		return callErr
	})
	if err != nil || resp == nil || resp.Data == nil {
		return nil, sharedErrors.ErrInternal
	}

	u := resp.Data
	return &salemongo.User{
		UserID:    u.Id,
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Email:     u.Email,
		CreatedAt: convert.TimePtr(u.CreatedAt),
		UpdatedAt: convert.TimePtr(u.UpdatedAt),
	}, nil
}

// FindByEmail returns the user with the password hash. A missing user comes
// back as ErrUserNotFound so consumers can distinguish it from transport
// failures via errors.Is.
func (r *Repository) FindByEmail(ctx context.Context, email string) (*salemongo.User, error) {
	var resp *pbusers.ApiResponseUserWithPassword
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.query.FindByEmail(ctx, &pbusers.FindByEmailUserRequest{Email: email})
		return callErr
	})
	if err != nil || resp == nil || resp.Data == nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}
	return userWithPasswordToModel(resp.Data), nil
}

// FindByEmailAndVerify mirrors the legacy auth repo behaviour: lookup by email
// restricted to verified users, used for credential verification on login.
func (r *Repository) FindByEmailAndVerify(ctx context.Context, email string) (*salemongo.User, error) {
	var resp *pbusers.ApiResponseUserWithPassword
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.query.FindByEmailAndVerify(ctx, &pbusers.FindByEmailUserRequest{Email: email})
		return callErr
	})
	if err != nil || resp == nil || resp.Data == nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}
	return userWithPasswordToModel(resp.Data), nil
}

func (r *Repository) FindByVerificationCode(ctx context.Context, code string) (*salemongo.User, error) {
	var resp *pbusers.ApiResponseUser
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.query.FindByVerificationCode(ctx, &pbusers.FindByVerificationCodeUserRequest{VerificationCode: code})
		return callErr
	})
	if err != nil || resp == nil || resp.Data == nil {
		return nil, user_errors.ErrUserNotFound.WithInternal(err)
	}
	return userToModel(resp.Data), nil
}

// CreateUser registers a user through the User service. The password must
// already be hashed by the caller; the verification fields pass through as
// auth-managed values.
func (r *Repository) CreateUser(ctx context.Context, request *requests.RegisterRequest) (*salemongo.User, error) {
	if r.command == nil {
		return nil, user_errors.ErrCreateUser
	}

	var resp *pbusers.ApiResponseUser
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.command.CreateUserRecord(ctx, &pbusers.CreateUserRecordRequest{
			Firstname:        request.FirstName,
			Lastname:         request.LastName,
			Email:            request.Email,
			Password:         request.Password,
			VerificationCode: request.VerifiedCode,
			IsVerified:       request.IsVerified,
		})
		return callErr
	})
	if err != nil || resp == nil || resp.Data == nil {
		return nil, user_errors.ErrCreateUser.WithInternal(err)
	}
	return userToModel(resp.Data), nil
}

func (r *Repository) UpdateUserIsVerified(ctx context.Context, userID int, isVerified bool) (*salemongo.User, error) {
	if r.command == nil {
		return nil, user_errors.ErrUpdateUserVerificationCode
	}

	var resp *pbusers.ApiResponseUser
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.command.UpdateUserIsVerified(ctx, &pbusers.UpdateUserIsVerifiedUserRequest{
			Id:         int32(userID),
			IsVerified: isVerified,
		})
		return callErr
	})
	if err != nil || resp == nil || resp.Data == nil {
		return nil, user_errors.ErrUpdateUserVerificationCode.WithInternal(err)
	}
	return userToModel(resp.Data), nil
}

func (r *Repository) UpdateUserPassword(ctx context.Context, userID int, password string) (*salemongo.User, error) {
	if r.command == nil {
		return nil, user_errors.ErrUpdateUserPassword
	}

	var resp *pbusers.ApiResponseUser
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.command.UpdateUserPassword(ctx, &pbusers.UpdateUserPasswordUserRequest{
			Id:       int32(userID),
			Password: password,
		})
		return callErr
	})
	if err != nil || resp == nil || resp.Data == nil {
		return nil, user_errors.ErrUpdateUserPassword.WithInternal(err)
	}
	return userToModel(resp.Data), nil
}

func userToModel(u *pbusers.UserResponse) *salemongo.User {
	if u == nil {
		return nil
	}
	return &salemongo.User{
		UserID:    u.Id,
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Email:     u.Email,
		CreatedAt: convert.TimePtr(u.CreatedAt),
		UpdatedAt: convert.TimePtr(u.UpdatedAt),
	}
}

func userWithPasswordToModel(u *pbusers.UserResponseWithPassword) *salemongo.User {
	if u == nil {
		return nil
	}
	return &salemongo.User{
		UserID:    u.Id,
		Firstname: u.Firstname,
		Lastname:  u.Lastname,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: convert.TimePtr(u.CreatedAt),
		UpdatedAt: convert.TimePtr(u.UpdatedAt),
	}
}
