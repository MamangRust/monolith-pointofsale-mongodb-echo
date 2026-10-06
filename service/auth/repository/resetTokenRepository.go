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
)

type resetTokenRepository struct {
	coll *mongo.Collection
}

func NewResetTokenRepository(coll *mongo.Collection) *resetTokenRepository {
	return &resetTokenRepository{coll: coll}
}

func (r *resetTokenRepository) FindByToken(ctx context.Context, code string) (*salemongo.ResetToken, error) {
	var rt salemongo.ResetToken
	err := r.coll.FindOne(ctx, bson.M{"token": code}).Decode(&rt)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, sharedErrors.ErrNotFound
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &rt, nil
}

func (r *resetTokenRepository) CreateResetToken(ctx context.Context, req *requests.CreateResetTokenRequest) (*salemongo.ResetToken, error) {
	resetTokenID, err := salemongo.NextID(ctx, r.coll.Database(), "reset_tokens")
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	expiryDate, err := time.Parse("2006-01-02 15:04:05", req.ExpiredAt)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	rt := &salemongo.ResetToken{
		ResetTokenID: resetTokenID,
		UserID:       int64(req.UserID),
		Token:        req.ResetToken,
		ExpiryDate:   expiryDate,
	}

	result, err := r.coll.InsertOne(ctx, rt)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	rt.ID = result.InsertedID.(primitive.ObjectID)
	return rt, nil
}

func (r *resetTokenRepository) DeleteResetToken(ctx context.Context, user_id int) error {
	_, err := r.coll.DeleteMany(ctx, bson.M{"user_id": user_id})
	if err != nil {
		return sharedErrors.ErrInternal.WithInternal(err)
	}
	return nil
}
