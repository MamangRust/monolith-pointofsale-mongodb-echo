package seeder

import (
	"context"
	"fmt"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type categorySeeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewCategorySeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *categorySeeder {
	return &categorySeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *categorySeeder) Seed() error {
	categoryNames := []string{
		"Electronics", "Clothing", "Groceries", "Toys", "Home & Kitchen",
		"Books", "Beauty & Health", "Sports & Outdoors", "Automotive", "Furniture",
	}

	categoryDescriptions := []string{
		"Best electronics products", "Latest fashion trends", "Fresh groceries",
		"Fun toys for kids", "Essentials for home & kitchen",
		"Books for all ages", "Beauty and health products",
		"Outdoor sports equipment", "Automotive accessories", "Stylish furniture",
	}

	coll := r.db.Collection("categories")
	now := time.Now().UTC()

	for i := 0; i < 10; i++ {
		name := categoryNames[i%len(categoryNames)]
		description := categoryDescriptions[i%len(categoryDescriptions)]
		slugCategory := fmt.Sprintf("%s-%d", name, i+1)

		categoryID, err := salemongo.NextID(r.ctx, r.db, "categories")
		if err != nil {
			r.logger.Error("Failed to generate category id:", zap.Error(err))
			return err
		}

		category := &salemongo.Category{
			CategoryID:   categoryID,
			Name:         name,
			Description:  ptrString(description),
			SlugCategory: ptrString(slugCategory),
			CreatedAt:    &now,
			UpdatedAt:    &now,
		}

		if _, err := coll.InsertOne(r.ctx, category); err != nil {
			r.logger.Error("Failed to create category:", zap.Error(err))
			return err
		}
		r.logger.Debug("Category created:", zap.String("slug", slugCategory))
	}

	r.logger.Info("Category seeding completed successfully.", zap.Int("count", 10))
	return nil
}
