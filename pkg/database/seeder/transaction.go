package seeder

import (
	"context"
	"math/rand"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type transactionSeeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewTransactionSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *transactionSeeder {
	return &transactionSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *transactionSeeder) Seed() error {
	orders, err := findAll[*salemongo.Order](r.ctx, r.db.Collection("orders"), bson.M{"deleted_at": nil}, 20)
	if err != nil {
		r.logger.Error("Failed to get transactions:", zap.Any("error", err))
		return err
	}

	merchants, err := findAll[*salemongo.Merchant](r.ctx, r.db.Collection("merchants"), bson.M{"deleted_at": nil}, 20)
	if err != nil {
		r.logger.Error("Failed to get transactions:", zap.Any("error", err))
		return err
	}

	if len(orders) == 0 || len(merchants) == 0 {
		r.logger.Error("No orders or merchants found, skipping transaction seeding")
		return nil
	}

	coll := r.db.Collection("transactions")
	now := time.Now().UTC()

	for i := 0; i < 10; i++ {
		selectedMerchant := merchants[rand.Intn(len(merchants))]
		selectedOrder := orders[rand.Intn(len(orders))]

		paymentMethod := "Credit Card"
		amount := int32(100 + i)
		changeAmount := int32(5 + i)
		paymentStatus := "completed"

		transactionID, err := salemongo.NextID(r.ctx, r.db, "transactions")
		if err != nil {
			r.logger.Error("Failed to generate transaction id:", zap.Any("error", err))
			return err
		}

		transaction := &salemongo.Transaction{
			TransactionID: transactionID,
			OrderID:       selectedOrder.OrderID,
			CashierID:     selectedOrder.CashierID,
			PaymentMethod: paymentMethod,
			Amount:        amount,
			ChangeAmount:  ptrInt32(changeAmount),
			PaymentStatus: paymentStatus,
			Status:        paymentStatus,
			TotalPrice:    selectedOrder.TotalPrice,
			MerchantID:    selectedMerchant.MerchantID,
			CreatedAt:     &now,
			UpdatedAt:     &now,
		}

		if _, err := coll.InsertOne(r.ctx, transaction); err != nil {
			r.logger.Error("Failed to create transaction:", zap.Any("error", err))
			return err
		}
	}

	r.logger.Info("Successfully seeded 10 transactions.", zap.Int("count", 10))
	return nil
}
