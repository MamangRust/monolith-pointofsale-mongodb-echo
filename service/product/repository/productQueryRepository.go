package repository

import (
	"context"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/product_errors"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type productQueryRepository struct {
	coll *mongodriver.Collection
}

func NewProductQueryRepository(coll *mongodriver.Collection) *productQueryRepository {
	return &productQueryRepository{coll: coll}
}

func (r *productQueryRepository) FindAllProducts(ctx context.Context, req *requests.FindAllProducts) ([]*ProductResult, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, product_errors.ErrFindAllProducts
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "product_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, product_errors.ErrFindAllProducts
	}
	defer cursor.Close(ctx)

	var products []salemongo.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, nil, product_errors.ErrFindAllProducts
	}

	var results []*ProductResult
	for _, p := range products {
		pr := &ProductResult{
			ProductID:    p.ProductID,
			MerchantID:   p.MerchantID,
			CategoryID:   p.CategoryID,
			Name:         p.Name,
			Description:  p.Description,
			Price:        p.Price,
			CountInStock: p.CountInStock,
			Brand:        p.Brand,
			Weight:       p.Weight,
			SlugProduct:  p.SlugProduct,
			ImageProduct: p.ImageProduct,
			TotalCount:   totalCount,
		}
		if p.CreatedAt != nil {
			pr.CreatedAt = p.CreatedAt
		}
		if p.UpdatedAt != nil {
			pr.UpdatedAt = p.UpdatedAt
		}
		results = append(results, pr)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *productQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllProducts) ([]*ProductResultDeleteAt, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, product_errors.ErrFindByActive
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "product_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, product_errors.ErrFindByActive
	}
	defer cursor.Close(ctx)

	var products []salemongo.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, nil, product_errors.ErrFindByActive
	}

	var results []*ProductResultDeleteAt
	for _, p := range products {
		pr := &ProductResultDeleteAt{
			ProductID:    p.ProductID,
			MerchantID:   p.MerchantID,
			CategoryID:   p.CategoryID,
			Name:         p.Name,
			Description:  p.Description,
			Price:        p.Price,
			CountInStock: p.CountInStock,
			Brand:        p.Brand,
			Weight:       p.Weight,
			SlugProduct:  p.SlugProduct,
			ImageProduct: p.ImageProduct,
			TotalCount:   totalCount,
		}
		if p.CreatedAt != nil {
			pr.CreatedAt = p.CreatedAt
		}
		if p.UpdatedAt != nil {
			pr.UpdatedAt = p.UpdatedAt
		}
		results = append(results, pr)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *productQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllProducts) ([]*ProductResultDeleteAt, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, product_errors.ErrFindByTrashed
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "product_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, product_errors.ErrFindByTrashed
	}
	defer cursor.Close(ctx)

	var products []salemongo.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, nil, product_errors.ErrFindByTrashed
	}

	var results []*ProductResultDeleteAt
	for _, p := range products {
		pr := &ProductResultDeleteAt{
			ProductID:    p.ProductID,
			MerchantID:   p.MerchantID,
			CategoryID:   p.CategoryID,
			Name:         p.Name,
			Description:  p.Description,
			Price:        p.Price,
			CountInStock: p.CountInStock,
			Brand:        p.Brand,
			Weight:       p.Weight,
			SlugProduct:  p.SlugProduct,
			ImageProduct: p.ImageProduct,
			TotalCount:   totalCount,
		}
		if p.CreatedAt != nil {
			pr.CreatedAt = p.CreatedAt
		}
		if p.UpdatedAt != nil {
			pr.UpdatedAt = p.UpdatedAt
		}
		if p.DeletedAt != nil {
			pr.DeletedAt = p.DeletedAt
		}
		results = append(results, pr)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *productQueryRepository) FindByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest) ([]*ProductByMerchantResult, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil, "merchant_id": req.MerchantID}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, product_errors.ErrFindByMerchant
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "product_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, product_errors.ErrFindByMerchant
	}
	defer cursor.Close(ctx)

	var products []salemongo.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, nil, product_errors.ErrFindByMerchant
	}

	var results []*ProductByMerchantResult
	for _, p := range products {
		pr := &ProductByMerchantResult{
			ProductID:    p.ProductID,
			Name:         p.Name,
			Description:  p.Description,
			Price:        p.Price,
			CountInStock: p.CountInStock,
			Brand:        p.Brand,
			ImageProduct: p.ImageProduct,
			TotalCount:   totalCount,
		}
		if p.CreatedAt != nil {
			pr.CreatedAt = p.CreatedAt
		}
		results = append(results, pr)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *productQueryRepository) FindByCategory(ctx context.Context, categoryID int32, req *requests.ProductByCategoryRequest) ([]*ProductByCategoryResult, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	// The caller resolves category_name → category_id through the category
	// service; this repository only reads the products collection it owns.
	filter := bson.M{"deleted_at": nil, "category_id": categoryID}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, product_errors.ErrFindByCategory
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "product_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, product_errors.ErrFindByCategory
	}
	defer cursor.Close(ctx)

	var products []salemongo.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, nil, product_errors.ErrFindByCategory
	}

	var results []*ProductByCategoryResult
	for _, p := range products {
		pr := &ProductByCategoryResult{
			ProductID:    p.ProductID,
			MerchantID:   p.MerchantID,
			CategoryID:   p.CategoryID,
			SlugProduct:  p.SlugProduct,
			Weight:       p.Weight,
			Name:         p.Name,
			Description:  p.Description,
			Price:        p.Price,
			CountInStock: p.CountInStock,
			TotalCount:   totalCount,
		}
		results = append(results, pr)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *productQueryRepository) FindById(ctx context.Context, product_id int) (*salemongo.Product, error) {
	var product salemongo.Product
	err := r.coll.FindOne(ctx, bson.M{"product_id": product_id, "deleted_at": nil}).Decode(&product)
	if err != nil {
		return nil, product_errors.ErrFindById
	}
	return &product, nil
}

func (r *productQueryRepository) FindByIdTrashed(ctx context.Context, id int) (*salemongo.Product, error) {
	var product salemongo.Product
	err := r.coll.FindOne(ctx, bson.M{"product_id": id, "deleted_at": bson.M{"$ne": nil}}).Decode(&product)
	if err != nil {
		return nil, product_errors.ErrFindById
	}
	return &product, nil
}
