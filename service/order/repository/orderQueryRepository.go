package repository

import (
	"context"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type orderQueryRepository struct {
	coll *mongodriver.Collection
}

func NewOrderQueryRepository(coll *mongodriver.Collection) OrderQueryRepository {
	return &orderQueryRepository{coll: coll}
}

func (r *orderQueryRepository) FindAllOrders(ctx context.Context, req *requests.FindAllOrders) ([]*OrderResult, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"order_id": bson.M{"$regex": req.Search}},
		}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, order_errors.ErrFindAllOrders
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "order_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, order_errors.ErrFindAllOrders
	}
	defer cursor.Close(ctx)

	var orders []salemongo.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, nil, order_errors.ErrFindAllOrders
	}

	var results []*OrderResult
	for _, o := range orders {
		or := &OrderResult{
			OrderID:    o.OrderID,
			MerchantID: o.MerchantID,
			CashierID:  o.CashierID,
			TotalPrice: o.TotalPrice,
			TotalCount: totalCount,
		}
		if o.CreatedAt != nil {
			or.CreatedAt = o.CreatedAt
		}
		if o.UpdatedAt != nil {
			or.UpdatedAt = o.UpdatedAt
		}
		results = append(results, or)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *orderQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllOrders) ([]*OrderResultDeleteAt, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"order_id": bson.M{"$regex": req.Search}},
		}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, order_errors.ErrFindByActive
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "order_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, order_errors.ErrFindByActive
	}
	defer cursor.Close(ctx)

	var orders []salemongo.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, nil, order_errors.ErrFindByActive
	}

	var results []*OrderResultDeleteAt
	for _, o := range orders {
		or := &OrderResultDeleteAt{
			OrderID:    o.OrderID,
			MerchantID: o.MerchantID,
			CashierID:  o.CashierID,
			TotalPrice: o.TotalPrice,
			TotalCount: totalCount,
		}
		if o.CreatedAt != nil {
			or.CreatedAt = o.CreatedAt
		}
		if o.UpdatedAt != nil {
			or.UpdatedAt = o.UpdatedAt
		}
		results = append(results, or)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *orderQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllOrders) ([]*OrderResultDeleteAt, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"order_id": bson.M{"$regex": req.Search}},
		}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, order_errors.ErrFindByTrashed
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "order_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, order_errors.ErrFindByTrashed
	}
	defer cursor.Close(ctx)

	var orders []salemongo.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, nil, order_errors.ErrFindByTrashed
	}

	var results []*OrderResultDeleteAt
	for _, o := range orders {
		or := &OrderResultDeleteAt{
			OrderID:    o.OrderID,
			MerchantID: o.MerchantID,
			CashierID:  o.CashierID,
			TotalPrice: o.TotalPrice,
			TotalCount: totalCount,
		}
		if o.CreatedAt != nil {
			or.CreatedAt = o.CreatedAt
		}
		if o.UpdatedAt != nil {
			or.UpdatedAt = o.UpdatedAt
		}
		if o.DeletedAt != nil {
			or.DeletedAt = o.DeletedAt
		}
		results = append(results, or)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *orderQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllOrderMerchant) ([]*OrderResult, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil, "merchant_id": req.MerchantID}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"order_id": bson.M{"$regex": req.Search}},
		}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, order_errors.ErrFindByMerchant
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "order_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, order_errors.ErrFindByMerchant
	}
	defer cursor.Close(ctx)

	var orders []salemongo.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, nil, order_errors.ErrFindByMerchant
	}

	var results []*OrderResult
	for _, o := range orders {
		or := &OrderResult{
			OrderID:    o.OrderID,
			MerchantID: o.MerchantID,
			CashierID:  o.CashierID,
			TotalPrice: o.TotalPrice,
			TotalCount: totalCount,
		}
		if o.CreatedAt != nil {
			or.CreatedAt = o.CreatedAt
		}
		if o.UpdatedAt != nil {
			or.UpdatedAt = o.UpdatedAt
		}
		results = append(results, or)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *orderQueryRepository) FindById(ctx context.Context, orderID int) (*salemongo.Order, error) {
	var order salemongo.Order
	err := r.coll.FindOne(ctx, bson.M{"order_id": orderID, "deleted_at": nil}).Decode(&order)
	if err != nil {
		return nil, order_errors.ErrFindById
	}
	return &order, nil
}

func (r *orderQueryRepository) FindByTrashedId(ctx context.Context, orderID int) (*salemongo.Order, error) {
	var order salemongo.Order
	err := r.coll.FindOne(ctx, bson.M{"order_id": orderID, "deleted_at": bson.M{"$ne": nil}}).Decode(&order)
	if err != nil {
		return nil, order_errors.ErrFindByTrashedId
	}
	return &order, nil
}
