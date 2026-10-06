package repository

import "go.mongodb.org/mongo-driver/mongo"

type Repositories struct {
	CategoryCommand         CategoryCommandRepository
	CategoryQuery           CategoryQueryRepository
	CategoryStats           CategoryStatsRepository
	CategoryStatsById       CategoryStatsByIdRepository
	CategoryStatsByMerchant CategoryStatsByMerchantRepository
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		CategoryCommand:         NewCategoryCommandRepository(db.Collection("categories")),
		CategoryQuery:           NewCategoryQueryRepository(db.Collection("categories")),
		CategoryStats:           NewCategoryStatsRepository(db),
		CategoryStatsById:       NewCategoryStatsByIdRepository(db),
		CategoryStatsByMerchant: NewCategoryStatsByMerchantRepository(db),
	}
}
