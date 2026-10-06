package seeder

import (
	"context"
	"fmt"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

type cashierSeeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewCashierSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *cashierSeeder {
	return &cashierSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *cashierSeeder) Seed() error {
	merchants, err := findAll[*salemongo.Merchant](r.ctx, r.db.Collection("merchants"), bson.M{"deleted_at": nil}, 20)
	if err != nil {
		r.logger.Error("Failed to fetch merchants:", zap.Any("err", err))
		return err
	}

	users, err := findAll[*salemongo.User](r.ctx, r.db.Collection("users"), bson.M{"deleted_at": nil}, 20)
	if err != nil {
		r.logger.Error("Failed to fetch users:", zap.Any("error", err))
		return err
	}

	if len(merchants) == 0 || len(users) == 0 {
		r.logger.Error("Merchants or Users not found. Seed operation aborted.")
		return fmt.Errorf("no merchants or users found")
	}

	coll := r.db.Collection("cashiers")
	now := time.Now().UTC()

	for i := 1; i <= 10; i++ {
		merchant := merchants[rand.Intn(len(merchants))]
		user := users[rand.Intn(len(users))]

		cashierID, err := salemongo.NextID(r.ctx, r.db, "cashiers")
		if err != nil {
			r.logger.Error("Failed to generate cashier id:", zap.Any("error", err))
			return err
		}

		cashier := &salemongo.Cashier{
			CashierID:  cashierID,
			MerchantID: merchant.MerchantID,
			UserID:     user.UserID,
			Name:       fmt.Sprintf("Cashier %d", i),
			CreatedAt:  &now,
			UpdatedAt:  &now,
		}

		if _, err := coll.InsertOne(r.ctx, cashier); err != nil {
			r.logger.Error("Failed to create cashier:", zap.Any("error", err))
			return err
		}
	}

	r.logger.Info("Cashier seeding completed successfully.", zap.Int("count", 10))
	return nil
}
