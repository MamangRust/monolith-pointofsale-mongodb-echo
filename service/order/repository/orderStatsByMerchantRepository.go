package repository

import (
	"context"
	"fmt"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type orderStatsByMerchantRepository struct {
	db *mongo.Database
}

func NewOrderStatsByMerchantRepository(db *mongo.Database) OrderStatByMerchantRepository {
	return &orderStatsByMerchantRepository{db: db}
}

func (r *orderStatsByMerchantRepository) GetMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*salemongo.MonthlyTotalRevenueByMerchantRow, error) {
	monthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	match := bson.M{
		"merchant_id": int32(req.MerchantID),
		"created_at":  bson.M{"$gte": monthStart, "$lt": monthStart.AddDate(0, 1, 0)},
		"deleted_at":  nil,
	}

	cursor, err := r.db.Collection("orders").Aggregate(ctx, revenuePipeline(match, true))
	if err != nil {
		return nil, order_errors.ErrGetMonthlyTotalRevenueByMerchant
	}
	var rows []revenueAggRow
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, order_errors.ErrGetMonthlyTotalRevenueByMerchant
	}

	res := make([]*salemongo.MonthlyTotalRevenueByMerchantRow, 0, len(rows))
	for _, row := range rows {
		res = append(res, &salemongo.MonthlyTotalRevenueByMerchantRow{
			Year:         row.Year,
			Month:        row.Month,
			TotalRevenue: row.TotalRevenue,
		})
	}
	return res, nil
}

func (r *orderStatsByMerchantRepository) GetYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*salemongo.YearlyTotalRevenueByMerchantRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	match := bson.M{
		"merchant_id": int32(req.MerchantID),
		"created_at":  bson.M{"$gte": yearStart, "$lt": yearStart.AddDate(1, 0, 0)},
		"deleted_at":  nil,
	}

	cursor, err := r.db.Collection("orders").Aggregate(ctx, revenuePipeline(match, false))
	if err != nil {
		return nil, order_errors.ErrGetYearlyTotalRevenueByMerchant
	}
	var rows []revenueAggRow
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, order_errors.ErrGetYearlyTotalRevenueByMerchant
	}

	if len(rows) == 0 {
		return []*salemongo.YearlyTotalRevenueByMerchantRow{{Year: fmt.Sprintf("%d", req.Year)}}, nil
	}
	res := make([]*salemongo.YearlyTotalRevenueByMerchantRow, 0, len(rows))
	for _, row := range rows {
		res = append(res, &salemongo.YearlyTotalRevenueByMerchantRow{Year: row.Year, TotalRevenue: row.TotalRevenue})
	}
	return res, nil
}

func (r *orderStatsByMerchantRepository) GetMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*salemongo.MonthlyOrderByMerchantRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	match := bson.M{
		"merchant_id": int32(req.MerchantID),
		"created_at":  bson.M{"$gte": yearStart, "$lt": yearStart.AddDate(1, 0, 0)},
		"deleted_at":  nil,
	}

	cursor, err := r.db.Collection("orders").Aggregate(ctx, orderPipeline(match, true))
	if err != nil {
		return nil, order_errors.ErrGetMonthlyOrderByMerchant
	}
	var rows []orderAggRow
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, order_errors.ErrGetMonthlyOrderByMerchant
	}

	res := make([]*salemongo.MonthlyOrderByMerchantRow, 0, len(rows))
	for _, row := range rows {
		res = append(res, &salemongo.MonthlyOrderByMerchantRow{
			Month:          row.Month,
			OrderCount:     row.OrderCount,
			TotalRevenue:   row.TotalRevenue,
			TotalItemsSold: row.TotalItemsSold,
		})
	}
	return res, nil
}

func (r *orderStatsByMerchantRepository) GetYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*salemongo.YearlyOrderByMerchantRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	match := bson.M{
		"merchant_id": int32(req.MerchantID),
		"created_at":  bson.M{"$gte": yearStart, "$lt": yearStart.AddDate(1, 0, 0)},
		"deleted_at":  nil,
	}

	cursor, err := r.db.Collection("orders").Aggregate(ctx, orderPipeline(match, false))
	if err != nil {
		return nil, order_errors.ErrGetYearlyOrderByMerchant
	}
	var rows []orderAggRow
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, order_errors.ErrGetYearlyOrderByMerchant
	}

	if len(rows) == 0 {
		return []*salemongo.YearlyOrderByMerchantRow{{Year: fmt.Sprintf("%d", req.Year)}}, nil
	}
	res := make([]*salemongo.YearlyOrderByMerchantRow, 0, len(rows))
	for _, row := range rows {
		res = append(res, &salemongo.YearlyOrderByMerchantRow{
			Year:               row.Year,
			OrderCount:         row.OrderCount,
			TotalRevenue:       row.TotalRevenue,
			TotalItemsSold:     row.TotalItemsSold,
			ActiveCashiers:     row.ActiveCashiers,
			UniqueProductsSold: row.UniqueProductsSold,
		})
	}
	return res, nil
}
