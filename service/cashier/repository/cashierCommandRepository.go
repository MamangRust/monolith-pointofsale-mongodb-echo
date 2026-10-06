package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type cashierCommandRepository struct {
	coll *mongo.Collection
}

func NewCashierCommandRepository(coll *mongo.Collection) CashierCommandRepository {
	return &cashierCommandRepository{coll: coll}
}

func (r *cashierCommandRepository) CreateCashier(ctx context.Context, request *requests.CreateCashierRequest) (*salemongo.Cashier, error) {
	cashierID, err := salemongo.NextID(ctx, r.coll.Database(), "cashiers")
	if err != nil {
		return nil, cashier_errors.ErrCreateCashier
	}

	now := time.Now()
	cashier := &salemongo.Cashier{
		CashierID:  cashierID,
		MerchantID: int32(request.MerchantID),
		UserID:     int32(request.UserID),
		Name:       request.Name,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}

	result, err := r.coll.InsertOne(ctx, cashier)
	if err != nil {
		return nil, cashier_errors.ErrCreateCashier
	}

	cashier.ID = result.InsertedID.(primitive.ObjectID)
	return cashier, nil
}

func (r *cashierCommandRepository) UpdateCashier(ctx context.Context, request *requests.UpdateCashierRequest) (*salemongo.Cashier, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"cashier_id": *request.CashierID},
		bson.M{
			"$set": bson.M{
				"name":       request.Name,
				"updated_at": now,
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var cashier salemongo.Cashier
	err := result.Decode(&cashier)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, cashier_errors.ErrUpdateCashier
		}
		return nil, cashier_errors.ErrUpdateCashier
	}
	return &cashier, nil
}

func (r *cashierCommandRepository) TrashedCashier(ctx context.Context, cashier_id int) (*salemongo.Cashier, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"cashier_id": cashier_id},
		bson.M{
			"$set": bson.M{
				"deleted_at": now,
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var cashier salemongo.Cashier
	err := result.Decode(&cashier)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, cashier_errors.ErrTrashedCashier
		}
		return nil, cashier_errors.ErrTrashedCashier
	}
	return &cashier, nil
}

func (r *cashierCommandRepository) RestoreCashier(ctx context.Context, cashier_id int) (*salemongo.Cashier, error) {
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"cashier_id": cashier_id},
		bson.M{
			"$unset": bson.M{
				"deleted_at": "",
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	var cashier salemongo.Cashier
	err := result.Decode(&cashier)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, cashier_errors.ErrRestoreCashier
		}
		return nil, cashier_errors.ErrRestoreCashier
	}
	return &cashier, nil
}

func (r *cashierCommandRepository) DeleteCashierPermanent(ctx context.Context, cashier_id int) (bool, error) {
	result, err := r.coll.DeleteOne(ctx, bson.M{"cashier_id": cashier_id})
	if err != nil {
		return false, cashier_errors.ErrDeleteCashierPermanent
	}
	return result.DeletedCount > 0, nil
}

func (r *cashierCommandRepository) RestoreAllCashier(ctx context.Context) (bool, error) {
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
		return false, cashier_errors.ErrRestoreAllCashiers
	}
	return result.ModifiedCount > 0, nil
}

func (r *cashierCommandRepository) DeleteAllCashierPermanent(ctx context.Context) (bool, error) {
	result, err := r.coll.DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, cashier_errors.ErrDeleteAllCashiersPermanent
	}
	return result.DeletedCount > 0, nil
}
