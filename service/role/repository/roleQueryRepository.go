package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type roleQueryRepository struct {
	coll *mongo.Collection
}

func NewRoleQueryRepository(coll *mongo.Collection) *roleQueryRepository {
	return &roleQueryRepository{coll: coll}
}

func (r *roleQueryRepository) FindAllRoles(ctx context.Context, req *requests.FindAllRoles) ([]*RoleResult, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["role_name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "role_id", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(req.PageSize))

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var roles []salemongo.Role
	if err := cursor.All(ctx, &roles); err != nil {
		return nil, err
	}

	var results []*RoleResult
	for _, role := range roles {
		rr := &RoleResult{
			RoleID:     role.RoleID,
			RoleName:   role.RoleName,
			TotalCount: total,
		}
		if role.CreatedAt != nil {
			rr.CreatedAt = role.CreatedAt
		}
		if role.UpdatedAt != nil {
			rr.UpdatedAt = role.UpdatedAt
		}
		if role.DeletedAt != nil {
			rr.DeletedAt = role.DeletedAt
		}
		results = append(results, rr)
	}

	return results, nil
}

func (r *roleQueryRepository) FindById(ctx context.Context, id int) (*salemongo.Role, error) {
	var role salemongo.Role
	err := r.coll.FindOne(ctx, bson.M{"role_id": id, "deleted_at": nil}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleQueryRepository) FindByName(ctx context.Context, name string) (*salemongo.Role, error) {
	var role salemongo.Role
	err := r.coll.FindOne(ctx, bson.M{"role_name": name, "deleted_at": nil}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleQueryRepository) FindByUserId(ctx context.Context, userID int) ([]*salemongo.Role, error) {
	// Find role IDs from user_roles collection first
	userRolesColl := r.coll.Database().Collection("user_roles")
	cursor, err := userRolesColl.Find(ctx, bson.M{"user_id": userID, "deleted_at": nil})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var userRoles []struct {
		RoleID int32 `bson:"role_id"`
	}
	if err := cursor.All(ctx, &userRoles); err != nil {
		return nil, err
	}

	if len(userRoles) == 0 {
		return nil, nil
	}

	// Extract role IDs
	roleIDs := make([]int32, 0, len(userRoles))
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	// Find roles by IDs
	rolesCursor, err := r.coll.Find(ctx, bson.M{
		"role_id":    bson.M{"$in": roleIDs},
		"deleted_at": nil,
	})
	if err != nil {
		return nil, err
	}
	defer rolesCursor.Close(ctx)

	var roles []*salemongo.Role
	if err := rolesCursor.All(ctx, &roles); err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *roleQueryRepository) FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*RoleResult, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["role_name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "role_id", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(req.PageSize))

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var roles []salemongo.Role
	if err := cursor.All(ctx, &roles); err != nil {
		return nil, err
	}

	var results []*RoleResult
	for _, role := range roles {
		rr := &RoleResult{
			RoleID:     role.RoleID,
			RoleName:   role.RoleName,
			TotalCount: total,
		}
		if role.CreatedAt != nil {
			rr.CreatedAt = role.CreatedAt
		}
		if role.UpdatedAt != nil {
			rr.UpdatedAt = role.UpdatedAt
		}
		if role.DeletedAt != nil {
			rr.DeletedAt = role.DeletedAt
		}
		results = append(results, rr)
	}

	return results, nil
}

func (r *roleQueryRepository) FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*RoleResult, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if req.Search != "" {
		filter["role_name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "role_id", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(req.PageSize))

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var roles []salemongo.Role
	if err := cursor.All(ctx, &roles); err != nil {
		return nil, err
	}

	var results []*RoleResult
	now := time.Now()
	for _, role := range roles {
		rr := &RoleResult{
			RoleID:     role.RoleID,
			RoleName:   role.RoleName,
			TotalCount: total,
		}
		if role.CreatedAt != nil {
			rr.CreatedAt = role.CreatedAt
		}
		if role.UpdatedAt != nil {
			rr.UpdatedAt = role.UpdatedAt
		}
		_ = now
		if role.DeletedAt != nil {
			rr.DeletedAt = role.DeletedAt
		}
		results = append(results, rr)
	}

	return results, nil
}
