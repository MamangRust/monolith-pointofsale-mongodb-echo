package repository

import (
	"context"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	orderitem_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/order_item_errors"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type orderItemQueryRepository struct {
	coll *mongodriver.Collection
}

func NewOrderItemQueryRepository(coll *mongodriver.Collection) OrderItemQueryRepository {
	return &orderItemQueryRepository{coll: coll}
}

func (r *orderItemQueryRepository) FindAllOrderItems(ctx context.Context, req *requests.FindAllOrderItems) ([]*OrderItemResult, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, orderitem_errors.ErrFindAllOrderItems
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "order_item_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, orderitem_errors.ErrFindAllOrderItems
	}
	defer cursor.Close(ctx)

	var items []salemongo.OrderItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, nil, orderitem_errors.ErrFindAllOrderItems
	}

	var results []*OrderItemResult
	for _, item := range items {
		cr := &OrderItemResult{
			OrderItemID: item.OrderItemID, OrderID: item.OrderID, ProductID: item.ProductID,
			Quantity: item.Quantity, Price: int64(item.Price), TotalCount: totalCount,
		}
		if item.CreatedAt != nil {
			cr.CreatedAt = item.CreatedAt
		}
		if item.UpdatedAt != nil {
			cr.UpdatedAt = item.UpdatedAt
		}
		results = append(results, cr)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *orderItemQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*OrderItemResultDeleteAt, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, orderitem_errors.ErrFindAllOrderItems
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "order_item_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, orderitem_errors.ErrFindAllOrderItems
	}
	defer cursor.Close(ctx)

	var items []salemongo.OrderItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, nil, orderitem_errors.ErrFindAllOrderItems
	}

	var results []*OrderItemResultDeleteAt
	for _, item := range items {
		cr := &OrderItemResultDeleteAt{
			OrderItemID: item.OrderItemID, OrderID: item.OrderID, ProductID: item.ProductID,
			Quantity: item.Quantity, Price: int64(item.Price), TotalCount: totalCount,
		}
		if item.CreatedAt != nil {
			cr.CreatedAt = item.CreatedAt
		}
		if item.UpdatedAt != nil {
			cr.UpdatedAt = item.UpdatedAt
		}
		results = append(results, cr)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *orderItemQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*OrderItemResultDeleteAt, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, orderitem_errors.ErrFindAllOrderItems
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "order_item_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, orderitem_errors.ErrFindAllOrderItems
	}
	defer cursor.Close(ctx)

	var items []salemongo.OrderItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, nil, orderitem_errors.ErrFindAllOrderItems
	}

	var results []*OrderItemResultDeleteAt
	for _, item := range items {
		cr := &OrderItemResultDeleteAt{
			OrderItemID: item.OrderItemID, OrderID: item.OrderID, ProductID: item.ProductID,
			Quantity: item.Quantity, Price: int64(item.Price), TotalCount: totalCount,
		}
		if item.CreatedAt != nil {
			cr.CreatedAt = item.CreatedAt
		}
		if item.UpdatedAt != nil {
			cr.UpdatedAt = item.UpdatedAt
		}
		if item.DeletedAt != nil {
			cr.DeletedAt = item.DeletedAt
		}
		results = append(results, cr)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *orderItemQueryRepository) FindOrderItemByOrder(ctx context.Context, orderID int) ([]*salemongo.OrderItem, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"order_id": orderID})
	if err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder
	}
	defer cursor.Close(ctx)

	var items []salemongo.OrderItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder
	}
	var result []*salemongo.OrderItem
	for i := range items {
		result = append(result, &items[i])
	}
	return result, nil
}
