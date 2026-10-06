package repository

import (
	"context"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/user_errors"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userQueryRepository struct {
	coll *mongodriver.Collection
}

func NewUserQueryRepository(coll *mongodriver.Collection) *userQueryRepository {
	return &userQueryRepository{coll: coll}
}

func (r *userQueryRepository) FindAllUsers(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"firstname": bson.M{"$regex": req.Search, "$options": "i"}},
			{"lastname": bson.M{"$regex": req.Search, "$options": "i"}},
			{"email": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "user_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []salemongo.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	var results []*UserResult
	for _, u := range users {
		ur := &UserResult{
			UserID:     u.UserID,
			Firstname:  u.Firstname,
			Lastname:   u.Lastname,
			Email:      u.Email,
			TotalCount: totalCount,
		}
		if u.CreatedAt != nil {
			ur.CreatedAt = u.CreatedAt
		}
		if u.UpdatedAt != nil {
			ur.UpdatedAt = u.UpdatedAt
		}
		results = append(results, ur)
	}
	return results, nil
}

func (r *userQueryRepository) FindById(ctx context.Context, userID int) (*salemongo.User, error) {
	var user salemongo.User
	err := r.coll.FindOne(ctx, bson.M{"user_id": userID, "deleted_at": nil}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userQueryRepository) FindByEmail(ctx context.Context, email string) (*salemongo.User, error) {
	var user salemongo.User
	err := r.coll.FindOne(ctx, bson.M{"email": email, "deleted_at": nil}).Decode(&user)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmailAndVerify is the login lookup: verified, active users only.
func (r *userQueryRepository) FindByEmailAndVerify(ctx context.Context, email string) (*salemongo.User, error) {
	var user salemongo.User
	err := r.coll.FindOne(ctx, bson.M{"email": email, "is_verified": true, "deleted_at": nil}).Decode(&user)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userQueryRepository) FindByVerificationCode(ctx context.Context, verificationCode string) (*salemongo.User, error) {
	var user salemongo.User
	err := r.coll.FindOne(ctx, bson.M{"verification_code": verificationCode, "deleted_at": nil}).Decode(&user)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return nil, user_errors.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"firstname": bson.M{"$regex": req.Search, "$options": "i"}},
			{"lastname": bson.M{"$regex": req.Search, "$options": "i"}},
			{"email": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "user_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []salemongo.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	var results []*UserResult
	for _, u := range users {
		ur := &UserResult{
			UserID:     u.UserID,
			Firstname:  u.Firstname,
			Lastname:   u.Lastname,
			Email:      u.Email,
			TotalCount: totalCount,
		}
		if u.CreatedAt != nil {
			ur.CreatedAt = u.CreatedAt
		}
		if u.UpdatedAt != nil {
			ur.UpdatedAt = u.UpdatedAt
		}
		results = append(results, ur)
	}
	return results, nil
}

func (r *userQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*UserResult, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"firstname": bson.M{"$regex": req.Search, "$options": "i"}},
			{"lastname": bson.M{"$regex": req.Search, "$options": "i"}},
			{"email": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "user_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []salemongo.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	var results []*UserResult
	for _, u := range users {
		ur := &UserResult{
			UserID:     u.UserID,
			Firstname:  u.Firstname,
			Lastname:   u.Lastname,
			Email:      u.Email,
			TotalCount: totalCount,
		}
		if u.CreatedAt != nil {
			ur.CreatedAt = u.CreatedAt
		}
		if u.UpdatedAt != nil {
			ur.UpdatedAt = u.UpdatedAt
		}
		if u.DeletedAt != nil {
			ur.DeletedAt = u.DeletedAt
		}
		results = append(results, ur)
	}
	return results, nil
}
