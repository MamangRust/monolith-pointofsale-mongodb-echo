package repository

import "go.mongodb.org/mongo-driver/mongo"

type Repositories struct {
	RoleCommand     RoleCommandRepository
	RoleQuery       RoleQueryRepository
	UserRoleCommand UserRoleCommandRepository
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		RoleCommand:     NewRoleCommandRepository(db.Collection("roles")),
		RoleQuery:       NewRoleQueryRepository(db.Collection("roles")),
		UserRoleCommand: NewUserRoleCommandRepository(db.Collection("user_roles")),
	}
}
