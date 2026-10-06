package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type orderStatsRepository struct {
	db *mongo.Database
}

func NewOrderStatsRepository(db *mongo.Database) OrderStatsRepository {
	return &orderStatsRepository{db: db}
}

// revenueAggRow is the decoded shape of the yearly/monthly revenue aggregation.
type revenueAggRow struct {
	Year         string `bson:"year"`
	Month        string `bson:"month"`
	TotalRevenue int32  `bson:"total_revenue"`
}

// orderAggRow is the decoded shape of the yearly/monthly order aggregation.
type orderAggRow struct {
	Year               string  `bson:"year"`
	Month              string  `bson:"month"`
	OrderCount         int64   `bson:"order_count"`
	TotalRevenue       float64 `bson:"total_revenue"`
	TotalItemsSold     int64   `bson:"total_items_sold"`
	ActiveCashiers     int64   `bson:"active_cashiers"`
	UniqueProductsSold int64   `bson:"unique_products_sold"`
}

// revenuePipeline aggregates SUM(total_price) over the orders collection.
func revenuePipeline(match bson.M, monthly bool) mongo.Pipeline {
	id := bson.M{"year": bson.M{"$year": "$created_at"}}
	project := bson.M{
		"_id":           0,
		"year":          bson.M{"$toString": "$_id.year"},
		"total_revenue": bson.M{"$toInt": "$total_revenue"},
	}
	sort := bson.D{{Key: "year", Value: 1}}
	if monthly {
		id["month"] = bson.M{"$month": "$created_at"}
		project["month"] = bson.M{"$toString": "$_id.month"}
		sort = bson.D{{Key: "year", Value: 1}, {Key: "month", Value: 1}}
	}

	return mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$group", Value: bson.M{"_id": id, "total_revenue": bson.M{"$sum": "$total_price"}}}},
		{{Key: "$project", Value: project}},
		{{Key: "$sort", Value: sort}},
	}
}

// orderPipeline aggregates order counts, revenue and item/cashier/product
// distinct counters, joining order_items with $lookup (the SQL LEFT JOIN).
func orderPipeline(match bson.M, monthly bool) mongo.Pipeline {
	id := bson.M{"year": bson.M{"$year": "$created_at"}}
	project := bson.M{
		"_id":                  0,
		"year":                 bson.M{"$toString": "$_id.year"},
		"order_count":          bson.M{"$size": "$order_ids"},
		"total_revenue":        bson.M{"$toDouble": "$total_revenue"},
		"total_items_sold":     "$total_items_sold",
		"active_cashiers":      bson.M{"$size": "$cashier_ids"},
		"unique_products_sold": bson.M{"$size": "$product_ids"},
	}
	sort := bson.D{{Key: "year", Value: 1}}
	if monthly {
		id = bson.M{"$dateToString": bson.M{"format": "%Y-%m", "date": "$created_at"}}
		project = bson.M{
			"_id":              0,
			"month":            "$_id",
			"order_count":      bson.M{"$size": "$order_ids"},
			"total_revenue":    bson.M{"$toDouble": "$total_revenue"},
			"total_items_sold": "$total_items_sold",
		}
		sort = bson.D{{Key: "month", Value: 1}}
	}

	return mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$lookup", Value: bson.M{
			"from": "order_items",
			"let":  bson.M{"oid": "$order_id"},
			"pipeline": bson.A{
				bson.M{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$order_id", "$$oid"}}}},
				bson.M{"$match": bson.M{"deleted_at": nil}},
			},
			"as": "items",
		}}},
		{{Key: "$unwind", Value: bson.M{"path": "$items", "preserveNullAndEmptyArrays": true}}},
		{{Key: "$group", Value: bson.M{
			"_id":              id,
			"order_ids":        bson.M{"$addToSet": "$order_id"},
			"cashier_ids":      bson.M{"$addToSet": "$cashier_id"},
			"product_ids":      bson.M{"$addToSet": "$items.product_id"},
			"total_revenue":    bson.M{"$sum": "$total_price"},
			"total_items_sold": bson.M{"$sum": "$items.quantity"},
		}}},
		{{Key: "$project", Value: project}},
		{{Key: "$sort", Value: sort}},
	}
}

func (r *orderStatsRepository) GetMonthlyTotalRevenue(ctx context.Context, req *requests.MonthTotalRevenue) ([]*salemongo.MonthlyTotalRevenueRow, error) {
	monthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	monthEnd := monthStart.AddDate(0, 1, 0)
	match := bson.M{"created_at": bson.M{"$gte": monthStart, "$lt": monthEnd}, "deleted_at": nil}

	cursor, err := r.db.Collection("orders").Aggregate(ctx, revenuePipeline(match, true))
	if err != nil {
		return nil, order_errors.ErrGetMonthlyTotalRevenue
	}
	var rows []revenueAggRow
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, order_errors.ErrGetMonthlyTotalRevenue
	}

	var res []*salemongo.MonthlyTotalRevenueRow
	for _, row := range rows {
		res = append(res, &salemongo.MonthlyTotalRevenueRow{
			Year:         row.Year,
			Month:        row.Month,
			TotalRevenue: row.TotalRevenue,
		})
	}
	return res, nil
}

func (r *orderStatsRepository) GetYearlyTotalRevenue(ctx context.Context, year int) ([]*salemongo.YearlyTotalRevenueRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	match := bson.M{"created_at": bson.M{"$gte": yearStart, "$lt": yearStart.AddDate(1, 0, 0)}, "deleted_at": nil}

	cursor, err := r.db.Collection("orders").Aggregate(ctx, revenuePipeline(match, false))
	if err != nil {
		return nil, order_errors.ErrGetYearlyTotalRevenue
	}
	var rows []revenueAggRow
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, order_errors.ErrGetYearlyTotalRevenue
	}

	if len(rows) == 0 {
		return []*salemongo.YearlyTotalRevenueRow{{Year: strconv.Itoa(year)}}, nil
	}
	res := make([]*salemongo.YearlyTotalRevenueRow, 0, len(rows))
	for _, row := range rows {
		res = append(res, &salemongo.YearlyTotalRevenueRow{Year: row.Year, TotalRevenue: row.TotalRevenue})
	}
	return res, nil
}

func (r *orderStatsRepository) GetMonthlyOrder(ctx context.Context, year int) ([]*salemongo.MonthlyOrderRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	match := bson.M{"created_at": bson.M{"$gte": yearStart, "$lt": yearStart.AddDate(1, 0, 0)}, "deleted_at": nil}

	cursor, err := r.db.Collection("orders").Aggregate(ctx, orderPipeline(match, true))
	if err != nil {
		return nil, order_errors.ErrGetMonthlyOrder
	}
	var rows []orderAggRow
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, order_errors.ErrGetMonthlyOrder
	}

	res := make([]*salemongo.MonthlyOrderRow, 0, len(rows))
	for _, row := range rows {
		res = append(res, &salemongo.MonthlyOrderRow{
			Month:          row.Month,
			OrderCount:     row.OrderCount,
			TotalRevenue:   row.TotalRevenue,
			TotalItemsSold: row.TotalItemsSold,
		})
	}
	return res, nil
}

func (r *orderStatsRepository) GetYearlyOrder(ctx context.Context, year int) ([]*salemongo.YearlyOrderRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	match := bson.M{"created_at": bson.M{"$gte": yearStart, "$lt": yearStart.AddDate(1, 0, 0)}, "deleted_at": nil}

	cursor, err := r.db.Collection("orders").Aggregate(ctx, orderPipeline(match, false))
	if err != nil {
		return nil, order_errors.ErrGetYearlyOrder
	}
	var rows []orderAggRow
	if err := cursor.All(ctx, &rows); err != nil {
		return nil, order_errors.ErrGetYearlyOrder
	}

	if len(rows) == 0 {
		return []*salemongo.YearlyOrderRow{{Year: fmt.Sprintf("%d", year)}}, nil
	}
	res := make([]*salemongo.YearlyOrderRow, 0, len(rows))
	for _, row := range rows {
		res = append(res, &salemongo.YearlyOrderRow{
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
