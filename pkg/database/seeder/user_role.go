package seeder

import (
	"context"
	"fmt"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

type userRoleSeeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewUserRoleSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *userRoleSeeder {
	return &userRoleSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *userRoleSeeder) Seed() error {
	users, err := findAll[*salemongo.User](r.ctx, r.db.Collection("users"), bson.M{"deleted_at": nil}, 20)
	if err != nil {
		r.logger.Error("failed to fetch users", zap.Error(err))
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	roles, err := findAll[*salemongo.Role](r.ctx, r.db.Collection("roles"), bson.M{}, 4)
	if err != nil {
		r.logger.Error("failed to fetch roles", zap.Error(err))
		return fmt.Errorf("failed to fetch roles: %w", err)
	}

	if len(users) == 0 || len(roles) == 0 {
		r.logger.Debug("no users or roles available for seeding")
		return nil
	}

	rand.Seed(uint64(time.Now().UnixNano()))

	coll := r.db.Collection("user_roles")
	now := time.Now().UTC()

	for _, user := range users {
		role := roles[rand.Intn(len(roles))]

		userRoleID, err := salemongo.NextID(r.ctx, r.db, "user_roles")
		if err != nil {
			return fmt.Errorf("failed to generate user_role id: %w", err)
		}

		userRole := &salemongo.UserRole{
			UserRoleID: userRoleID,
			UserID:     user.UserID,
			RoleID:     role.RoleID,
			CreatedAt:  &now,
			UpdatedAt:  &now,
		}

		if _, err := coll.InsertOne(r.ctx, userRole); err != nil {
			r.logger.Error("failed to assign role to user", zap.String("user", user.Email), zap.String("role", role.RoleName), zap.Error(err))
			return fmt.Errorf("failed to assign role %s to user %s: %w", role.RoleName, user.Email, err)
		}
	}

	r.logger.Info("user roles assigned successfully", zap.Int("users", len(users)), zap.Int("roles", len(roles)))
	return nil
}

// findAll returns up to limit documents matching filter, newest first.
func findAll[T any](ctx context.Context, coll *mongo.Collection, filter bson.M, limit int64) ([]T, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(limit)

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var out []T
	if err := cursor.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}
