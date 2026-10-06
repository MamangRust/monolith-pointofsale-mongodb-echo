package repository

import "go.mongodb.org/mongo-driver/mongo"

type Repositories struct {
	OrderItemQuery   OrderItemQueryRepository
	OrderItemCommand OrderItemCommandRepository
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		OrderItemQuery:   NewOrderItemQueryRepository(db.Collection("order_items")),
		OrderItemCommand: NewOrderItemCommandRepository(db.Collection("order_items")),
	}
}
