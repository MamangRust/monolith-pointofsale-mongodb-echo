package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/transaction_errors"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type transactionQueryRepository struct {
	coll *mongodriver.Collection
}

func NewTransactionQueryRepository(coll *mongodriver.Collection) *transactionQueryRepository {
	return &transactionQueryRepository{coll: coll}
}

func (r *transactionQueryRepository) FindAllTransactions(ctx context.Context, req *requests.FindAllTransaction) ([]*TransactionResult, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["payment_method"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, transaction_errors.ErrFindAllTransactions
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, transaction_errors.ErrFindAllTransactions
	}
	defer cursor.Close(ctx)

	var txns []salemongo.Transaction
	if err := cursor.All(ctx, &txns); err != nil {
		return nil, nil, transaction_errors.ErrFindAllTransactions
	}

	var out []*TransactionResult
	for _, tx := range txns {
		cr := &TransactionResult{
			TransactionID: tx.TransactionID,
			OrderID:       tx.OrderID,
			MerchantID:    tx.MerchantID,
			PaymentMethod: tx.PaymentMethod,
			Amount:        int32(tx.TotalPrice),
			ChangeAmount:  gintptr32(0),
			PaymentStatus: tx.Status,
			TotalCount:    totalCount,
		}
		if tx.CreatedAt != nil {
			cr.CreatedAt = tx.CreatedAt
		}
		out = append(out, cr)
	}
	return out, ptrint(int(totalCount)), nil
}

func (r *transactionQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllTransaction) ([]*TransactionResultDeleteAt, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["payment_method"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, transaction_errors.ErrFindByActive
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, transaction_errors.ErrFindByActive
	}
	defer cursor.Close(ctx)

	var txns []salemongo.Transaction
	if err := cursor.All(ctx, &txns); err != nil {
		return nil, nil, transaction_errors.ErrFindByActive
	}

	var out []*TransactionResultDeleteAt
	for _, tx := range txns {
		cr := &TransactionResultDeleteAt{
			TransactionID: tx.TransactionID,
			OrderID:       tx.OrderID,
			MerchantID:    tx.MerchantID,
			PaymentMethod: tx.PaymentMethod,
			Amount:        int32(tx.TotalPrice),
			ChangeAmount:  gintptr32(0),
			PaymentStatus: tx.Status,
			TotalCount:    totalCount,
		}
		if tx.CreatedAt != nil {
			cr.CreatedAt = tx.CreatedAt
		}
		if tx.DeletedAt != nil {
			cr.DeletedAt = tx.DeletedAt
		}
		out = append(out, cr)
	}
	return out, ptrint(int(totalCount)), nil
}

func (r *transactionQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllTransaction) ([]*TransactionResultDeleteAt, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if req.Search != "" {
		filter["payment_method"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, transaction_errors.ErrFindByTrashed
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "deleted_at", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, transaction_errors.ErrFindByTrashed
	}
	defer cursor.Close(ctx)

	var txns []salemongo.Transaction
	if err := cursor.All(ctx, &txns); err != nil {
		return nil, nil, transaction_errors.ErrFindByTrashed
	}

	var out []*TransactionResultDeleteAt
	for _, tx := range txns {
		cr := &TransactionResultDeleteAt{
			TransactionID: tx.TransactionID,
			OrderID:       tx.OrderID,
			MerchantID:    tx.MerchantID,
			PaymentMethod: tx.PaymentMethod,
			Amount:        int32(tx.TotalPrice),
			ChangeAmount:  gintptr32(0),
			PaymentStatus: tx.Status,
			TotalCount:    totalCount,
		}
		if tx.CreatedAt != nil {
			cr.CreatedAt = tx.CreatedAt
		}
		if tx.DeletedAt != nil {
			cr.DeletedAt = tx.DeletedAt
		}
		out = append(out, cr)
	}
	return out, ptrint(int(totalCount)), nil
}

func (r *transactionQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*TransactionByMerchantResult, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil, "merchant_id": req.MerchantID}
	if req.Search != "" {
		filter["payment_method"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, transaction_errors.ErrFindByMerchant
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, transaction_errors.ErrFindByMerchant
	}
	defer cursor.Close(ctx)

	var txns []salemongo.Transaction
	if err := cursor.All(ctx, &txns); err != nil {
		return nil, nil, transaction_errors.ErrFindByMerchant
	}

	var out []*TransactionByMerchantResult
	for _, tx := range txns {
		cr := &TransactionByMerchantResult{
			TransactionID: tx.TransactionID,
			OrderID:       tx.OrderID,
			MerchantID:    tx.MerchantID,
			PaymentMethod: tx.PaymentMethod,
			Amount:        int32(tx.TotalPrice),
			ChangeAmount:  gintptr32(0),
			PaymentStatus: tx.Status,
			TotalCount:    totalCount,
		}
		if tx.CreatedAt != nil {
			cr.CreatedAt = tx.CreatedAt
		}
		out = append(out, cr)
	}
	return out, ptrint(int(totalCount)), nil
}

func (r *transactionQueryRepository) FindById(ctx context.Context, transaction_id int) (*salemongo.Transaction, error) {
	var t salemongo.Transaction
	err := r.coll.FindOne(ctx, bson.M{"transaction_id": transaction_id, "deleted_at": nil}).Decode(&t)
	if err != nil {
		return nil, transaction_errors.ErrFindById
	}
	return &t, nil
}

func (r *transactionQueryRepository) FindByOrderId(ctx context.Context, order_id int) (*salemongo.Transaction, error) {
	var t salemongo.Transaction
	err := r.coll.FindOne(ctx, bson.M{"order_id": order_id, "deleted_at": nil}).Decode(&t)
	if err != nil {
		return nil, transaction_errors.ErrFindByOrderId
	}
	return &t, nil
}

// parseTxTime converts a gRPC timestamp string to *time.Time.
func parseTxTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return nil
		}
	}
	return &t
}

func ptrint(v int) *int        { return &v }
func gintptr32(v int32) *int32 { return &v }
