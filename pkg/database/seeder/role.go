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
)

type roleSeeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewRoleSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *roleSeeder {
	return &roleSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *roleSeeder) Seed() error {
	randomRoles := []string{
		"ROLE_ADMIN",
		"Admin Access 1",
		"Super Admin",
		"Admin",
		"Store Manager",
		"Cashier",
		"Inventory Staff",
		"Support",
		"Auditor",
		"Viewer",
	}

	totalRoles := len(randomRoles)
	coll := r.db.Collection("roles")
	now := time.Now().UTC()

	for i, roleName := range randomRoles {
		roleID, err := salemongo.NextID(r.ctx, r.db, "roles")
		if err != nil {
			r.logger.Error("failed to generate role id", zap.Int("role", i+1), zap.String("roleName", roleName), zap.Error(err))
			return fmt.Errorf("failed to generate role id for %s: %w", roleName, err)
		}

		role := &salemongo.Role{
			RoleID:    roleID,
			RoleName:  roleName,
			CreatedAt: &now,
			UpdatedAt: &now,
		}

		_, err = coll.UpdateOne(r.ctx,
			bson.M{"role_name": roleName},
			bson.M{"$setOnInsert": role},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			r.logger.Error("failed to seed role", zap.Int("role", i+1), zap.String("roleName", roleName), zap.Error(err))
			return fmt.Errorf("failed to seed role %d (%s): %w", i+1, roleName, err)
		}
	}

	r.logger.Info("role seeded successfully", zap.Int("totalRoles", totalRoles))
	return nil
}
