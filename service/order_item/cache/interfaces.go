package mencache

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-order-item/repository"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type OrderItemQueryCache interface {
	GetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems) ([]*repository.OrderItemResult, *int, bool)
	SetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems, data []*repository.OrderItemResult, total *int)

	GetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*repository.OrderItemResultDeleteAt, *int, bool)
	SetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems, data []*repository.OrderItemResultDeleteAt, total *int)

	GetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*repository.OrderItemResultDeleteAt, *int, bool)
	SetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems, data []*repository.OrderItemResultDeleteAt, total *int)

	GetCachedOrderItems(ctx context.Context, orderID int) ([]*salemongo.OrderItem, bool)
	SetCachedOrderItems(ctx context.Context, data []*salemongo.OrderItem)
}

// OrderItemCommandCache invalidates the read caches after a mutation so a
// create/update/delete never leaves stale line items behind (TTL is 5m).
type OrderItemCommandCache interface {
	DeleteCachedOrderItems(ctx context.Context, orderID int)
	DeleteCachedOrderItemAllCache(ctx context.Context)
}
