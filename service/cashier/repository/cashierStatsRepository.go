package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type cashierStatsRepository struct {
	db *mongo.Database
}

func NewCashierStatsRepository(db *mongo.Database) CashierStatsRepository {
	return &cashierStatsRepository{
		db: db,
	}
}

func (r *cashierStatsRepository) GetMonthlyTotalSales(ctx context.Context, req *requests.MonthTotalSales) ([]*salemongo.CashierMonthlyTotalSalesRow, error) {
	coll := r.db.Collection("orders")
	currentStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevStart := currentStart.AddDate(0, -1, 0)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"created_at": bson.M{"$gte": prevStart, "$lte": currentStart.AddDate(0, 1, -1)},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"year":  bson.M{"$year": "$created_at"},
				"month": bson.M{"$month": "$created_at"},
			},
			"total_sales": bson.M{"$sum": "$total_price"},
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":         0,
			"year":        "$_id.year",
			"month":       "$_id.month",
			"total_sales": 1,
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "year", Value: 1}, {Key: "month", Value: 1}}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, cashier_errors.ErrGetMonthlyTotalSales.WithInternal(err)
	}
	var results []*salemongo.CashierMonthlyTotalSalesRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, cashier_errors.ErrGetMonthlyTotalSales.WithInternal(err)
	}
	return results, nil
}

func (r *cashierStatsRepository) GetYearlyTotalSales(ctx context.Context, year int) ([]*salemongo.CashierYearlyTotalSalesRow, error) {
	coll := r.db.Collection("orders")
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"created_at": bson.M{"$gte": start, "$lt": end},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":         bson.M{"year": bson.M{"$year": "$created_at"}},
			"total_sales": bson.M{"$sum": "$total_price"},
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":         0,
			"year":        "$_id.year",
			"total_sales": 1,
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "year", Value: 1}}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, cashier_errors.ErrGetYearlyTotalSales.WithInternal(err)
	}
	var results []*salemongo.CashierYearlyTotalSalesRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, cashier_errors.ErrGetYearlyTotalSales.WithInternal(err)
	}
	return results, nil
}

func (r *cashierStatsRepository) GetMonthyCashier(ctx context.Context, year int) ([]*salemongo.CashierMonthlySaleRow, error) {
	coll := r.db.Collection("orders")
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"created_at": bson.M{"$gte": start, "$lt": end},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"cashier_id": "$cashier_id",
				"month":      bson.M{"$month": "$created_at"},
			},
			"order_count": bson.M{"$sum": 1},
			"total_sales": bson.M{"$sum": "$total_price"},
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "cashiers",
			"localField":   "_id.cashier_id",
			"foreignField": "cashier_id",
			"as":           "cashier",
		}}},
		{{Key: "$unwind", Value: bson.M{"path": "$cashier", "preserveNullAndEmptyArrays": true}}},
		{{Key: "$project", Value: bson.M{
			"_id":          0,
			"cashier_id":   "$_id.cashier_id",
			"cashier_name": bson.M{"$ifNull": []interface{}{"$cashier.name", ""}},
			"month":        "$_id.month",
			"order_count":  1,
			"total_sales":  1,
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "month", Value: 1}}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, cashier_errors.ErrGetMonthlyCashier.WithInternal(err)
	}
	var results []*salemongo.CashierMonthlySaleRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, cashier_errors.ErrGetMonthlyCashier.WithInternal(err)
	}
	return results, nil
}

func (r *cashierStatsRepository) GetYearlyCashier(ctx context.Context, year int) ([]*salemongo.CashierYearlySaleRow, error) {
	coll := r.db.Collection("orders")
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"created_at": bson.M{"$gte": start, "$lt": end},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":         bson.M{"cashier_id": "$cashier_id", "year": bson.M{"$year": "$created_at"}},
			"order_count": bson.M{"$sum": 1},
			"total_sales": bson.M{"$sum": "$total_price"},
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "cashiers",
			"localField":   "_id.cashier_id",
			"foreignField": "cashier_id",
			"as":           "cashier",
		}}},
		{{Key: "$unwind", Value: bson.M{"path": "$cashier", "preserveNullAndEmptyArrays": true}}},
		{{Key: "$project", Value: bson.M{
			"_id":          0,
			"cashier_id":   "$_id.cashier_id",
			"cashier_name": bson.M{"$ifNull": []interface{}{"$cashier.name", ""}},
			"year":         "$_id.year",
			"order_count":  1,
			"total_sales":  1,
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "year", Value: 1}}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, cashier_errors.ErrGetYearlyCashier.WithInternal(err)
	}
	var results []*salemongo.CashierYearlySaleRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, cashier_errors.ErrGetYearlyCashier.WithInternal(err)
	}
	return results, nil
}
