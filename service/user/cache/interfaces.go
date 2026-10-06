package mencache

import (
	"context"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-user/repository"
)

type UserQueryCache interface {
	GetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, bool)
	SetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers, data []*repository.UserResult, total *int)
	GetCachedUserCache(ctx context.Context, id int) (*salemongo.User, bool)
	SetCachedUserCache(ctx context.Context, data *salemongo.User)
	GetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, bool)
	SetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers, data []*repository.UserResult, total *int)
	GetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers) ([]*repository.UserResult, *int, bool)
	SetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers, data []*repository.UserResult, total *int)
}

type UserCommandCache interface {
	DeleteUserCache(ctx context.Context, id int)
	DeleteUserAllCache(ctx context.Context)
}
