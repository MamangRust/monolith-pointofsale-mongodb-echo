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

type merchantDocumentCommandRepository struct {
	coll *mongodriver.Collection
}

func NewMerchantDocumentCommandRepository(coll *mongodriver.Collection) MerchantDocumentCommandRepository {
	return &merchantDocumentCommandRepository{coll: coll}
}

func (r *merchantDocumentCommandRepository) CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*salemongo.MerchantDocument, error) {
	documentID, err := salemongo.NextID(ctx, r.coll.Database(), "merchant_documents")
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	now := time.Now()
	emptyNote := ""
	doc := &salemongo.MerchantDocument{
		DocumentID:   documentID,
		MerchantID:   int32(request.MerchantID),
		DocumentType: request.DocumentType,
		DocumentUrl:  request.DocumentUrl,
		Status:       "pending",
		Note:         &emptyNote,
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}
	result, err := r.coll.InsertOne(ctx, doc)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	doc.ID = result.InsertedID.(primitive.ObjectID)
	return doc, nil
}

func (r *merchantDocumentCommandRepository) UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*salemongo.MerchantDocument, error) {
	if request.DocumentID == nil {
		return nil, sharedErrors.ErrBadRequest.WithMessage("document id is required")
	}
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"document_id": *request.DocumentID},
		bson.M{"$set": bson.M{
			"document_type": request.DocumentType,
			"document_url":  request.DocumentUrl,
			"status":        request.Status,
			"note":          request.Note,
			"updated_at":    now,
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var doc salemongo.MerchantDocument
	if err := result.Decode(&doc); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*salemongo.MerchantDocument, error) {
	if request.DocumentID == nil {
		return nil, sharedErrors.ErrBadRequest.WithMessage("document id is required")
	}
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"document_id": *request.DocumentID},
		bson.M{"$set": bson.M{
			"status":     request.Status,
			"note":       request.Note,
			"updated_at": now,
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var doc salemongo.MerchantDocument
	if err := result.Decode(&doc); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) TrashedMerchantDocument(ctx context.Context, documentID int) (*salemongo.MerchantDocument, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"document_id": documentID, "deleted_at": nil},
		bson.M{"$set": bson.M{"deleted_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var doc salemongo.MerchantDocument
	if err := result.Decode(&doc); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) RestoreMerchantDocument(ctx context.Context, documentID int) (*salemongo.MerchantDocument, error) {
	now := time.Now()
	result := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"document_id": documentID, "deleted_at": bson.M{"$ne": nil}},
		bson.M{"$unset": bson.M{"deleted_at": ""}, "$set": bson.M{"updated_at": now}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)
	var doc salemongo.MerchantDocument
	if err := result.Decode(&doc); err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) DeleteMerchantDocumentPermanent(ctx context.Context, documentID int) (bool, error) {
	result, err := r.coll.DeleteOne(ctx, bson.M{"document_id": documentID, "deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}
	return result.DeletedCount > 0, nil
}

func (r *merchantDocumentCommandRepository) RestoreAllMerchantDocument(ctx context.Context) (bool, error) {
	_, err := r.coll.UpdateMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}}, bson.M{"$unset": bson.M{"deleted_at": ""}})
	if err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}
	return true, nil
}

func (r *merchantDocumentCommandRepository) DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, error) {
	_, err := r.coll.DeleteMany(ctx, bson.M{"deleted_at": bson.M{"$ne": nil}})
	if err != nil {
		return false, sharedErrors.ErrInternal.WithInternal(err)
	}
	return true, nil
}
