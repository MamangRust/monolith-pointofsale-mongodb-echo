package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	userrole_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/user_role_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRoleCommandRepository struct {
	coll *mongo.Collection
}

func NewUserRoleCommandRepository(coll *mongo.Collection) UserRoleCommandRepository {
	return &userRoleCommandRepository{coll: coll}
}

func (r *userRoleCommandRepository) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*salemongo.UserRole, error) {
	userRoleID, err := salemongo.NextID(ctx, r.coll.Database(), "user_roles")
	if err != nil {
		return nil, userrole_errors.ErrAssignRoleToUser.WithInternal(err)
	}

	now := time.Now()
	userRole := &salemongo.UserRole{
		UserRoleID: userRoleID,
		UserID:     int32(req.UserId),
		RoleID:     int32(req.RoleId),
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}

	_, err = r.coll.InsertOne(ctx, userRole)
	if err != nil {
		return nil, userrole_errors.ErrAssignRoleToUser.WithInternal(err)
	}
	return userRole, nil
}

func (r *userRoleCommandRepository) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{
		"user_id":    int32(req.UserId),
		"role_id":    int32(req.RoleId),
		"deleted_at": bson.M{"$exists": false},
	})
	if err != nil {
		return userrole_errors.ErrRemoveRole.WithInternal(err)
	}
	return nil
}
