// Package orderitem adapts the OrderItem service gRPC API into the shared
// domain model. Query and bulk are exposed as separate interfaces so consumers
// can depend on the narrower one.
package orderitem

import (
	"context"

	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter/convert"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-pkg/resilience"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	orderitem_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/order_item_errors"
)

// QueryRepository reads order items from the OrderItem service over gRPC.
type QueryRepository interface {
	FindOrderItemByOrder(ctx context.Context, orderID int) ([]*salemongo.OrderItem, error)
	CalculateTotalPrice(ctx context.Context, orderID int) (*int32, error)
}

// CommandRepository mutates order items via the OrderItem service over gRPC.
// order_items stays owned by the order_item service; consumers reach it only
// through this adapter.
type CommandRepository interface {
	CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*salemongo.OrderItem, error)
	UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*salemongo.OrderItem, error)
	DeleteOrderItem(ctx context.Context, orderItemID int) error
}

// BulkRepository enumerates order items page by page. It is deliberately a
// separate narrow interface: offline consumers (the stats backfill) need the
// whole collection, while ordinary consumers only ever read order items for a
// single order and should not depend on the extra method.
type BulkRepository interface {
	FindAllOrderItems(ctx context.Context, page, pageSize int) ([]*salemongo.OrderItem, int, error)
}

// Repository implements QueryRepository, CommandRepository and BulkRepository
// over a single pair of gRPC clients, so consumers can share one instance
// across all roles.
type Repository struct {
	query   pborder_items.OrderItemQueryServiceClient
	command pborder_items.OrderItemCommandServiceClient
	guard   *resilience.DependencyGuard
}

// New builds an order-item adapter. Passing zero options leaves the guard nil,
// which makes DependencyGuard.Call a plain passthrough. Query-only consumers
// may pass nil for the command client.
func New(query pborder_items.OrderItemQueryServiceClient, command pborder_items.OrderItemCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{query: query, command: command}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *Repository) CalculateTotalPrice(ctx context.Context, orderID int) (*int32, error) {
	items, err := r.FindOrderItemByOrder(ctx, orderID)
	if err != nil {
		return nil, orderitem_errors.ErrCalculateTotalPrice
	}

	var total int32 = 0
	for _, item := range items {
		if item != nil {
			total += item.Quantity * int32(item.Price)
		}
	}

	return &total, nil
}

func (r *Repository) FindOrderItemByOrder(ctx context.Context, orderID int) ([]*salemongo.OrderItem, error) {
	var resp *pborder_items.ApiResponsesOrderItem
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.query.FindOrderItemByOrder(ctx, &pborder_items.FindByIdOrderItemRequest{Id: int32(orderID)})
		return callErr
	})
	if err != nil || resp == nil || resp.Data == nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder
	}

	var res []*salemongo.OrderItem
	for _, item := range resp.Data {
		if item == nil {
			continue
		}
		res = append(res, mapOrderItem(item))
	}

	return res, nil
}

// FindAllOrderItems returns one page of active order items plus the total
// number of active rows, so callers can walk the collection page by page.
func (r *Repository) FindAllOrderItems(ctx context.Context, page, pageSize int) ([]*salemongo.OrderItem, int, error) {
	var resp *pborder_items.ApiResponsePaginationOrderItem
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.query.FindAll(ctx, &pborder_items.FindAllOrderItemRequest{
			Page:     int32(page),
			PageSize: int32(pageSize),
		})
		return callErr
	})
	if err != nil || resp == nil {
		return nil, 0, orderitem_errors.ErrFindAllOrderItems
	}

	items := make([]*salemongo.OrderItem, 0, len(resp.Data))
	for _, item := range resp.Data {
		if item == nil {
			continue
		}
		items = append(items, mapOrderItem(item))
	}

	total := 0
	if resp.Pagination != nil {
		total = int(resp.Pagination.TotalRecords)
	}
	return items, total, nil
}

func (r *Repository) CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*salemongo.OrderItem, error) {
	if req == nil || r.command == nil {
		return nil, orderitem_errors.ErrCreateOrderItem
	}

	var resp *pborder_items.ApiResponseOrderItem
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.command.CreateOrderItem(ctx, &pborder_items.CreateOrderItemRecordRequest{
			OrderId:   int32(req.OrderID),
			ProductId: int32(req.ProductID),
			Quantity:  int32(req.Quantity),
			Price:     int32(req.Price),
		})
		return callErr
	})
	if err != nil || resp == nil || resp.Data == nil {
		return nil, orderitem_errors.ErrCreateOrderItem
	}

	return mapOrderItem(resp.Data), nil
}

func (r *Repository) UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*salemongo.OrderItem, error) {
	if req == nil || r.command == nil {
		return nil, orderitem_errors.ErrUpdateOrderItem
	}

	var resp *pborder_items.ApiResponseOrderItem
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.command.UpdateOrderItem(ctx, &pborder_items.UpdateOrderItemRecordRequest{
			OrderItemId: int32(req.OrderItemID),
			OrderId:     int32(req.OrderID),
			ProductId:   int32(req.ProductID),
			Quantity:    int32(req.Quantity),
			Price:       int32(req.Price),
		})
		return callErr
	})
	if err != nil || resp == nil || resp.Data == nil {
		return nil, orderitem_errors.ErrUpdateOrderItem
	}

	return mapOrderItem(resp.Data), nil
}

// DeleteOrderItem removes one order item by its own id (the rollback path of
// order-service).
func (r *Repository) DeleteOrderItem(ctx context.Context, orderItemID int) error {
	if r.command == nil {
		return orderitem_errors.ErrDeleteOrderItemPermanent
	}

	var resp *pborder_items.ApiResponseOrderItemDelete
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.command.DeleteOrderItem(ctx, &pborder_items.DeleteOrderItemRecordRequest{
			OrderItemId: int32(orderItemID),
		})
		return callErr
	})
	if err != nil || resp == nil {
		return orderitem_errors.ErrDeleteOrderItemPermanent
	}

	return nil
}

// mapOrderItem converts a proto order item into the shared domain model.
func mapOrderItem(item *pborder_items.OrderItemResponse) *salemongo.OrderItem {
	if item == nil {
		return nil
	}
	return &salemongo.OrderItem{
		OrderItemID: item.Id,
		OrderID:     item.OrderId,
		ProductID:   item.ProductId,
		Quantity:    item.Quantity,
		Price:       int64(item.Price),
		CreatedAt:   convert.TimePtr(item.CreatedAt),
		UpdatedAt:   convert.TimePtr(item.UpdatedAt),
	}
}
