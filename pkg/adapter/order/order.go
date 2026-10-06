// Package order adapts the Order service gRPC API into the shared domain model.
// Query and bulk are exposed as separate interfaces so consumers depend on the
// narrower one: ordinary consumers read a single order, while offline consumers
// (the stats backfill) enumerate the whole collection.
package order

import (
	"context"

	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter/convert"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-pkg/resilience"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
)

// QueryRepository reads a single order from the Order service over gRPC.
type QueryRepository interface {
	FindById(ctx context.Context, id int) (*salemongo.Order, error)
}

// BulkRepository enumerates orders page by page. It is deliberately a separate
// narrow interface: offline consumers (the stats backfill) need the whole
// collection, while ordinary consumers only ever read one order at a time and
// should not depend on the extra method.
type BulkRepository interface {
	FindAll(ctx context.Context, page, pageSize int) ([]*salemongo.Order, int, error)
}

// Repository implements QueryRepository and BulkRepository over a single gRPC
// client, so consumers can share one instance across both roles.
type Repository struct {
	client pborders.OrderQueryServiceClient
	guard  *resilience.DependencyGuard
}

// New builds an order adapter. Passing zero options leaves the guard nil, which
// makes DependencyGuard.Call a plain passthrough.
func New(client pborders.OrderQueryServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *Repository) FindById(ctx context.Context, id int) (*salemongo.Order, error) {
	var resp *pborders.ApiResponseOrder
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.client.FindById(ctx, &pborders.FindByIdOrderRequest{Id: int32(id)})
		return callErr
	})
	if err != nil || resp == nil || resp.Data == nil {
		return nil, order_errors.ErrFindById
	}

	return mapOrder(resp.Data), nil
}

// FindAll returns one page of active orders plus the total number of active
// rows, so callers can walk the collection page by page.
func (r *Repository) FindAll(ctx context.Context, page, pageSize int) ([]*salemongo.Order, int, error) {
	var resp *pborders.ApiResponsePaginationOrder
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.client.FindAll(ctx, &pborders.FindAllOrderRequest{
			Page:     int32(page),
			PageSize: int32(pageSize),
		})
		return callErr
	})
	if err != nil || resp == nil {
		return nil, 0, order_errors.ErrFindAllOrders
	}

	orders := make([]*salemongo.Order, 0, len(resp.Data))
	for _, o := range resp.Data {
		if o == nil {
			continue
		}
		orders = append(orders, mapOrder(o))
	}

	total := 0
	if resp.Pagination != nil {
		total = int(resp.Pagination.TotalRecords)
	}
	return orders, total, nil
}

// mapOrder converts a proto order into the shared domain model.
func mapOrder(o *pborders.OrderResponse) *salemongo.Order {
	if o == nil {
		return nil
	}
	return &salemongo.Order{
		OrderID:    o.Id,
		MerchantID: o.MerchantId,
		CashierID:  o.CashierId,
		TotalPrice: int64(o.TotalPrice),
		CreatedAt:  convert.TimePtr(o.CreatedAt),
		UpdatedAt:  convert.TimePtr(o.UpdatedAt),
	}
}
