package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-merchant/repository"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type MerchantQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchants) ([]*repository.MerchantResult, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*repository.MerchantResultDeleteAt, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*repository.MerchantResultDeleteAt, *int, error)
	FindById(ctx context.Context, merchantID int) (*salemongo.Merchant, error)
}

type MerchantCommandService interface {
	CreateMerchant(ctx context.Context, req *requests.CreateMerchantRequest) (*salemongo.Merchant, error)
	UpdateMerchant(ctx context.Context, req *requests.UpdateMerchantRequest) (*salemongo.Merchant, error)
	UpdateMerchantStatus(ctx context.Context, req *requests.UpdateMerchantStatusRequest) (*salemongo.Merchant, error)
	TrashedMerchant(ctx context.Context, merchantID int) (*salemongo.Merchant, error)
	RestoreMerchant(ctx context.Context, merchantID int) (*salemongo.Merchant, error)
	DeleteMerchantPermanent(ctx context.Context, merchantID int) (bool, error)
	RestoreAllMerchant(ctx context.Context) (bool, error)
	DeleteAllMerchantPermanent(ctx context.Context) (bool, error)
}

type MerchantDocumentQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*repository.MerchantDocumentResult, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*repository.MerchantDocumentResultDeleteAt, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*repository.MerchantDocumentResultDeleteAt, *int, error)
	FindById(ctx context.Context, documentID int) (*salemongo.MerchantDocument, error)
}

type MerchantDocumentCommandService interface {
	CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*salemongo.MerchantDocument, error)
	UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*salemongo.MerchantDocument, error)
	UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*salemongo.MerchantDocument, error)
	TrashedMerchantDocument(ctx context.Context, documentID int) (*salemongo.MerchantDocument, error)
	RestoreMerchantDocument(ctx context.Context, documentID int) (*salemongo.MerchantDocument, error)
	DeleteMerchantDocumentPermanent(ctx context.Context, documentID int) (bool, error)
	RestoreAllMerchantDocument(ctx context.Context) (bool, error)
	DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, error)
}
