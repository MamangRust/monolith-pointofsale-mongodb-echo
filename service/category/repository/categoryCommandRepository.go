package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/category_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type categoryCommandRepository struct {
	coll *mongo.Collection
}

func NewCategoryCommandRepository(coll *mongo.Collection) CategoryCommandRepository {
	return &categoryCommandRepository{coll: coll}
}

func (r *categoryCommandRepository) CreateCategory(ctx context.Context, request *requests.CreateCategoryRequest) (*salemongo.Category, error) {
	categoryID, err := salemongo.NextID(ctx, r.coll.Database(), "categories")
	if err != nil {
		return nil, category_errors.ErrCreateCategory
	}

	now := time.Now()
	category := &salemongo.Category{
		CategoryID:   categoryID,
		Name:         request.Name,
		Description:  &request.Description,
		SlugCategory: request.SlugCategory,
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}

	result, err := r.coll.InsertOne(ctx, category)
	if err != nil {
		return nil, category_errors.ErrCreateCategory
	}

	category.ID = result.InsertedID.(primitive.ObjectID)
	return category, nil
}

func (r *categoryCommandRepository) UpdateCategory(ctx context.Context, request *requests.UpdateCategoryRequest) (*salemongo.Category, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"category_id": *request.CategoryID},
		bson.M{
			"$set": bson.M{
				"name":          request.Name,
				"description":   request.Description,
				"slug_category": request.SlugCategory,
				"updated_at":    now,
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var category salemongo.Category
	err := result.Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, category_errors.ErrUpdateCategory
		}
		return nil, category_errors.ErrUpdateCategory
	}
	return &category, nil
}

func (r *categoryCommandRepository) TrashedCategory(ctx context.Context, category_id int) (*salemongo.Category, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"category_id": category_id},
		bson.M{
			"$set": bson.M{
				"deleted_at": now,
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var category salemongo.Category
	err := result.Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, category_errors.ErrTrashedCategory
		}
		return nil, category_errors.ErrTrashedCategory
	}
	return &category, nil
}

func (r *categoryCommandRepository) RestoreCategory(ctx context.Context, category_id int) (*salemongo.Category, error) {
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"category_id": category_id},
		bson.M{
			"$unset": bson.M{
				"deleted_at": "",
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var category salemongo.Category
	err := result.Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, category_errors.ErrRestoreCategory
		}
		return nil, category_errors.ErrRestoreCategory
	}
	return &category, nil
}

func (r *categoryCommandRepository) DeleteCategoryPermanently(ctx context.Context, category_id int) (bool, error) {
	result, err := r.coll.DeleteOne(ctx, bson.M{"category_id": category_id})
	if err != nil {
		return false, category_errors.ErrDeleteCategoryPermanently
	}
	return result.DeletedCount > 0, nil
}

func (r *categoryCommandRepository) RestoreAllCategories(ctx context.Context) (bool, error) {
	result, err := r.coll.UpdateMany(
		ctx,
		bson.M{"deleted_at": bson.M{"$ne": nil}},
		bson.M{
			"$unset": bson.M{
				"deleted_at": "",
			},
		},
	)
	if err != nil {
		return false, category_errors.ErrRestoreAllCategories
	}
	return result.ModifiedCount > 0, nil
}

func (r *categoryCommandRepository) DeleteAllPermanentCategories(ctx context.Context) (bool, error) {
	result, err := r.coll.DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, category_errors.ErrDeleteAllPermanentCategories
	}
	return result.DeletedCount > 0, nil
}
