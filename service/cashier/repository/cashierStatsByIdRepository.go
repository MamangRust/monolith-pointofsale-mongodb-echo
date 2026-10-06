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

type cashierStatsByIdRepository struct {
	db *mongo.Database
}

func NewCashierStatsByIdRepository(db *mongo.Database) CashierStatByIdRepository {
	return &cashierStatsByIdRepository{
		db: db,
	}
}

func (r *cashierStatsByIdRepository) GetMonthlyTotalSalesById(ctx context.Context, req *requests.MonthTotalSalesCashier) ([]*salemongo.CashierMonthlyTotalSalesRow, error) {
	coll := r.db.Collection("orders")
	currentStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	prevStart := currentStart.AddDate(0, -1, 0)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"cashier_id": int32(req.CashierID),
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
		return nil, cashier_errors.ErrGetMonthlyTotalSalesById.WithInternal(err)
	}
	var results []*salemongo.CashierMonthlyTotalSalesRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, cashier_errors.ErrGetMonthlyTotalSalesById.WithInternal(err)
	}
	return results, nil
}

func (r *cashierStatsByIdRepository) GetYearlyTotalSalesById(ctx context.Context, req *requests.YearTotalSalesCashier) ([]*salemongo.CashierYearlyTotalSalesRow, error) {
	coll := r.db.Collection("orders")
	start := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"cashier_id": int32(req.CashierID),
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
		return nil, cashier_errors.ErrGetYearlyTotalSalesById.WithInternal(err)
	}
	var results []*salemongo.CashierYearlyTotalSalesRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, cashier_errors.ErrGetYearlyTotalSalesById.WithInternal(err)
	}
	return results, nil
}

func (r *cashierStatsByIdRepository) GetMonthlyCashierById(ctx context.Context, req *requests.MonthCashierId) ([]*salemongo.CashierMonthlySaleByIdRow, error) {
	coll := r.db.Collection("orders")
	start := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"cashier_id": int32(req.CashierID),
			"created_at": bson.M{"$gte": start, "$lt": end},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":         bson.M{"month": bson.M{"$month": "$created_at"}},
			"order_count": bson.M{"$sum": 1},
			"total_sales": bson.M{"$sum": "$total_price"},
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "cashiers",
			"localField":   "cashier_id",
			"foreignField": "cashier_id",
			"as":           "cashier",
		}}},
		{{Key: "$unwind", Value: bson.M{"path": "$cashier", "preserveNullAndEmptyArrays": true}}},
		{{Key: "$project", Value: bson.M{
			"_id":          0,
			"cashier_id":   int32(req.CashierID),
			"cashier_name": bson.M{"$ifNull": []interface{}{"$cashier.name", ""}},
			"month":        "$_id.month",
			"order_count":  1,
			"total_sales":  1,
		}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, cashier_errors.ErrGetMonthlyCashierById.WithInternal(err)
	}
	var results []*salemongo.CashierMonthlySaleByIdRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, cashier_errors.ErrGetMonthlyCashierById.WithInternal(err)
	}
	return results, nil
}

func (r *cashierStatsByIdRepository) GetYearlyCashierById(ctx context.Context, req *requests.YearCashierId) ([]*salemongo.CashierYearlySaleRow, error) {
	coll := r.db.Collection("orders")
	start := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"cashier_id": int32(req.CashierID),
			"created_at": bson.M{"$gte": start, "$lt": end},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":         bson.M{"year": bson.M{"$year": "$created_at"}},
			"order_count": bson.M{"$sum": 1},
			"total_sales": bson.M{"$sum": "$total_price"},
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "cashiers",
			"localField":   "cashier_id",
			"foreignField": "cashier_id",
			"as":           "cashier",
		}}},
		{{Key: "$unwind", Value: bson.M{"path": "$cashier", "preserveNullAndEmptyArrays": true}}},
		{{Key: "$project", Value: bson.M{
			"_id":          0,
			"cashier_id":   int32(req.CashierID),
			"cashier_name": bson.M{"$ifNull": []interface{}{"$cashier.name", ""}},
			"year":         "$_id.year",
			"order_count":  1,
			"total_sales":  1,
		}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, cashier_errors.ErrGetYearlyCashierById.WithInternal(err)
	}
	var results []*salemongo.CashierYearlySaleRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, cashier_errors.ErrGetYearlyCashierById.WithInternal(err)
	}
	return results, nil
}
