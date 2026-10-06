// Package product adapts the Product service gRPC API into the shared domain
// model. Query and bulk are exposed as separate interfaces so consumers can
// depend on the narrower one.
package product

import (
	"context"

	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter/convert"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-pkg/resilience"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/product_errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// QueryRepository reads products from the Product service over gRPC.
type QueryRepository interface {
	FindById(ctx context.Context, id int) (*salemongo.Product, error)
}

// CommandRepository mutates product stock through the Product service, which
// owns the products collection. Consumers must never write to products
// directly.
type CommandRepository interface {
	DecrementProductCountStock(ctx context.Context, productID int, quantity int) (*salemongo.Product, error)
	IncrementProductCountStock(ctx context.Context, productID int, quantity int) (*salemongo.Product, error)
}

// BulkRepository enumerates products page by page. It is deliberately a
// separate narrow interface: offline consumers (the stats backfill) need the
// whole catalog, while ordinary consumers only ever read one product at a time
// and should not depend on the extra method.
type BulkRepository interface {
	FindAll(ctx context.Context, page, pageSize int) ([]*salemongo.Product, int, error)
}

// Repository implements QueryRepository, CommandRepository and BulkRepository
// over a single pair of gRPC clients, so consumers can share one instance
// across all roles.
type Repository struct {
	query   pbproducts.ProductQueryServiceClient
	command pbproducts.ProductCommandServiceClient
	guard   *resilience.DependencyGuard
}

// New builds a product adapter. Passing zero options leaves the guard nil,
// which makes DependencyGuard.Call a plain passthrough. Query-only consumers
// may pass nil for the command client.
func New(query pbproducts.ProductQueryServiceClient, command pbproducts.ProductCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{query: query, command: command}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *Repository) FindById(ctx context.Context, id int) (*salemongo.Product, error) {
	var resp *pbproducts.ApiResponseProduct
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.query.FindById(ctx, &pbproducts.FindByIdProductRequest{Id: int32(id)})
		return callErr
	})
	if err != nil || resp == nil || resp.Data == nil {
		return nil, product_errors.ErrFindById
	}

	return mapResponse(resp.Data), nil
}

// FindAll returns one page of active products plus the total number of active
// rows, so callers can walk the catalog page by page.
func (r *Repository) FindAll(ctx context.Context, page, pageSize int) ([]*salemongo.Product, int, error) {
	var resp *pbproducts.ApiResponsePaginationProduct
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.query.FindAll(ctx, &pbproducts.FindAllProductRequest{
			Page:     int32(page),
			PageSize: int32(pageSize),
		})
		return callErr
	})
	if err != nil || resp == nil {
		return nil, 0, product_errors.ErrFindAllProducts
	}

	products := make([]*salemongo.Product, 0, len(resp.Data))
	for _, p := range resp.Data {
		if p == nil {
			continue
		}
		products = append(products, mapResponse(p))
	}

	total := 0
	if resp.Pagination != nil {
		total = int(resp.Pagination.TotalRecords)
	}
	return products, total, nil
}

// DecrementProductCountStock removes quantity from a product's stock via the
// Product service. Insufficient stock comes back as
// order_errors.ErrInsufficientProductStock (gRPC FailedPrecondition).
func (r *Repository) DecrementProductCountStock(ctx context.Context, productID int, quantity int) (*salemongo.Product, error) {
	return r.adjustProductCountStock(ctx, productID, quantity, true)
}

func (r *Repository) IncrementProductCountStock(ctx context.Context, productID int, quantity int) (*salemongo.Product, error) {
	return r.adjustProductCountStock(ctx, productID, quantity, false)
}

// adjustProductCountStock delegates the stock mutation to the Product service,
// which owns the products collection.
func (r *Repository) adjustProductCountStock(ctx context.Context, productID int, quantity int, decrement bool) (*salemongo.Product, error) {
	if r.command == nil {
		return nil, product_errors.ErrUpdateProductCountStock
	}

	req := &pbproducts.AdjustProductStockRequest{
		ProductId: int32(productID),
		Quantity:  int32(quantity),
	}

	var resp *pbproducts.ApiResponseProduct
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		if decrement {
			resp, callErr = r.command.DecrementStock(ctx, req)
		} else {
			resp, callErr = r.command.IncrementStock(ctx, req)
		}
		return callErr
	})
	if err != nil {
		if status.Code(err) == codes.FailedPrecondition {
			return nil, order_errors.ErrInsufficientProductStock
		}
		return nil, product_errors.ErrUpdateProductCountStock
	}
	if resp == nil || resp.Data == nil {
		return nil, product_errors.ErrUpdateProductCountStock
	}

	return mapResponse(resp.Data), nil
}

// mapResponse converts a proto product into the shared domain model.
func mapResponse(p *pbproducts.ProductResponse) *salemongo.Product {
	if p == nil {
		return nil
	}
	return &salemongo.Product{
		ProductID:    p.Id,
		MerchantID:   p.MerchantId,
		CategoryID:   p.CategoryId,
		Name:         p.Name,
		Description:  convert.NullableString(p.Description),
		Price:        p.Price,
		CountInStock: p.CountInStock,
		Brand:        convert.NullableString(p.Brand),
		Weight:       convert.NullableInt32(p.Weight),
		SlugProduct:  convert.NullableString(p.SlugProduct),
		ImageProduct: convert.NullableString(p.ImageProduct),
		Barcode:      convert.NullableString(p.Barcode),
		CreatedAt:    convert.TimePtr(p.CreatedAt),
		UpdatedAt:    convert.TimePtr(p.UpdatedAt),
	}
}
