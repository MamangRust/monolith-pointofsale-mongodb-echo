package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type roleCommandRepository struct {
	coll *mongo.Collection
}

func NewRoleCommandRepository(coll *mongo.Collection) *roleCommandRepository {
	return &roleCommandRepository{coll: coll}
}

func timePtr(t time.Time) *time.Time { return &t }

func (r *roleCommandRepository) CreateRole(ctx context.Context, req *requests.CreateRoleRequest) (*salemongo.Role, error) {
	roleID, err := salemongo.NextID(ctx, r.coll.Database(), "roles")
	if err != nil {
		return nil, err
	}

	now := time.Now()
	role := &salemongo.Role{
		RoleID:    roleID,
		RoleName:  req.Name,
		CreatedAt: timePtr(now),
		UpdatedAt: timePtr(now),
	}

	result, err := r.coll.InsertOne(ctx, role)
	if err != nil {
		return nil, err
	}

	role.ID = result.InsertedID.(primitive.ObjectID)
	return role, nil
}

func (r *roleCommandRepository) UpdateRole(ctx context.Context, req *requests.UpdateRoleRequest) (*salemongo.Role, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"role_id": *req.ID},
		bson.M{
			"$set": bson.M{
				"role_name":  req.Name,
				"updated_at": now,
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var role salemongo.Role
	err := result.Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleCommandRepository) TrashedRole(ctx context.Context, id int) (*salemongo.Role, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"role_id": id},
		bson.M{
			"$set": bson.M{
				"deleted_at": now,
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var role salemongo.Role
	err := result.Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleCommandRepository) RestoreRole(ctx context.Context, id int) (*salemongo.Role, error) {
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"role_id": id},
		bson.M{
			"$unset": bson.M{
				"deleted_at": "",
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var role salemongo.Role
	err := result.Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleCommandRepository) DeleteRolePermanent(ctx context.Context, id int) (bool, error) {
	result, err := r.coll.DeleteOne(ctx, bson.M{"role_id": id})
	if err != nil {
		return false, err
	}
	return result.DeletedCount > 0, nil
}

func (r *roleCommandRepository) RestoreAllRole(ctx context.Context) (bool, error) {
	result, err := r.coll.UpdateMany(
		ctx,
		bson.M{"deleted_at": bson.M{"$ne": nil}},
		bson.M{
			"$unset": bson.M{
				"deleted_at": "",
			},
		},
	)
	if err != nil {
		return false, err
	}
	return result.ModifiedCount > 0, nil
}

func (r *roleCommandRepository) DeleteAllRolePermanent(ctx context.Context) (bool, error) {
	result, err := r.coll.DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, err
	}
	return result.DeletedCount > 0, nil
}
