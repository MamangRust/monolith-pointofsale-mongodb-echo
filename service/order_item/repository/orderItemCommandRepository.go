package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	orderitem_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/order_item_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type orderItemCommandRepository struct {
	coll *mongodriver.Collection
}

func NewOrderItemCommandRepository(coll *mongodriver.Collection) OrderItemCommandRepository {
	return &orderItemCommandRepository{coll: coll}
}

func (r *orderItemCommandRepository) CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*salemongo.OrderItem, error) {
	orderItemID, err := salemongo.NextID(ctx, r.coll.Database(), "order_items")
	if err != nil {
		return nil, orderitem_errors.ErrCreateOrderItem
	}

	now := time.Now()
	item := &salemongo.OrderItem{
		OrderItemID: orderItemID,
		OrderID:     int32(req.OrderID),
		ProductID:   int32(req.ProductID),
		Quantity:    int32(req.Quantity),
		Price:       int64(req.Price),
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}
	result, err := r.coll.InsertOne(ctx, item)
	if err != nil {
		return nil, orderitem_errors.ErrCreateOrderItem
	}
	item.ID = result.InsertedID.(primitive.ObjectID)
	return item, nil
}

func (r *orderItemCommandRepository) UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*salemongo.OrderItem, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"order_item_id": req.OrderItemID},
		bson.M{"$set": bson.M{
			"product_id": int32(req.ProductID),
			"quantity":   int32(req.Quantity),
			"price":      int64(req.Price),
			"updated_at": now,
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var item salemongo.OrderItem
	if err := result.Decode(&item); err != nil {
		return nil, orderitem_errors.ErrUpdateOrderItem
	}
	return &item, nil
}

// DeleteOrderItem removes one order item by its own id (the rollback path of
// order-service). The filter must target order_item_id — not order_id — so a
// rollback actually removes the item that was created.
func (r *orderItemCommandRepository) DeleteOrderItem(ctx context.Context, orderItemID int) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"order_item_id": orderItemID})
	if err != nil {
		return orderitem_errors.ErrDeleteOrderItemPermanent
	}
	return nil
}
