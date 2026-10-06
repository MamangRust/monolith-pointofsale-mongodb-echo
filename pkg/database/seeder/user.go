package seeder

import (
	"context"
	"fmt"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-pkg/hash"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type userSeeder struct {
	db     *mongo.Database
	hash   hash.HashPassword
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewUserSeeder(db *mongo.Database, hash hash.HashPassword, ctx context.Context, logger logger.LoggerInterface) *userSeeder {
	return &userSeeder{
		db:     db,
		hash:   hash,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *userSeeder) Seed() error {
	coll := r.db.Collection("users")
	now := time.Now().UTC()

	for i := 1; i <= 10; i++ {
		email := fmt.Sprintf("user_%s@example.com", uuid.NewString())
		rawPassword := fmt.Sprintf("password%d", i)

		hashedPassword, err := r.hash.HashPassword(rawPassword)
		if err != nil {
			r.logger.Error("failed to hash password", zap.Int("user", i), zap.Error(err))
			return fmt.Errorf("failed to hash password for user %d: %w", i, err)
		}

		userID, err := salemongo.NextID(r.ctx, r.db, "users")
		if err != nil {
			r.logger.Error("failed to generate user id", zap.Int("user", i), zap.Error(err))
			return fmt.Errorf("failed to generate user id for user %d: %w", i, err)
		}

		verified := true
		user := &salemongo.User{
			UserID:     userID,
			Firstname:  fmt.Sprintf("User%d", i),
			Lastname:   fmt.Sprintf("Last%d", i),
			Email:      email,
			Password:   hashedPassword,
			IsVerified: &verified,
			CreatedAt:  &now,
			UpdatedAt:  &now,
		}

		if _, err := coll.InsertOne(r.ctx, user); err != nil {
			r.logger.Error("failed to seed user", zap.Int("user", i), zap.Error(err))
			return fmt.Errorf("failed to seed user %d: %w", i, err)
		}

		if i > 5 {
			if _, err := coll.UpdateOne(r.ctx,
				bson.M{"user_id": userID},
				bson.M{"$set": bson.M{"deleted_at": now, "updated_at": now}},
			); err != nil {
				r.logger.Error("failed to trash user", zap.Int("user", i), zap.Error(err))
				return fmt.Errorf("failed to trash user %d: %w", i, err)
			}
		}
	}

	r.logger.Info("User seeding completed successfully", zap.Int("totalUsers", 10))
	return nil
}
