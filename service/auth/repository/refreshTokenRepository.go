package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type refreshTokenRepository struct {
	coll *mongo.Collection
}

func NewRefreshTokenRepository(coll *mongo.Collection) *refreshTokenRepository {
	return &refreshTokenRepository{coll: coll}
}

func (r *refreshTokenRepository) FindByToken(ctx context.Context, token string) (*salemongo.RefreshToken, error) {
	var rt salemongo.RefreshToken
	err := r.coll.FindOne(ctx, bson.M{"token": token}).Decode(&rt)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, sharedErrors.ErrNotFound
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &rt, nil
}

func (r *refreshTokenRepository) FindByUserId(ctx context.Context, user_id int) (*salemongo.RefreshToken, error) {
	var rt salemongo.RefreshToken
	err := r.coll.FindOne(ctx, bson.M{"user_id": user_id}).Decode(&rt)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, sharedErrors.ErrNotFound
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &rt, nil
}

func (r *refreshTokenRepository) CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*salemongo.RefreshToken, error) {
	refreshTokenID, err := salemongo.NextID(ctx, r.coll.Database(), "refresh_tokens")
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	now := time.Now()
	rt := &salemongo.RefreshToken{
		RefreshTokenID: refreshTokenID,
		UserID:         int32(req.UserId),
		Token:          req.Token,
		Expiration:     expirationTime,
		CreatedAt:      &now,
		UpdatedAt:      &now,
	}

	result, err := r.coll.InsertOne(ctx, rt)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	rt.ID = result.InsertedID.(primitive.ObjectID)
	return rt, nil
}

func (r *refreshTokenRepository) UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*salemongo.RefreshToken, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": req.UserId},
		bson.M{
			"$set": bson.M{
				"token":      req.Token,
				"expiration": expirationTime,
				"updated_at": now,
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var rt salemongo.RefreshToken
	err = result.Decode(&rt)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, sharedErrors.ErrNotFound
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &rt, nil
}

func (r *refreshTokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.coll.DeleteOne(ctx, bson.M{"token": token})
	if err != nil {
		return sharedErrors.ErrInternal.WithInternal(err)
	}
	return nil
}

func (r *refreshTokenRepository) DeleteRefreshTokenByUserId(ctx context.Context, user_id int) error {
	_, err := r.coll.DeleteMany(ctx, bson.M{"user_id": user_id})
	if err != nil {
		return sharedErrors.ErrInternal.WithInternal(err)
	}
	return nil
}
