package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/category_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type categoryStatsByMerchantRepository struct {
	db *mongo.Database
}

func NewCategoryStatsByMerchantRepository(db *mongo.Database) CategoryStatsByMerchantRepository {
	return &categoryStatsByMerchantRepository{
		db: db,
	}
}

func (r *categoryStatsByMerchantRepository) GetMonthlyTotalPriceByMerchant(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*salemongo.CategoryMonthlyTotalPriceByMerchantRow, error) {
	coll := r.db.Collection("order_items")
	currentStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevStart := currentStart.AddDate(0, -1, 0)

	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.M{
			"from":         "products",
			"localField":   "product_id",
			"foreignField": "product_id",
			"as":           "product",
		}}},
		{{Key: "$unwind", Value: "$product"}},
		{{Key: "$match", Value: bson.M{
			"product.merchant_id": int32(req.MerchantID),
			"created_at":          bson.M{"$gte": prevStart, "$lte": currentStart.AddDate(0, 1, -1)},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"year":  bson.M{"$year": "$created_at"},
				"month": bson.M{"$month": "$created_at"},
			},
			"total_price": bson.M{"$sum": bson.M{"$multiply": []interface{}{"$price", "$quantity"}}},
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":         0,
			"year":        "$_id.year",
			"month":       "$_id.month",
			"total_price": 1,
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "year", Value: 1}, {Key: "month", Value: 1}}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, category_errors.ErrGetMonthlyTotalPriceByMerchant.WithInternal(err)
	}
	var results []*salemongo.CategoryMonthlyTotalPriceByMerchantRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, category_errors.ErrGetMonthlyTotalPriceByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *categoryStatsByMerchantRepository) GetYearlyTotalPricesByMerchant(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*salemongo.CategoryYearlyTotalPriceByMerchantRow, error) {
	coll := r.db.Collection("order_items")
	start := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)

	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.M{
			"from":         "products",
			"localField":   "product_id",
			"foreignField": "product_id",
			"as":           "product",
		}}},
		{{Key: "$unwind", Value: "$product"}},
		{{Key: "$match", Value: bson.M{
			"product.merchant_id": int32(req.MerchantID),
			"created_at":          bson.M{"$gte": start, "$lt": end},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":         bson.M{"year": bson.M{"$year": "$created_at"}},
			"total_price": bson.M{"$sum": bson.M{"$multiply": []interface{}{"$price", "$quantity"}}},
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":         0,
			"year":        "$_id.year",
			"total_price": 1,
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "year", Value: 1}}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, category_errors.ErrGetYearlyTotalPricesByMerchant.WithInternal(err)
	}
	var results []*salemongo.CategoryYearlyTotalPriceByMerchantRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, category_errors.ErrGetYearlyTotalPricesByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *categoryStatsByMerchantRepository) GetMonthPriceByMerchant(ctx context.Context, req *requests.MonthPriceMerchant) ([]*salemongo.CategoryMonthlyByMerchantRow, error) {
	coll := r.db.Collection("order_items")
	start := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)

	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.M{
			"from":         "products",
			"localField":   "product_id",
			"foreignField": "product_id",
			"as":           "product",
		}}},
		{{Key: "$unwind", Value: "$product"}},
		{{Key: "$match", Value: bson.M{
			"product.merchant_id": int32(req.MerchantID),
			"created_at":          bson.M{"$gte": start, "$lt": end},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":           bson.M{"category_id": "$product.category_id", "month": bson.M{"$month": "$created_at"}},
			"total_revenue": bson.M{"$sum": bson.M{"$multiply": []interface{}{"$price", "$quantity"}}},
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "categories",
			"localField":   "_id.category_id",
			"foreignField": "category_id",
			"as":           "category",
		}}},
		{{Key: "$unwind", Value: bson.M{"path": "$category", "preserveNullAndEmptyArrays": true}}},
		{{Key: "$project", Value: bson.M{
			"_id":           0,
			"category_id":   "$_id.category_id",
			"category_name": bson.M{"$ifNull": []interface{}{"$category.name", ""}},
			"month":         "$_id.month",
			"total_revenue": 1,
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "month", Value: 1}}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, category_errors.ErrGetMonthPriceByMerchant.WithInternal(err)
	}
	var results []*salemongo.CategoryMonthlyByMerchantRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, category_errors.ErrGetMonthPriceByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *categoryStatsByMerchantRepository) GetYearPriceByMerchant(ctx context.Context, req *requests.YearPriceMerchant) ([]*salemongo.CategoryYearlyByMerchantRow, error) {
	coll := r.db.Collection("order_items")
	start := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)

	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.M{
			"from":         "products",
			"localField":   "product_id",
			"foreignField": "product_id",
			"as":           "product",
		}}},
		{{Key: "$unwind", Value: "$product"}},
		{{Key: "$match", Value: bson.M{
			"product.merchant_id": int32(req.MerchantID),
			"created_at":          bson.M{"$gte": start, "$lt": end},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":           bson.M{"category_id": "$product.category_id", "year": bson.M{"$year": "$created_at"}},
			"total_revenue": bson.M{"$sum": bson.M{"$multiply": []interface{}{"$price", "$quantity"}}},
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "categories",
			"localField":   "_id.category_id",
			"foreignField": "category_id",
			"as":           "category",
		}}},
		{{Key: "$unwind", Value: bson.M{"path": "$category", "preserveNullAndEmptyArrays": true}}},
		{{Key: "$project", Value: bson.M{
			"_id":           0,
			"category_id":   "$_id.category_id",
			"category_name": bson.M{"$ifNull": []interface{}{"$category.name", ""}},
			"year":          "$_id.year",
			"total_revenue": 1,
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "year", Value: 1}}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, category_errors.ErrGetYearPriceByMerchant.WithInternal(err)
	}
	var results []*salemongo.CategoryYearlyByMerchantRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, category_errors.ErrGetYearPriceByMerchant.WithInternal(err)
	}
	return results, nil
}
