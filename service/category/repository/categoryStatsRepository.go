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

type categoryStatsRepository struct {
	db *mongo.Database
}

func NewCategoryStatsRepository(db *mongo.Database) CategoryStatsRepository {
	return &categoryStatsRepository{
		db: db,
	}
}

func (r *categoryStatsRepository) GetMonthlyTotalPrice(ctx context.Context, req *requests.MonthTotalPrice) ([]*salemongo.CategoryMonthlyTotalPriceRow, error) {
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
			"created_at": bson.M{"$gte": prevStart, "$lte": currentStart.AddDate(0, 1, -1)},
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
		return nil, category_errors.ErrGetMonthlyTotalPrice.WithInternal(err)
	}
	var results []*salemongo.CategoryMonthlyTotalPriceRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, category_errors.ErrGetMonthlyTotalPrice.WithInternal(err)
	}
	return results, nil
}

func (r *categoryStatsRepository) GetYearlyTotalPrices(ctx context.Context, year int) ([]*salemongo.CategoryYearlyTotalPriceRow, error) {
	coll := r.db.Collection("order_items")
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
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
			"created_at": bson.M{"$gte": start, "$lt": end},
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
		return nil, category_errors.ErrGetYearlyTotalPrices.WithInternal(err)
	}
	var results []*salemongo.CategoryYearlyTotalPriceRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, category_errors.ErrGetYearlyTotalPrices.WithInternal(err)
	}
	return results, nil
}

func (r *categoryStatsRepository) GetMonthPrice(ctx context.Context, year int) ([]*salemongo.CategoryMonthlyRow, error) {
	coll := r.db.Collection("order_items")
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
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
			"created_at": bson.M{"$gte": start, "$lt": end},
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
		return nil, category_errors.ErrGetMonthPrice.WithInternal(err)
	}
	var results []*salemongo.CategoryMonthlyRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, category_errors.ErrGetMonthPrice.WithInternal(err)
	}
	return results, nil
}

func (r *categoryStatsRepository) GetYearPrice(ctx context.Context, year int) ([]*salemongo.CategoryYearlyRow, error) {
	coll := r.db.Collection("order_items")
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
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
			"created_at": bson.M{"$gte": start, "$lt": end},
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
		return nil, category_errors.ErrGetYearPrice.WithInternal(err)
	}
	var results []*salemongo.CategoryYearlyRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, category_errors.ErrGetYearPrice.WithInternal(err)
	}
	return results, nil
}
