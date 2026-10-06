package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type merchantCommandRepository struct {
	coll *mongodriver.Collection
}

func NewMerchantCommandRepository(coll *mongodriver.Collection) MerchantCommandRepository {
	return &merchantCommandRepository{coll: coll}
}

func timePtr(t time.Time) *time.Time { return &t }

func (r *merchantCommandRepository) CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*salemongo.Merchant, error) {
	merchantID, err := salemongo.NextID(ctx, r.coll.Database(), "merchants")
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	now := time.Now()
	merchant := &salemongo.Merchant{
		MerchantID:   merchantID,
		UserID:       int32(request.UserID),
		Name:         request.Name,
		Description:  &request.Description,
		Address:      &request.Address,
		ContactEmail: &request.ContactEmail,
		ContactPhone: &request.ContactPhone,
		Status:       "inactive",
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}
	result, err := r.coll.InsertOne(ctx, merchant)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	merchant.ID = result.InsertedID.(primitive.ObjectID)
	return merchant, nil
}

func (r *merchantCommandRepository) UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*salemongo.Merchant, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"merchant_id": *request.MerchantID, "deleted_at": nil},
		bson.M{"$set": bson.M{
			"name":          request.Name,
			"description":   request.Description,
			"address":       request.Address,
			"contact_email": request.ContactEmail,
			"contact_phone": request.ContactPhone,
			"status":        request.Status,
			"updated_at":    now,
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var merchant salemongo.Merchant
	if err := result.Decode(&merchant); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &merchant, nil
}

func (r *merchantCommandRepository) UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*salemongo.Merchant, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"merchant_id": *request.MerchantID, "deleted_at": nil},
		bson.M{"$set": bson.M{
			"status":     request.Status,
			"updated_at": now,
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var merchant salemongo.Merchant
	if err := result.Decode(&merchant); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &merchant, nil
}

func (r *merchantCommandRepository) TrashedMerchant(ctx context.Context, merchantID int) (*salemongo.Merchant, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"merchant_id": merchantID, "deleted_at": nil},
		bson.M{"$set": bson.M{"deleted_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var merchant salemongo.Merchant
	if err := result.Decode(&merchant); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &merchant, nil
}

func (r *merchantCommandRepository) RestoreMerchant(ctx context.Context, merchantID int) (*salemongo.Merchant, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"merchant_id": merchantID, "deleted_at": bson.M{"$ne": nil}},
		bson.M{"$unset": bson.M{"deleted_at": ""}, "$set": bson.M{"updated_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var merchant salemongo.Merchant
	if err := result.Decode(&merchant); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &merchant, nil
}

func (r *merchantCommandRepository) DeleteMerchantPermanent(ctx context.Context, merchantID int) (bool, error) {
	result, err := r.coll.DeleteOne(ctx, bson.M{"merchant_id": merchantID, "deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}
	return result.DeletedCount > 0, nil
}

func (r *merchantCommandRepository) RestoreAllMerchant(ctx context.Context) (bool, error) {
	_, err := r.coll.UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$unset": bson.M{"deleted_at": ""}})
	if err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}
	return true, nil
}

func (r *merchantCommandRepository) DeleteAllMerchantPermanent(ctx context.Context) (bool, error) {
	_, err := r.coll.DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}
	return true, nil
}
