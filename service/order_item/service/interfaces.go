package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-order-item/repository"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type OrderItemQueryService interface {
	FindAllOrderItems(ctx context.Context, req *requests.FindAllOrderItems) ([]*repository.OrderItemResult, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*repository.OrderItemResultDeleteAt, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*repository.OrderItemResultDeleteAt, *int, error)
	FindOrderItemByOrder(ctx context.Context, orderID int) ([]*salemongo.OrderItem, error)
}

// OrderItemCommandService owns the write side of order_items. order-service
// reaches it over gRPC for both the order flow and its compensating rollback.
type OrderItemCommandService interface {
	CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*salemongo.OrderItem, error)
	UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*salemongo.OrderItem, error)
	DeleteOrderItem(ctx context.Context, orderItemID int) error
}
