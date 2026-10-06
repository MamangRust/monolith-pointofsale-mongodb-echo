package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

// MerchantResult is the result type for paginated merchant queries.
type MerchantResult struct {
	MerchantID   int32
	UserID       int32
	Name         string
	Description  *string
	Address      *string
	ContactEmail *string
	ContactPhone *string
	Status       string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	DeletedAt    *time.Time
	TotalCount   int64
}

// MerchantResultDeleteAt is the result type for paginated merchant queries with deleted_at.
type MerchantResultDeleteAt struct {
	MerchantID   int32
	UserID       int32
	Name         string
	Description  *string
	Address      *string
	ContactEmail *string
	ContactPhone *string
	Status       string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	DeletedAt    *time.Time
	TotalCount   int64
}

// MerchantDocumentResult is the result type for paginated merchant document queries.
type MerchantDocumentResult struct {
	DocumentID   int32
	MerchantID   int32
	DocumentType string
	DocumentUrl  string
	Status       string
	Note         *string
	UploadedAt   string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	DeletedAt    *time.Time
	TotalCount   int64
}

// MerchantDocumentResultDeleteAt is the result type for paginated merchant document queries with deleted_at.
type MerchantDocumentResultDeleteAt struct {
	DocumentID   int32
	MerchantID   int32
	DocumentType string
	DocumentUrl  string
	Status       string
	Note         *string
	UploadedAt   string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	DeletedAt    *time.Time
	TotalCount   int64
}

type MerchantDocumentQueryRepository interface {
	FindAllDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResult, *int, error)
	FindById(ctx context.Context, id int) (*salemongo.MerchantDocument, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResultDeleteAt, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*MerchantDocumentResultDeleteAt, *int, error)
}

type MerchantDocumentCommandRepository interface {
	CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*salemongo.MerchantDocument, error)
	UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*salemongo.MerchantDocument, error)
	UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*salemongo.MerchantDocument, error)
	TrashedMerchantDocument(ctx context.Context, merchantDocumentID int) (*salemongo.MerchantDocument, error)
	RestoreMerchantDocument(ctx context.Context, merchantDocumentID int) (*salemongo.MerchantDocument, error)
	DeleteMerchantDocumentPermanent(ctx context.Context, merchantDocumentID int) (bool, error)
	RestoreAllMerchantDocument(ctx context.Context) (bool, error)
	DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, error)
}

type MerchantQueryRepository interface {
	FindAllMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*MerchantResult, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*MerchantResultDeleteAt, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*MerchantResultDeleteAt, *int, error)
	FindById(ctx context.Context, userID int) (*salemongo.Merchant, error)
}

type MerchantCommandRepository interface {
	CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*salemongo.Merchant, error)
	UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*salemongo.Merchant, error)
	UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*salemongo.Merchant, error)
	TrashedMerchant(ctx context.Context, merchantID int) (*salemongo.Merchant, error)
	RestoreMerchant(ctx context.Context, merchantID int) (*salemongo.Merchant, error)
	DeleteMerchantPermanent(ctx context.Context, merchantID int) (bool, error)
	RestoreAllMerchant(ctx context.Context) (bool, error)
	DeleteAllMerchantPermanent(ctx context.Context) (bool, error)
}

type UserQueryRepository interface {
	FindById(ctx context.Context, userID int) (*salemongo.User, error)
}

// sqlc row type aliases for handler compatibility
type MerchantRow = MerchantResult
type MerchantActiveRow = MerchantResultDeleteAt
type MerchantTrashedRow = MerchantResultDeleteAt
type MerchantDocumentRow = MerchantDocumentResult
type MerchantDocumentActiveRow = MerchantDocumentResultDeleteAt
type MerchantDocumentTrashedRow = MerchantDocumentResultDeleteAt
