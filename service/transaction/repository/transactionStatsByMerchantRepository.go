package repository

import (
	"context"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/transaction_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type transactionStatsByMerchantRepository struct {
	db *mongo.Database
}

func NewTransactionStatsByMerchantRepository(db *mongo.Database) TransactionStatsByMerchantRepository {
	return &transactionStatsByMerchantRepository{
		db: db,
	}
}

func (r *transactionStatsByMerchantRepository) GetMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*salemongo.TransactionMonthlyAmountSuccessByMerchantRow, error) {
	start, end := monthRange(req.Year, req.Month)
	match := bson.M{
		"merchant_id": int32(req.MerchantID),
		"created_at":  bson.M{"$gte": start, "$lt": end},
		"deleted_at":  nil,
	}

	cursor, err := r.db.Collection("transactions").Aggregate(ctx, amountPipeline(match, true))
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountSuccessByMerchant.WithInternal(err)
	}
	var results []*salemongo.TransactionMonthlyAmountSuccessByMerchantRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountSuccessByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*salemongo.TransactionYearlyAmountSuccessByMerchantRow, error) {
	start, end := yearRange(req.Year)
	match := bson.M{
		"merchant_id": int32(req.MerchantID),
		"created_at":  bson.M{"$gte": start, "$lt": end},
		"deleted_at":  nil,
	}

	cursor, err := r.db.Collection("transactions").Aggregate(ctx, amountPipeline(match, false))
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountSuccessByMerchant.WithInternal(err)
	}
	var results []*salemongo.TransactionYearlyAmountSuccessByMerchantRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountSuccessByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *transactionStatsByMerchantRepository) GetMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*salemongo.TransactionMonthlyAmountFailedByMerchantRow, error) {
	start, end := monthRange(req.Year, req.Month)
	match := bson.M{
		"merchant_id":    int32(req.MerchantID),
		"created_at":     bson.M{"$gte": start, "$lt": end},
		"deleted_at":     nil,
		"payment_status": "failed",
	}

	cursor, err := r.db.Collection("transactions").Aggregate(ctx, failedPipeline(match, true))
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountFailedByMerchant.WithInternal(err)
	}
	var results []*salemongo.TransactionMonthlyAmountFailedByMerchantRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountFailedByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*salemongo.TransactionYearlyAmountFailedByMerchantRow, error) {
	start, end := yearRange(req.Year)
	match := bson.M{
		"merchant_id":    int32(req.MerchantID),
		"created_at":     bson.M{"$gte": start, "$lt": end},
		"deleted_at":     nil,
		"payment_status": "failed",
	}

	cursor, err := r.db.Collection("transactions").Aggregate(ctx, failedPipeline(match, false))
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountFailedByMerchant.WithInternal(err)
	}
	var results []*salemongo.TransactionYearlyAmountFailedByMerchantRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountFailedByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *transactionStatsByMerchantRepository) GetMonthlyTransactionMethodByMerchantSuccess(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*salemongo.TransactionMonthlyMethodByMerchantSuccessRow, error) {
	start, end := monthRange(req.Year, req.Month)
	match := bson.M{
		"merchant_id":    int32(req.MerchantID),
		"created_at":     bson.M{"$gte": start, "$lt": end},
		"deleted_at":     nil,
		"payment_status": "completed",
	}

	cursor, err := r.db.Collection("transactions").Aggregate(ctx, methodPipeline(match, true))
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethodByMerchant.WithInternal(err)
	}
	var results []*salemongo.TransactionMonthlyMethodByMerchantSuccessRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethodByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyTransactionMethodByMerchantSuccess(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*salemongo.TransactionYearlyMethodByMerchantSuccessRow, error) {
	start, end := yearRange(req.Year)
	match := bson.M{
		"merchant_id":    int32(req.MerchantID),
		"created_at":     bson.M{"$gte": start, "$lt": end},
		"deleted_at":     nil,
		"payment_status": "completed",
	}

	cursor, err := r.db.Collection("transactions").Aggregate(ctx, methodPipeline(match, false))
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethodByMerchant.WithInternal(err)
	}
	var results []*salemongo.TransactionYearlyMethodByMerchantSuccessRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethodByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *transactionStatsByMerchantRepository) GetMonthlyTransactionMethodByMerchantFailed(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*salemongo.TransactionMonthlyMethodByMerchantFailedRow, error) {
	start, end := monthRange(req.Year, req.Month)
	match := bson.M{
		"merchant_id":    int32(req.MerchantID),
		"created_at":     bson.M{"$gte": start, "$lt": end},
		"deleted_at":     nil,
		"payment_status": "failed",
	}

	cursor, err := r.db.Collection("transactions").Aggregate(ctx, methodPipeline(match, true))
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethodByMerchant.WithInternal(err)
	}
	var results []*salemongo.TransactionMonthlyMethodByMerchantFailedRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethodByMerchant.WithInternal(err)
	}
	return results, nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyTransactionMethodByMerchantFailed(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*salemongo.TransactionYearlyMethodByMerchantFailedRow, error) {
	start, end := yearRange(req.Year)
	match := bson.M{
		"merchant_id":    int32(req.MerchantID),
		"created_at":     bson.M{"$gte": start, "$lt": end},
		"deleted_at":     nil,
		"payment_status": "failed",
	}

	cursor, err := r.db.Collection("transactions").Aggregate(ctx, methodPipeline(match, false))
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethodByMerchant.WithInternal(err)
	}
	var results []*salemongo.TransactionYearlyMethodByMerchantFailedRow
	if err := cursor.All(ctx, &results); err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethodByMerchant.WithInternal(err)
	}
	return results, nil
}
