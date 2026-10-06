package repository

import (
	"context"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type merchantQueryRepository struct {
	coll *mongodriver.Collection
}

func NewMerchantQueryRepository(coll *mongodriver.Collection) MerchantQueryRepository {
	return &merchantQueryRepository{coll: coll}
}

func fmtStr(t interface{}) string {
	switch v := t.(type) {
	case *salemongo.Merchant:
		if v.CreatedAt != nil {
			return v.CreatedAt.Format("2006-01-02 15:04:05")
		}
	}
	return ""
}

func (r *merchantQueryRepository) FindAllMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*MerchantResult, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "merchant_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var merchants []salemongo.Merchant
	if err := cursor.All(ctx, &merchants); err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var results []*MerchantResult
	for _, m := range merchants {
		mr := &MerchantResult{
			MerchantID:   m.MerchantID,
			UserID:       m.UserID,
			Name:         m.Name,
			Description:  m.Description,
			Address:      m.Address,
			ContactEmail: m.ContactEmail,
			ContactPhone: m.ContactPhone,
			Status:       m.Status,
			TotalCount:   totalCount,
		}
		if m.CreatedAt != nil {
			mr.CreatedAt = m.CreatedAt
		}
		if m.UpdatedAt != nil {
			mr.UpdatedAt = m.UpdatedAt
		}
		results = append(results, mr)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *merchantQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*MerchantResultDeleteAt, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "merchant_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var merchants []salemongo.Merchant
	if err := cursor.All(ctx, &merchants); err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var results []*MerchantResultDeleteAt
	for _, m := range merchants {
		mr := &MerchantResultDeleteAt{
			MerchantID:   m.MerchantID,
			UserID:       m.UserID,
			Name:         m.Name,
			Description:  m.Description,
			Address:      m.Address,
			ContactEmail: m.ContactEmail,
			ContactPhone: m.ContactPhone,
			Status:       m.Status,
			TotalCount:   totalCount,
		}
		if m.CreatedAt != nil {
			mr.CreatedAt = m.CreatedAt
		}
		if m.UpdatedAt != nil {
			mr.UpdatedAt = m.UpdatedAt
		}
		results = append(results, mr)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *merchantQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*MerchantResultDeleteAt, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "merchant_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var merchants []salemongo.Merchant
	if err := cursor.All(ctx, &merchants); err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var results []*MerchantResultDeleteAt
	for _, m := range merchants {
		mr := &MerchantResultDeleteAt{
			MerchantID:   m.MerchantID,
			UserID:       m.UserID,
			Name:         m.Name,
			Description:  m.Description,
			Address:      m.Address,
			ContactEmail: m.ContactEmail,
			ContactPhone: m.ContactPhone,
			Status:       m.Status,
			TotalCount:   totalCount,
		}
		if m.CreatedAt != nil {
			mr.CreatedAt = m.CreatedAt
		}
		if m.UpdatedAt != nil {
			mr.UpdatedAt = m.UpdatedAt
		}
		if m.DeletedAt != nil {
			mr.DeletedAt = m.DeletedAt
		}
		results = append(results, mr)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *merchantQueryRepository) FindById(ctx context.Context, merchantID int) (*salemongo.Merchant, error) {
	var merchant salemongo.Merchant
	err := r.coll.FindOne(ctx, bson.M{"merchant_id": merchantID, "deleted_at": nil}).Decode(&merchant)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &merchant, nil
}
