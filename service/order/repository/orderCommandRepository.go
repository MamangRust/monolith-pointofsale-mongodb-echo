package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type orderCommandRepository struct {
	coll *mongodriver.Collection
}

func NewOrderCommandRepository(coll *mongodriver.Collection) OrderCommandRepository {
	return &orderCommandRepository{coll: coll}
}

func ptrTime(t time.Time) *time.Time { return &t }

func (r *orderCommandRepository) DeleteOrder(ctx context.Context, orderID int) error {
	now := time.Now()
	_, err := r.coll.UpdateOne(ctx, bson.M{"order_id": orderID}, bson.M{"$set": bson.M{"deleted_at": now}})
	if err != nil {
		return order_errors.ErrDeleteOrderPermanent
	}
	return nil
}

func (r *orderCommandRepository) CreateOrder(ctx context.Context, request *requests.CreateOrderRecordRequest) (*salemongo.Order, error) {
	orderID, err := salemongo.NextID(ctx, r.coll.Database(), "orders")
	if err != nil {
		return nil, order_errors.ErrCreateOrder
	}

	order := &salemongo.Order{
		OrderID:    orderID,
		MerchantID: int32(request.MerchantID),
		CashierID:  int32(request.CashierID),
		TotalPrice: int64(request.TotalPrice),
		CreatedAt:  ptrTime(time.Now()),
		UpdatedAt:  ptrTime(time.Now()),
	}
	result, err := r.coll.InsertOne(ctx, order)
	if err != nil {
		return nil, order_errors.ErrCreateOrder
	}
	order.ID = result.InsertedID.(primitive.ObjectID)
	return order, nil
}

func (r *orderCommandRepository) FindAllTrashed(ctx context.Context) ([]*salemongo.Order, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return nil, order_errors.ErrFindAllTrashed
	}
	defer cursor.Close(ctx)
	var orders []*salemongo.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, order_errors.ErrFindAllTrashed
	}
	return orders, nil
}

func (r *orderCommandRepository) UpdateOrder(ctx context.Context, request *requests.UpdateOrderRecordRequest) (*salemongo.Order, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"order_id": request.OrderID, "deleted_at": nil},
		bson.M{"$set": bson.M{"total_price": int64(request.TotalPrice), "updated_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var order salemongo.Order
	if err := result.Decode(&order); err != nil {
		return nil, order_errors.ErrUpdateOrder
	}
	return &order, nil
}

func (r *orderCommandRepository) TrashedOrder(ctx context.Context, orderID int) (*salemongo.Order, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"order_id": orderID, "deleted_at": nil},
		bson.M{"$set": bson.M{"deleted_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var order salemongo.Order
	if err := result.Decode(&order); err != nil {
		return nil, order_errors.ErrTrashedOrder.WithInternal(err)
	}
	return &order, nil
}

func (r *orderCommandRepository) RestoreOrder(ctx context.Context, orderID int) (*salemongo.Order, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"order_id": orderID, "deleted_at": bson.M{"$ne": nil}},
		bson.M{"$unset": bson.M{"deleted_at": ""}, "$set": bson.M{"updated_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var order salemongo.Order
	if err := result.Decode(&order); err != nil {
		return nil, order_errors.ErrRestoreOrderNotFound
	}
	return &order, nil
}

func (r *orderCommandRepository) DeleteOrderPermanent(ctx context.Context, orderID int) (bool, error) {
	result, err := r.coll.DeleteOne(ctx, bson.M{"order_id": orderID, "deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, order_errors.ErrDeleteOrderPermanent.WithInternal(err)
	}
	if result.DeletedCount == 0 {
		return false, order_errors.ErrDeleteOrderPermanentNotFound
	}
	return true, nil
}

func (r *orderCommandRepository) DeleteAllOrderPermanent(ctx context.Context) (bool, error) {
	_, err := r.coll.DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, order_errors.ErrDeleteAllOrderPermanent.WithInternal(err)
	}
	return true, nil
}
