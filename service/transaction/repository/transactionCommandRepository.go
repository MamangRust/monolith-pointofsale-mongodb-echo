package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/transaction_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type transactionCommandRepository struct {
	coll *mongodriver.Collection
}

func NewTransactionCommandRepository(coll *mongodriver.Collection) *transactionCommandRepository {
	return &transactionCommandRepository{coll: coll}
}

func (r *transactionCommandRepository) CreateTransaction(ctx context.Context, request *requests.CreateTransactionRequest) (*salemongo.Transaction, error) {
	transactionID, err := salemongo.NextID(ctx, r.coll.Database(), "transactions")
	if err != nil {
		return nil, transaction_errors.ErrCreateTransaction
	}

	now := time.Now()
	paymentStatus := derefStr(request.PaymentStatus)
	if paymentStatus == "" {
		paymentStatus = "success"
	}
	t := &salemongo.Transaction{
		TransactionID: transactionID,
		OrderID:       int32(request.OrderID),
		MerchantID:    int32(request.MerchantID),
		CashierID:     int32(request.CashierID),
		PaymentMethod: request.PaymentMethod,
		TotalPrice:    int64(request.Amount),
		Amount:        int32(request.Amount),
		PaymentStatus: paymentStatus,
		Status:        paymentStatus,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}
	result, err := r.coll.InsertOne(ctx, t)
	if err != nil {
		return nil, transaction_errors.ErrCreateTransaction
	}
	t.ID = result.InsertedID.(primitive.ObjectID)
	return t, nil
}

func (r *transactionCommandRepository) UpdateTransaction(ctx context.Context, request *requests.UpdateTransactionRequest) (*salemongo.Transaction, error) {
	now := time.Now()
	setFields := bson.M{
		"merchant_id":    int32(request.MerchantID),
		"payment_method": request.PaymentMethod,
		"total_price":    int64(request.Amount),
		"amount":         int32(request.Amount),
		"order_id":       int32(request.OrderID),
		"updated_at":     now,
	}
	if request.PaymentStatus != nil {
		setFields["payment_status"] = *request.PaymentStatus
		setFields["status"] = *request.PaymentStatus
	}
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"transaction_id": *request.TransactionID, "deleted_at": nil},
		bson.M{"$set": setFields},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var t salemongo.Transaction
	if err := result.Decode(&t); err != nil {
		return nil, transaction_errors.ErrUpdateTransaction
	}
	return &t, nil
}

func (r *transactionCommandRepository) TrashTransaction(ctx context.Context, transaction_id int) (*salemongo.Transaction, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"transaction_id": transaction_id, "deleted_at": nil},
		bson.M{"$set": bson.M{"deleted_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var t salemongo.Transaction
	if err := result.Decode(&t); err != nil {
		return nil, transaction_errors.ErrTrashTransaction
	}
	return &t, nil
}

func (r *transactionCommandRepository) RestoreTransaction(ctx context.Context, transaction_id int) (*salemongo.Transaction, error) {
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"transaction_id": transaction_id, "deleted_at": bson.M{"$ne": nil}},
		bson.M{"$unset": bson.M{"deleted_at": ""}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var t salemongo.Transaction
	if err := result.Decode(&t); err != nil {
		return nil, transaction_errors.ErrRestoreTransaction
	}
	return &t, nil
}

func (r *transactionCommandRepository) DeleteTransactionPermanently(ctx context.Context, transaction_id int) (bool, error) {
	result, err := r.coll.DeleteOne(ctx, bson.M{"transaction_id": transaction_id})
	if err != nil {
		return false, transaction_errors.ErrDeleteTransactionPermanently
	}
	if result.DeletedCount == 0 {
		return false, transaction_errors.ErrDeleteTransactionPermanently
	}
	return true, nil
}

func (r *transactionCommandRepository) RestoreAllTransactions(ctx context.Context) (bool, error) {
	_, err := r.coll.UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$unset": bson.M{"deleted_at": ""}})
	if err != nil {
		return false, transaction_errors.ErrRestoreAllTransactions
	}
	return true, nil
}

func (r *transactionCommandRepository) DeleteAllTransactionPermanent(ctx context.Context) (bool, error) {
	_, err := r.coll.DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, transaction_errors.ErrDeleteAllTransactionPermanent
	}
	return true, nil
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
