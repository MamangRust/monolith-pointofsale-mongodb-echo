package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type OrderItemResult struct {
	OrderItemID int32
	OrderID     int32
	ProductID   int32
	Quantity    int32
	Price       int64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	TotalCount  int64
}

type OrderItemResultDeleteAt struct {
	OrderItemID int32
	OrderID     int32
	ProductID   int32
	Quantity    int32
	Price       int64
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
	TotalCount  int64
}

type OrderItemQueryRepository interface {
	FindAllOrderItems(ctx context.Context, req *requests.FindAllOrderItems) ([]*OrderItemResult, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*OrderItemResultDeleteAt, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*OrderItemResultDeleteAt, *int, error)
	FindOrderItemByOrder(ctx context.Context, orderID int) ([]*salemongo.OrderItem, error)
}

// OrderItemCommandRepository owns the write side of order_items. order-service
// reaches it over gRPC for both the order flow and its compensating rollback.
type OrderItemCommandRepository interface {
	CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*salemongo.OrderItem, error)
	UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*salemongo.OrderItem, error)
	DeleteOrderItem(ctx context.Context, orderItemID int) error
}

// sqlc row type aliases for handler compatibility
type OrderItemRow = OrderItemResult
type OrderItemActiveRow = OrderItemResultDeleteAt
type OrderItemTrashedRow = OrderItemResultDeleteAt
