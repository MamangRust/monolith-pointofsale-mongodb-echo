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

type productSeeder struct {
	db     *mongo.Database
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewProductSeeder(db *mongo.Database, ctx context.Context, logger logger.LoggerInterface) *productSeeder {
	return &productSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *productSeeder) Seed() error {
	merchants, err := findAll[*salemongo.Merchant](r.ctx, r.db.Collection("merchants"), bson.M{"deleted_at": nil}, 20)
	if err != nil {
		r.logger.Error("Failed to get merchants:", zap.Any("error", err))
		return err
	}

	categories, err := findAll[*salemongo.Category](r.ctx, r.db.Collection("categories"), bson.M{"deleted_at": nil}, 20)
	if err != nil {
		r.logger.Error("Failed to get categories:", zap.Any("error", err))
		return err
	}

	if len(merchants) == 0 || len(categories) == 0 {
		r.logger.Error("No merchants or categories found, skipping seeding")
		return nil
	}

	productNames := []string{
		"Smartphone", "Laptop", "Wireless Earbuds", "Gaming Mouse", "Mechanical Keyboard",
		"Smartwatch", "Power Bank", "Bluetooth Speaker", "External Hard Drive", "Monitor",
	}
	brands := []string{"Samsung", "Apple", "Sony", "Logitech", "Razer", "Xiaomi", "HP", "Dell", "Acer", "Asus"}
	images := []string{
		"image1.jpg", "image2.jpg", "image3.jpg", "image4.jpg", "image5.jpg",
		"image6.jpg", "image7.jpg", "image8.jpg", "image9.jpg", "image10.jpg",
	}

	coll := r.db.Collection("products")
	now := time.Now().UTC()

	for i := 0; i < 10; i++ {
		merchant := merchants[rand.Intn(len(merchants))]
		category := categories[rand.Intn(len(categories))]
		name := productNames[rand.Intn(len(productNames))]
		brand := brands[rand.Intn(len(brands))]
		price := int32(rand.Intn(5000000) + 50000)
		countInStock := int32(rand.Intn(100) + 1)
		weight := int32(rand.Intn(5000) + 100)
		slug := fmt.Sprintf("%s-%d", name, rand.Intn(1000))
		image := images[rand.Intn(len(images))]
		barcode := fmt.Sprintf("BC-%d", rand.Intn(9999999))

		productID, err := salemongo.NextID(r.ctx, r.db, "products")
		if err != nil {
			r.logger.Error("Failed to generate product id:", zap.Error(err))
			return err
		}

		product := &salemongo.Product{
			ProductID:    productID,
			MerchantID:   merchant.MerchantID,
			CategoryID:   category.CategoryID,
			Name:         name,
			Description:  ptrString(fmt.Sprintf("Description for %s", name)),
			Price:        price,
			CountInStock: countInStock,
			Brand:        ptrString(brand),
			Weight:       ptrInt32(weight),
			SlugProduct:  ptrString(slug),
			ImageProduct: ptrString(image),
			Barcode:      ptrString(barcode),
			CreatedAt:    &now,
			UpdatedAt:    &now,
		}

		if _, err := coll.InsertOne(r.ctx, product); err != nil {
			r.logger.Error("Failed to create product:", zap.Error(err))
			return err
		}
	}

	r.logger.Info("Product seeding completed successfully.", zap.Int("count", 10))
	return nil
}
