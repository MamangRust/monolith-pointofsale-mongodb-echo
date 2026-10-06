package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/transaction_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type transactonStatsRepository struct {
	db *mongo.Database
}

func NewTransactionStatsRepository(db *mongo.Database) TransactionStatsRepository {
	return &transactonStatsRepository{
		db: db,
	}
}

// --- shared aggregation helpers -------------------------------------------------

// monthRange returns the [start, end) bounds of the given month.
func monthRange(year, month int) (time.Time, time.Time) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	return start, start.AddDate(0, 1, 0)
}

// yearRange returns the [start, end) bounds of the given year.
func yearRange(year int) (time.Time, time.Time) {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	return start, start.AddDate(1, 0, 0)
}

func bsonString(expr string) bson.M { return bson.M{"$toString": expr} }
func bsonInt(expr string) bson.M    { return bson.M{"$toInt": expr} }
func bsonDouble(expr string) bson.M { return bson.M{"$toDouble": expr} }

func sortKeys(monthly bool) bson.D {
	if monthly {
		return bson.D{{Key: "year", Value: 1}, {Key: "month", Value: 1}}
	}
	return bson.D{{Key: "year", Value: 1}}
}

// amountPipeline aggregates successful transaction amounts, grouped by month or year.
func amountPipeline(match bson.M, monthly bool) mongo.Pipeline {
	id := bson.M{"year": bson.M{"$year": "$created_at"}}
	project := bson.M{
		"_id":           0,
		"year":          bsonString("$_id.year"),
		"total_amount":  bsonInt("$total_amount"),
		"total_success": bsonInt("$total_success"),
	}
	if monthly {
		id["month"] = bson.M{"$month": "$created_at"}
		project["month"] = bsonString("$_id.month")
	}

	return mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$group", Value: bson.M{
			"_id":           id,
			"total_amount":  bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$payment_status", "completed"}}, "$amount", 0}}},
			"total_success": bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$payment_status", "completed"}}, 1, 0}}},
		}}},
		{{Key: "$project", Value: project}},
		{{Key: "$sort", Value: sortKeys(monthly)}},
	}
}

// failedPipeline aggregates failed transaction counts, grouped by month or year.
func failedPipeline(match bson.M, monthly bool) mongo.Pipeline {
	id := bson.M{"year": bson.M{"$year": "$created_at"}}
	project := bson.M{
		"_id":          0,
		"year":         bsonString("$_id.year"),
		"total_failed": bsonInt("$total_failed"),
		"total_amount": bsonInt("$total_amount"),
	}
	if monthly {
		id["month"] = bson.M{"$month": "$created_at"}
		project["month"] = bsonString("$_id.month")
	}

	return mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$group", Value: bson.M{
			"_id":          id,
			"total_failed": bson.M{"$sum": 1},
			"total_amount": bson.M{"$sum": "$amount"},
		}}},
		{{Key: "$project", Value: project}},
		{{Key: "$sort", Value: sortKeys(monthly)}},
	}
}

// methodPipeline aggregates payment-method usage, grouped by month or year.
func methodPipeline(match bson.M, monthly bool) mongo.Pipeline {
	id := bson.M{
		"payment_method": "$payment_method",
		"year":           bson.M{"$year": "$created_at"},
	}
	project := bson.M{
		"_id":                0,
		"year":               bsonString("$_id.year"),
		"payment_method":     "$_id.payment_method",
		"total_transactions": "$total_transactions",
		"total_amount":       bsonDouble("$total_amount"),
	}
	if monthly {
		id["month"] = bson.M{"$month": "$created_at"}
		delete(project, "year")
		project["month"] = bsonString("$_id.month")
	}

	return mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$group", Value: bson.M{
			"_id":                id,
			"total_transactions": bson.M{"$sum": 1},
			"total_amount":       bson.M{"$sum": "$amount"},
		}}},
		{{Key: "$project", Value: project}},
		{{Key: "$sort", Value: sortKeys(monthly)}},
	}
}

// --- amount stats ---------------------------------------------------------------

func (r *transactonStatsRepository) GetMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*salemongo.TransactionMonthlyAmountSuccessRow, error) {
	start, end := monthRange(req.Year, req.Month)
	match := bson.M{"created_at": bson.M{"$gte": start, "$lt": end}, "deleted_at": nil}

	cursor, err := r.db.Collection("transactions").Aggregate(ctx, amountPipeline(match, true))
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountSuccess.WithInternal(err)
	}
	var results []*salemongo.TransactionMonthlyAmountSuccessRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountSuccess.WithInternal(err)
	}
	return results, nil
}

func (r *transactonStatsRepository) GetYearlyAmountSuccess(ctx context.Context, year int) ([]*salemongo.TransactionYearlyAmountSuccessRow, error) {
	start, end := yearRange(year)
	match := bson.M{"created_at": bson.M{"$gte": start, "$lt": end}, "deleted_at": nil}

	cursor, err := r.db.Collection("transactions").Aggregate(ctx, amountPipeline(match, false))
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountSuccess.WithInternal(err)
	}
	var results []*salemongo.TransactionYearlyAmountSuccessRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountSuccess.WithInternal(err)
	}
	return results, nil
}

func (r *transactonStatsRepository) GetMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*salemongo.TransactionMonthlyAmountFailedRow, error) {
	start, end := monthRange(req.Year, req.Month)
	match := bson.M{"created_at": bson.M{"$gte": start, "$lt": end}, "deleted_at": nil, "payment_status": "failed"}

	cursor, err := r.db.Collection("transactions").Aggregate(ctx, failedPipeline(match, true))
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountFailed.WithInternal(err)
	}
	var results []*salemongo.TransactionMonthlyAmountFailedRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountFailed.WithInternal(err)
	}
	return results, nil
}

func (r *transactonStatsRepository) GetYearlyAmountFailed(ctx context.Context, year int) ([]*salemongo.TransactionYearlyAmountFailedRow, error) {
	start, end := yearRange(year)
	match := bson.M{"created_at": bson.M{"$gte": start, "$lt": end}, "deleted_at": nil, "payment_status": "failed"}

	cursor, err := r.db.Collection("transactions").Aggregate(ctx, failedPipeline(match, false))
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountFailed.WithInternal(err)
	}
	var results []*salemongo.TransactionYearlyAmountFailedRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountFailed.WithInternal(err)
	}
	return results, nil
}

// --- method stats ---------------------------------------------------------------

func (r *transactonStatsRepository) GetMonthlyTransactionMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*salemongo.TransactionMonthlyMethodSuccessRow, error) {
	start, end := monthRange(req.Year, req.Month)
	match := bson.M{"created_at": bson.M{"$gte": start, "$lt": end}, "deleted_at": nil, "payment_status": "completed"}

	cursor, err := r.db.Collection("transactions").Aggregate(ctx, methodPipeline(match, true))
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethod.WithInternal(err)
	}
	var results []*salemongo.TransactionMonthlyMethodSuccessRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethod.WithInternal(err)
	}
	return results, nil
}

func (r *transactonStatsRepository) GetYearlyTransactionMethodSuccess(ctx context.Context, year int) ([]*salemongo.TransactionYearlyMethodSuccessRow, error) {
	start, end := yearRange(year)
	match := bson.M{"created_at": bson.M{"$gte": start, "$lt": end}, "deleted_at": nil, "payment_status": "completed"}

	cursor, err := r.db.Collection("transactions").Aggregate(ctx, methodPipeline(match, false))
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethod.WithInternal(err)
	}
	var results []*salemongo.TransactionYearlyMethodSuccessRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethod.WithInternal(err)
	}
	return results, nil
}

func (r *transactonStatsRepository) GetMonthlyTransactionMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*salemongo.TransactionMonthlyMethodFailedRow, error) {
	start, end := monthRange(req.Year, req.Month)
	match := bson.M{"created_at": bson.M{"$gte": start, "$lt": end}, "deleted_at": nil, "payment_status": "failed"}

	cursor, err := r.db.Collection("transactions").Aggregate(ctx, methodPipeline(match, true))
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethod.WithInternal(err)
	}
	var results []*salemongo.TransactionMonthlyMethodFailedRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethod.WithInternal(err)
	}
	return results, nil
}

func (r *transactonStatsRepository) GetYearlyTransactionMethodFailed(ctx context.Context, year int) ([]*salemongo.TransactionYearlyMethodFailedRow, error) {
	start, end := yearRange(year)
	match := bson.M{"created_at": bson.M{"$gte": start, "$lt": end}, "deleted_at": nil, "payment_status": "failed"}

	cursor, err := r.db.Collection("transactions").Aggregate(ctx, methodPipeline(match, false))
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethod.WithInternal(err)
	}
	var results []*salemongo.TransactionYearlyMethodFailedRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethod.WithInternal(err)
	}
	return results, nil
}
