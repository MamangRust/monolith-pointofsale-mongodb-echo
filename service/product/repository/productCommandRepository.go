package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/product_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type productCommandRepository struct {
	coll *mongodriver.Collection
}

func NewProductCommandRepository(coll *mongodriver.Collection) *productCommandRepository {
	return &productCommandRepository{coll: coll}
}

func timePtr(t time.Time) *time.Time { return &t }

func (r *productCommandRepository) CreateProduct(ctx context.Context, request *requests.CreateProductRequest) (*salemongo.Product, error) {
	productID, err := salemongo.NextID(ctx, r.coll.Database(), "products")
	if err != nil {
		return nil, product_errors.ErrCreateProduct
	}

	now := time.Now()
	weight := int32(request.Weight)
	product := &salemongo.Product{
		ProductID:    productID,
		MerchantID:   int32(request.MerchantID),
		CategoryID:   int32(request.CategoryID),
		Name:         request.Name,
		Description:  &request.Description,
		Price:        int32(request.Price),
		CountInStock: int32(request.CountInStock),
		Brand:        &request.Brand,
		Weight:       &weight,
		SlugProduct:  request.SlugProduct,
		ImageProduct: &request.ImageProduct,
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}
	result, err := r.coll.InsertOne(ctx, product)
	if err != nil {
		return nil, product_errors.ErrCreateProduct
	}
	product.ID = result.InsertedID.(primitive.ObjectID)
	return product, nil
}

func (r *productCommandRepository) UpdateProduct(ctx context.Context, request *requests.UpdateProductRequest) (*salemongo.Product, error) {
	now := time.Now()
	weight := int32(request.Weight)
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"product_id": *request.ProductID, "deleted_at": nil},
		bson.M{"$set": bson.M{
			"category_id":    int32(request.CategoryID),
			"name":           request.Name,
			"description":    request.Description,
			"price":          int32(request.Price),
			"count_in_stock": int32(request.CountInStock),
			"brand":          request.Brand,
			"weight":         weight,
			"image_product":  request.ImageProduct,
			"updated_at":     now,
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var product salemongo.Product
	if err := result.Decode(&product); err != nil {
		return nil, product_errors.ErrUpdateProduct
	}
	return &product, nil
}

func (r *productCommandRepository) UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*salemongo.Product, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"product_id": product_id, "deleted_at": nil},
		bson.M{"$set": bson.M{"count_in_stock": int32(stock), "updated_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var product salemongo.Product
	if err := result.Decode(&product); err != nil {
		return nil, product_errors.ErrUpdateProductCountStock
	}
	return &product, nil
}

func (r *productCommandRepository) DecrementProductCountStock(ctx context.Context, productID int, quantity int) (*salemongo.Product, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"product_id": productID, "deleted_at": nil},
		bson.M{"$inc": bson.M{"count_in_stock": -int32(quantity)}, "$set": bson.M{"updated_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var product salemongo.Product
	if err := result.Decode(&product); err != nil {
		return nil, product_errors.ErrUpdateProductCountStock
	}
	if product.CountInStock < 0 {
		// Revert so a failed reservation never leaves a negative stock behind.
		_, _ = r.coll.UpdateOne(ctx, bson.M{"product_id": productID}, bson.M{"$inc": bson.M{"count_in_stock": int32(quantity)}})
		return nil, product_errors.ErrInsufficientProductStock
	}
	return &product, nil
}

func (r *productCommandRepository) IncrementProductCountStock(ctx context.Context, productID int, quantity int) (*salemongo.Product, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"product_id": productID, "deleted_at": nil},
		bson.M{"$inc": bson.M{"count_in_stock": int32(quantity)}, "$set": bson.M{"updated_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var product salemongo.Product
	if err := result.Decode(&product); err != nil {
		return nil, product_errors.ErrUpdateProductCountStock
	}
	return &product, nil
}

func (r *productCommandRepository) TrashedProduct(ctx context.Context, product_id int) (*salemongo.Product, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"product_id": product_id, "deleted_at": nil},
		bson.M{"$set": bson.M{"deleted_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var product salemongo.Product
	if err := result.Decode(&product); err != nil {
		return nil, product_errors.ErrTrashedProduct
	}
	return &product, nil
}

func (r *productCommandRepository) RestoreProduct(ctx context.Context, product_id int) (*salemongo.Product, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"product_id": product_id, "deleted_at": bson.M{"$ne": nil}},
		bson.M{"$unset": bson.M{"deleted_at": ""}, "$set": bson.M{"updated_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var product salemongo.Product
	if err := result.Decode(&product); err != nil {
		return nil, product_errors.ErrRestoreProduct
	}
	return &product, nil
}

func (r *productCommandRepository) DeleteProductPermanent(ctx context.Context, product_id int) (bool, error) {
	result, err := r.coll.DeleteOne(ctx, bson.M{"product_id": product_id, "deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, err
	}
	return result.DeletedCount > 0, nil
}

func (r *productCommandRepository) RestoreAllProducts(ctx context.Context) (bool, error) {
	_, err := r.coll.UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$unset": bson.M{"deleted_at": ""}})
	return true, err
}

func (r *productCommandRepository) DeleteAllProductPermanent(ctx context.Context) (bool, error) {
	_, err := r.coll.DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	return true, err
}
