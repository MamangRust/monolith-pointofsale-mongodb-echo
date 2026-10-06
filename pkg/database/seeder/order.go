package seeder

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

type orderSeeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewOrderSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *orderSeeder {
	return &orderSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *orderSeeder) Seed() error {
	merchants, err := findAll[*salemongo.Merchant](r.ctx, r.db.Collection("merchants"), bson.M{"deleted_at": nil}, 20)
	if err != nil {
		r.logger.Error("Failed to get merchants", zap.Error(err))
		return err
	}

	cashiers, err := findAll[*salemongo.Cashier](r.ctx, r.db.Collection("cashiers"), bson.M{"deleted_at": nil}, 20)
	if err != nil {
		r.logger.Error("Failed to get cashiers", zap.Error(err))
		return err
	}

	if len(merchants) == 0 || len(cashiers) == 0 {
		r.logger.Error("No merchants or cashiers found, skipping order seeding")
		return nil
	}

	ordersColl := r.db.Collection("orders")
	itemsColl := r.db.Collection("order_items")
	now := time.Now().UTC()

	for i := 0; i < 10; i++ {
		merchant := merchants[rand.Intn(len(merchants))]
		cashier := cashiers[rand.Intn(len(cashiers))]
		totalPrice := int32(rand.Intn(500000) + 50000)

		orderID, err := salemongo.NextID(r.ctx, r.db, "orders")
		if err != nil {
			r.logger.Error("Failed to generate order id", zap.Error(err))
			return err
		}

		order := &salemongo.Order{
			OrderID:    orderID,
			MerchantID: merchant.MerchantID,
			CashierID:  cashier.CashierID,
			TotalPrice: int64(totalPrice),
			CreatedAt:  &now,
			UpdatedAt:  &now,
		}

		if _, err := ordersColl.InsertOne(r.ctx, order); err != nil {
			r.logger.Error("Failed to create order", zap.Error(err))
			return err
		}

		products, err := findAll[*salemongo.Product](r.ctx, r.db.Collection("products"),
			bson.M{"merchant_id": merchant.MerchantID, "deleted_at": nil}, 10)
		if err != nil {
			r.logger.Error("Failed to get products", zap.Error(err))
			return err
		}

		if len(products) == 0 {
			r.logger.Debug("No products found for merchant", zap.Int32("merchant_id", merchant.MerchantID))
			continue
		}

		for j := 0; j < rand.Intn(5)+1; j++ {
			product := products[rand.Intn(len(products))]
			quantity := int32(rand.Intn(5) + 1)
			price := int64(product.Price) * int64(quantity)

			orderItemID, err := salemongo.NextID(r.ctx, r.db, "order_items")
			if err != nil {
				r.logger.Error("Failed to generate order item id", zap.Error(err))
				return err
			}

			orderItem := &salemongo.OrderItem{
				OrderItemID: orderItemID,
				OrderID:     orderID,
				ProductID:   product.ProductID,
				Quantity:    quantity,
				Price:       price,
				CreatedAt:   &now,
				UpdatedAt:   &now,
			}

			if _, err := itemsColl.InsertOne(r.ctx, orderItem); err != nil {
				r.logger.Error("Failed to create order item", zap.Error(err))
				return err
			}
		}
	}

	r.logger.Info("Order seeding completed successfully.", zap.Int("count", 10))
	return nil
}
