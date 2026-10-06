// Package category adapts the Category service gRPC API into the shared domain
// model.
package category

import (
	"context"

	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter/convert"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-pkg/resilience"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/category_errors"
)

// Repository reads categories from the Category service over gRPC.
type Repository interface {
	FindById(ctx context.Context, id int) (*salemongo.Category, error)
	FindByIds(ctx context.Context, ids []int) ([]*salemongo.Category, error)
	FindByName(ctx context.Context, name string) (*salemongo.Category, error)
}

type repository struct {
	client pbcategories.CategoryQueryServiceClient
	guard  *resilience.DependencyGuard
}

// New builds a category adapter. Passing zero options leaves the guard nil,
// which makes DependencyGuard.Call a plain passthrough.
func New(client pbcategories.CategoryQueryServiceClient, opts ...adapter.GuardOption) Repository {
	r := &repository{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *repository) FindById(ctx context.Context, id int) (*salemongo.Category, error) {
	var resp *pbcategories.ApiResponseCategory
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.client.FindById(ctx, &pbcategories.FindByIdCategoryRequest{Id: int32(id)})
		return callErr
	})
	if err != nil || resp == nil || resp.Data == nil {
		return nil, category_errors.ErrFindById
	}

	return mapCategory(resp.Data), nil
}

// FindByIds resolves many categories in a single call. Missing or soft-deleted
// ids are simply absent from the result; callers must tolerate that.
func (r *repository) FindByIds(ctx context.Context, ids []int) ([]*salemongo.Category, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	req := &pbcategories.FindByIdsCategoryRequest{Ids: make([]int32, 0, len(ids))}
	for _, id := range ids {
		req.Ids = append(req.Ids, int32(id))
	}

	var resp *pbcategories.ApiResponsesCategory
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.client.FindByIds(ctx, req)
		return callErr
	})
	if err != nil || resp == nil {
		return nil, category_errors.ErrFindByIds
	}

	categories := make([]*salemongo.Category, 0, len(resp.Data))
	for _, c := range resp.Data {
		if c == nil {
			continue
		}
		categories = append(categories, mapCategory(c))
	}
	return categories, nil
}

func (r *repository) FindByName(ctx context.Context, name string) (*salemongo.Category, error) {
	var resp *pbcategories.ApiResponseCategory
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.client.FindByName(ctx, &pbcategories.FindByNameCategoryRequest{Name: name})
		return callErr
	})
	if err != nil || resp == nil || resp.Data == nil {
		return nil, category_errors.ErrFindByName
	}

	return mapCategory(resp.Data), nil
}

// mapCategory converts a proto category into the shared domain model.
func mapCategory(c *pbcategories.CategoryResponse) *salemongo.Category {
	if c == nil {
		return nil
	}
	return &salemongo.Category{
		CategoryID:   c.Id,
		Name:         c.Name,
		Description:  convert.NullableString(c.Description),
		SlugCategory: convert.NullableString(c.SlugCategory),
		CreatedAt:    convert.TimePtr(c.CreatedAt),
		UpdatedAt:    convert.TimePtr(c.UpdatedAt),
	}
}
