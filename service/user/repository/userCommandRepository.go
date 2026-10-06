package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userCommandRepository struct {
	coll *mongodriver.Collection
}

func NewUserCommandRepository(coll *mongodriver.Collection) *userCommandRepository {
	return &userCommandRepository{coll: coll}
}

func timePtr(t time.Time) *time.Time { return &t }

func (r *userCommandRepository) CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*salemongo.User, error) {
	userID, err := salemongo.NextID(ctx, r.coll.Database(), "users")
	if err != nil {
		return nil, err
	}

	user := &salemongo.User{
		UserID:    userID,
		Firstname: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
		CreatedAt: timePtr(time.Now()),
		UpdatedAt: timePtr(time.Now()),
	}
	result, err := r.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

// CreateUserRecord stores a registration sent by auth-service. The password
// arrives already hashed (auth owns hashing), so unlike CreateUser this path
// must not hash again — doing so would make the stored hash unusable at login.
func (r *userCommandRepository) CreateUserRecord(ctx context.Context, request *requests.RegisterRequest) (*salemongo.User, error) {
	userID, err := salemongo.NextID(ctx, r.coll.Database(), "users")
	if err != nil {
		return nil, err
	}

	now := time.Now()
	isVerified := request.IsVerified
	user := &salemongo.User{
		UserID:           userID,
		Firstname:        request.FirstName,
		Lastname:         request.LastName,
		Email:            request.Email,
		Password:         request.Password,
		VerificationCode: request.VerifiedCode,
		IsVerified:       &isVerified,
		CreatedAt:        &now,
		UpdatedAt:        &now,
	}
	result, err := r.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (r *userCommandRepository) UpdateUserIsVerified(ctx context.Context, userID int, isVerified bool) (*salemongo.User, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": userID},
		bson.M{"$set": bson.M{
			"is_verified": isVerified,
			"updated_at":  now,
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var user salemongo.User
	if err := result.Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserPassword stores a password hash supplied by auth-service; it does
// not hash the value itself.
func (r *userCommandRepository) UpdateUserPassword(ctx context.Context, userID int, password string) (*salemongo.User, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": userID},
		bson.M{"$set": bson.M{
			"password":   password,
			"updated_at": now,
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var user salemongo.User
	if err := result.Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userCommandRepository) UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*salemongo.User, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": *request.UserID},
		bson.M{"$set": bson.M{
			"firstname":  request.FirstName,
			"lastname":   request.LastName,
			"email":      request.Email,
			"password":   request.Password,
			"updated_at": now,
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var user salemongo.User
	if err := result.Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userCommandRepository) TrashedUser(ctx context.Context, userID int) (*salemongo.User, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": userID},
		bson.M{"$set": bson.M{"deleted_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var user salemongo.User
	if err := result.Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userCommandRepository) RestoreUser(ctx context.Context, userID int) (*salemongo.User, error) {
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": userID, "deleted_at": bson.M{"$ne": nil}},
		bson.M{"$unset": bson.M{"deleted_at": ""}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var user salemongo.User
	if err := result.Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userCommandRepository) DeleteUserPermanent(ctx context.Context, userID int) (bool, error) {
	result, err := r.coll.DeleteOne(ctx, bson.M{"user_id": userID})
	if err != nil {
		return false, err
	}
	return result.DeletedCount > 0, nil
}

func (r *userCommandRepository) RestoreAllUser(ctx context.Context) (bool, error) {
	_, err := r.coll.UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$unset": bson.M{"deleted_at": ""}})
	return err == nil, err
}

func (r *userCommandRepository) DeleteAllUserPermanent(ctx context.Context) (bool, error) {
	_, err := r.coll.DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	return err == nil, err
}
