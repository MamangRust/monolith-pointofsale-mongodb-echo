package repository

import (
	"context"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type cashierQueryRepository struct {
	coll *mongo.Collection
}

func NewCashierQueryRepository(coll *mongo.Collection) CashierQueryRepository {
	return &cashierQueryRepository{coll: coll}
}

func (r *cashierQueryRepository) FindAllCashiers(ctx context.Context, req *requests.FindAllCashiers) ([]*CashierResult, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	// Count total
	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, cashier_errors.ErrFindAllCashiers
	}

	// Find with pagination
	opts := options.Find().
		SetSort(bson.D{{Key: "cashier_id", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(req.PageSize))

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, nil, cashier_errors.ErrFindAllCashiers
	}
	defer cursor.Close(ctx)

	var cashiers []salemongo.Cashier
	if err := cursor.All(ctx, &cashiers); err != nil {
		return nil, nil, cashier_errors.ErrFindAllCashiers
	}

	var results []*CashierResult
	for _, c := range cashiers {
		cr := &CashierResult{
			CashierID:  c.CashierID,
			MerchantID: c.MerchantID,
			UserID:     c.UserID,
			Name:       c.Name,
			TotalCount: totalCount,
		}
		if c.CreatedAt != nil {
			cr.CreatedAt = c.CreatedAt
		}
		if c.UpdatedAt != nil {
			cr.UpdatedAt = c.UpdatedAt
		}
		results = append(results, cr)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *cashierQueryRepository) FindById(ctx context.Context, cashier_id int) (*salemongo.Cashier, error) {
	var cashier salemongo.Cashier
	err := r.coll.FindOne(ctx, bson.M{"cashier_id": cashier_id}).Decode(&cashier)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, cashier_errors.ErrFindCashierById
		}
		return nil, cashier_errors.ErrFindCashierById
	}
	return &cashier, nil
}

func (r *cashierQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllCashiers) ([]*CashierResultDeleteAt, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	// Count total
	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, cashier_errors.ErrFindActiveCashiers
	}

	// Find with pagination
	opts := options.Find().
		SetSort(bson.D{{Key: "cashier_id", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(req.PageSize))

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, nil, cashier_errors.ErrFindActiveCashiers
	}
	defer cursor.Close(ctx)

	var cashiers []salemongo.Cashier
	if err := cursor.All(ctx, &cashiers); err != nil {
		return nil, nil, cashier_errors.ErrFindActiveCashiers
	}

	var results []*CashierResultDeleteAt
	for _, c := range cashiers {
		cr := &CashierResultDeleteAt{
			CashierID:  c.CashierID,
			MerchantID: c.MerchantID,
			UserID:     c.UserID,
			Name:       c.Name,
			TotalCount: totalCount,
		}
		if c.CreatedAt != nil {
			cr.CreatedAt = c.CreatedAt
		}
		if c.UpdatedAt != nil {
			cr.UpdatedAt = c.UpdatedAt
		}
		if c.DeletedAt != nil {
			cr.DeletedAt = c.DeletedAt
		}
		results = append(results, cr)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *cashierQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllCashiers) ([]*CashierResultDeleteAt, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	// Count total
	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, cashier_errors.ErrFindTrashedCashiers
	}

	// Find with pagination
	opts := options.Find().
		SetSort(bson.D{{Key: "cashier_id", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(req.PageSize))

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, nil, cashier_errors.ErrFindTrashedCashiers
	}
	defer cursor.Close(ctx)

	var cashiers []salemongo.Cashier
	if err := cursor.All(ctx, &cashiers); err != nil {
		return nil, nil, cashier_errors.ErrFindTrashedCashiers
	}

	var results []*CashierResultDeleteAt
	for _, c := range cashiers {
		cr := &CashierResultDeleteAt{
			CashierID:  c.CashierID,
			MerchantID: c.MerchantID,
			UserID:     c.UserID,
			Name:       c.Name,
			TotalCount: totalCount,
		}
		if c.CreatedAt != nil {
			cr.CreatedAt = c.CreatedAt
		}
		if c.UpdatedAt != nil {
			cr.UpdatedAt = c.UpdatedAt
		}
		if c.DeletedAt != nil {
			cr.DeletedAt = c.DeletedAt
		}
		results = append(results, cr)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *cashierQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant) ([]*CashierResult, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"merchant_id": req.MerchantID}
	if req.Search != "" {
		filter["name"] = bson.M{"$regex": req.Search, "$options": "i"}
	}

	// Count total
	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, cashier_errors.ErrFindCashiersByMerchant
	}

	// Find with pagination
	opts := options.Find().
		SetSort(bson.D{{Key: "cashier_id", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(req.PageSize))

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, nil, cashier_errors.ErrFindCashiersByMerchant
	}
	defer cursor.Close(ctx)

	var cashiers []salemongo.Cashier
	if err := cursor.All(ctx, &cashiers); err != nil {
		return nil, nil, cashier_errors.ErrFindCashiersByMerchant
	}

	var results []*CashierResult
	for _, c := range cashiers {
		cr := &CashierResult{
			CashierID:  c.CashierID,
			MerchantID: c.MerchantID,
			UserID:     c.UserID,
			Name:       c.Name,
			TotalCount: totalCount,
		}
		if c.CreatedAt != nil {
			cr.CreatedAt = c.CreatedAt
		}
		if c.UpdatedAt != nil {
			cr.UpdatedAt = c.UpdatedAt
		}
		results = append(results, cr)
	}
	total := int(totalCount)
	return results, &total, nil
}
