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

type merchantDocumentQueryRepository struct {
	coll *mongodriver.Collection
}

func NewMerchantDocumentQueryRepository(coll *mongodriver.Collection) MerchantDocumentQueryRepository {
	return &merchantDocumentQueryRepository{coll: coll}
}

func (r *merchantDocumentQueryRepository) FindAllDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResult, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"document_type": bson.M{"$regex": req.Search, "$options": "i"}},
			{"document_url": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "document_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var docs []salemongo.MerchantDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var results []*MerchantDocumentResult
	for _, d := range docs {
		dr := &MerchantDocumentResult{
			DocumentID:   d.DocumentID,
			MerchantID:   d.MerchantID,
			DocumentType: d.DocumentType,
			DocumentUrl:  d.DocumentUrl,
			Status:       d.Status,
			Note:         d.Note,
			TotalCount:   totalCount,
		}
		if d.CreatedAt != nil {
			dr.CreatedAt = d.CreatedAt
		}
		if d.UpdatedAt != nil {
			dr.UpdatedAt = d.UpdatedAt
		}
		results = append(results, dr)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *merchantDocumentQueryRepository) FindById(ctx context.Context, id int) (*salemongo.MerchantDocument, error) {
	var doc salemongo.MerchantDocument
	err := r.coll.FindOne(ctx, bson.M{"document_id": id, "deleted_at": nil}).Decode(&doc)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResultDeleteAt, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": nil}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"document_type": bson.M{"$regex": req.Search, "$options": "i"}},
			{"document_url": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "document_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var docs []salemongo.MerchantDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var results []*MerchantDocumentResultDeleteAt
	for _, d := range docs {
		dr := &MerchantDocumentResultDeleteAt{
			DocumentID:   d.DocumentID,
			MerchantID:   d.MerchantID,
			DocumentType: d.DocumentType,
			DocumentUrl:  d.DocumentUrl,
			Status:       d.Status,
			Note:         d.Note,
			TotalCount:   totalCount,
		}
		if d.CreatedAt != nil {
			dr.CreatedAt = d.CreatedAt
		}
		if d.UpdatedAt != nil {
			dr.UpdatedAt = d.UpdatedAt
		}
		results = append(results, dr)
	}
	total := int(totalCount)
	return results, &total, nil
}

func (r *merchantDocumentQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResultDeleteAt, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	filter := bson.M{"deleted_at": bson.M{"$ne": nil}}
	if req.Search != "" {
		filter["$or"] = []bson.M{
			{"document_type": bson.M{"$regex": req.Search, "$options": "i"}},
			{"document_url": bson.M{"$regex": req.Search, "$options": "i"}},
		}
	}

	totalCount, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	cursor, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "document_id", Value: -1}}).SetSkip(int64(offset)).SetLimit(int64(req.PageSize)))
	if err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	defer cursor.Close(ctx)

	var docs []salemongo.MerchantDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var results []*MerchantDocumentResultDeleteAt
	for _, d := range docs {
		dr := &MerchantDocumentResultDeleteAt{
			DocumentID:   d.DocumentID,
			MerchantID:   d.MerchantID,
			DocumentType: d.DocumentType,
			DocumentUrl:  d.DocumentUrl,
			Status:       d.Status,
			Note:         d.Note,
			TotalCount:   totalCount,
		}
		if d.CreatedAt != nil {
			dr.CreatedAt = d.CreatedAt
		}
		if d.UpdatedAt != nil {
			dr.UpdatedAt = d.UpdatedAt
		}
		if d.DeletedAt != nil {
			dr.DeletedAt = d.DeletedAt
		}
		results = append(results, dr)
	}
	total := int(totalCount)
	return results, &total, nil
}
