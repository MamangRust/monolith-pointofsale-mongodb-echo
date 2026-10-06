package repository

import (
	"context"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/category_errors"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type categoryQueryRepository struct {
	coll *mongodriver.Collection
}

func NewCategoryQueryRepository(coll *mongodriver.Collection) CategoryQueryRepository {
	return &categoryQueryRepository{coll: coll}
}

func (r *categoryQueryRepository) FindAllCategory(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResult, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, category_errors.ErrFindAllCategory
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "category_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, category_errors.ErrFindAllCategory
	}
	defer cursor.Close(ctx)

	var categories []salemongo.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, nil, category_errors.ErrFindAllCategory
	}

	var results []*CategoryResult
	for _, c := range categories {
		cr := &CategoryResult{
			CategoryID:   c.CategoryID,
			Name:         c.Name,
			Description:  c.Description,
			SlugCategory: c.SlugCategory,
			TotalCount:   totalCount,
		}
		if c.CreatedAt != nil {
			cr.CreatedAt = c.CreatedAt
		}
		if c.UpdatedAt != nil {
			cr.UpdatedAt = c.UpdatedAt
		}
		results = append(results, cr)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *categoryQueryRepository) FindById(ctx context.Context, category_id int) (*salemongo.Category, error) {
	var category salemongo.Category
	err := r.coll.FindOne(ctx, bson.M{"category_id": category_id, "deleted_at": nil}).Decode(&category)
	if err != nil {
		return nil, category_errors.ErrFindById
	}
	return &category, nil
}

func (r *categoryQueryRepository) FindByIdTrashed(ctx context.Context, category_id int) (*salemongo.Category, error) {
	var category salemongo.Category
	err := r.coll.FindOne(ctx, bson.M{"category_id": category_id}).Decode(&category)
	if err != nil {
		return nil, category_errors.ErrFindById
	}
	return &category, nil
}

func (r *categoryQueryRepository) FindByNameAndId(ctx context.Context, req *requests.CategoryNameAndId) (*salemongo.Category, error) {
	var category salemongo.Category
	err := r.coll.FindOne(ctx, bson.M{"name": req.Name, "category_id": bson.M{"$ne": req.CategoryID}, "deleted_at": nil}).Decode(&category)
	if err != nil {
		return nil, category_errors.ErrFindByNameAndId
	}
	return &category, nil
}

func (r *categoryQueryRepository) FindByName(ctx context.Context, name string) (*salemongo.Category, error) {
	var category salemongo.Category
	err := r.coll.FindOne(ctx, bson.M{"name": name, "deleted_at": nil}).Decode(&category)
	if err != nil {
		return nil, category_errors.ErrFindByName
	}
	return &category, nil
}

// FindByIds resolves many categories in a single $in query. Missing or
// soft-deleted ids are simply absent from the result; callers must tolerate
// that (batch name-fill for product listings).
func (r *categoryQueryRepository) FindByIds(ctx context.Context, ids []int) ([]*salemongo.Category, error) {
	reqIDs := make([]int32, 0, len(ids))
	for _, id := range ids {
		reqIDs = append(reqIDs, int32(id))
	}

	cursor, err := r.coll.Find(ctx, bson.M{"category_id": bson.M{"$in": reqIDs}, "deleted_at": nil})
	if err != nil {
		return nil, category_errors.ErrFindByIds
	}
	defer cursor.Close(ctx)

	var categories []*salemongo.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, category_errors.ErrFindByIds
	}
	return categories, nil
}

func (r *categoryQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResultDeleteAt, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, category_errors.ErrFindByActive
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "category_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, category_errors.ErrFindByActive
	}
	defer cursor.Close(ctx)

	var categories []salemongo.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, nil, category_errors.ErrFindByActive
	}

	var results []*CategoryResultDeleteAt
	for _, c := range categories {
		cr := &CategoryResultDeleteAt{
			CategoryID:   c.CategoryID,
			Name:         c.Name,
			Description:  c.Description,
			SlugCategory: c.SlugCategory,
			TotalCount:   totalCount,
		}
		if c.CreatedAt != nil {
			cr.CreatedAt = c.CreatedAt
		}
		if c.UpdatedAt != nil {
			cr.UpdatedAt = c.UpdatedAt
		}
		results = append(results, cr)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *categoryQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*CategoryResultDeleteAt, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, category_errors.ErrFindByTrashed
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "category_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, category_errors.ErrFindByTrashed
	}
	defer cursor.Close(ctx)

	var categories []salemongo.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, nil, category_errors.ErrFindByTrashed
	}

	var results []*CategoryResultDeleteAt
	for _, c := range categories {
		cr := &CategoryResultDeleteAt{
			CategoryID:   c.CategoryID,
			Name:         c.Name,
			Description:  c.Description,
			SlugCategory: c.SlugCategory,
			TotalCount:   totalCount,
		}
		if c.CreatedAt != nil {
			cr.CreatedAt = c.CreatedAt
		}
		if c.UpdatedAt != nil {
			cr.UpdatedAt = c.UpdatedAt
		}
		if c.DeletedAt != nil {
			cr.DeletedAt = c.DeletedAt
		}
		results = append(results, cr)
	}
	total := int(totalCount)
	return results, &total, nil
}
