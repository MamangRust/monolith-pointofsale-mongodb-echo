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
)

type merchantSeeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *merchantSeeder {
	return &merchantSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *merchantSeeder) Seed() error {
	users, err := findAll[*salemongo.User](r.ctx, r.db.Collection("users"), bson.M{"deleted_at": nil}, 20)
	if err != nil {
		r.logger.Error("Failed to fetch merchants:", zap.Any("error", err))
		return err
	}

	if len(users) == 0 {
		r.logger.Error("No users found, skipping merchant seeding")
		return fmt.Errorf("no users found")
	}

	coll := r.db.Collection("merchants")
	now := time.Now().UTC()

	for i := 1; i <= 10; i++ {
		userID := users[i%len(users)].UserID

		merchantID, err := salemongo.NextID(r.ctx, r.db, "merchants")
		if err != nil {
			r.logger.Error("Failed to generate merchant id:", zap.Error(err))
			return err
		}

		merchant := &salemongo.Merchant{
			MerchantID:   merchantID,
			UserID:       userID,
			Name:         fmt.Sprintf("Toko %d", i),
			Description:  ptrString(fmt.Sprintf("Deskripsi untuk Toko %d", i)),
			Address:      ptrString(fmt.Sprintf("Jl. Toko %d", i)),
			ContactEmail: ptrString(fmt.Sprintf("toko%d@example.com", i)),
			ContactPhone: ptrString(fmt.Sprintf("0812345678%d", i)),
			Status:       "active",
			CreatedAt:    &now,
			UpdatedAt:    &now,
		}

		if _, err := coll.InsertOne(r.ctx, merchant); err != nil {
			r.logger.Error("Failed to create merchant:", zap.Error(err))
			return err
		}
	}

	r.logger.Info("Merchant seeding completed successfully.", zap.Int("count", 10))
	return nil
}
