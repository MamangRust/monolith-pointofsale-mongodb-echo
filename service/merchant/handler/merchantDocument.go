package handler

import (
	pbmerchants_document "github.com/MamangRust/monolith-point-of-sale-pb/merchants_document"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/MamangRust/monolith-point-of-sale-merchant/repository"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/convert"
)

func mapMerchantDocument(doc *salemongo.MerchantDocument) *pbmerchants_document.MerchantDocument {
	if doc == nil {
		return nil
	}
	return &pbmerchants_document.MerchantDocument{
		DocumentId:   int32(doc.DocumentID),
		MerchantId:   int32(doc.MerchantID),
		DocumentType: doc.DocumentType,
		DocumentUrl:  doc.DocumentUrl,
		Status:       doc.Status,
		Note:         mapSqlNullString(doc.Note),
		UploadedAt:   mapSqlNullTimePtr(doc.CreatedAt),
		UpdatedAt:    mapSqlNullTimePtr(doc.UpdatedAt),
	}
}

func mapResponsesGetMerchantDocumentsRow(docs []*repository.MerchantDocumentRow) []*pbmerchants_document.MerchantDocument {
	var res []*pbmerchants_document.MerchantDocument
	for _, doc := range docs {
		res = append(res, &pbmerchants_document.MerchantDocument{
			DocumentId:   int32(doc.DocumentID),
			MerchantId:   int32(doc.MerchantID),
			DocumentType: doc.DocumentType,
			DocumentUrl:  doc.DocumentUrl,
			Status:       doc.Status,
			Note:         mapSqlNullString(doc.Note),
			UploadedAt:   mapSqlNullTimePtr(doc.CreatedAt),
			UpdatedAt:    mapSqlNullTimePtr(doc.UpdatedAt),
		})
	}
	return res
}

func mapResponsesGetActiveMerchantDocumentsRow(docs []*repository.MerchantDocumentActiveRow) []*pbmerchants_document.MerchantDocument {
	var res []*pbmerchants_document.MerchantDocument
	for _, doc := range docs {
		res = append(res, &pbmerchants_document.MerchantDocument{
			DocumentId:   int32(doc.DocumentID),
			MerchantId:   int32(doc.MerchantID),
			DocumentType: doc.DocumentType,
			DocumentUrl:  doc.DocumentUrl,
			Status:       doc.Status,
			Note:         mapSqlNullString(doc.Note),
			UploadedAt:   mapSqlNullTimePtr(doc.CreatedAt),
			UpdatedAt:    mapSqlNullTimePtr(doc.UpdatedAt),
		})
	}
	return res
}

func mapMerchantDocumentDeleteAt(doc *salemongo.MerchantDocument) *pbmerchants_document.MerchantDocumentDeleteAt {
	if doc == nil {
		return nil
	}
	var deletedAt *wrapperspb.StringValue
	if doc.DeletedAt != nil {
		deletedAt = convert.TimeToWrappers(doc.DeletedAt)
	}

	return &pbmerchants_document.MerchantDocumentDeleteAt{
		DocumentId:   int32(doc.DocumentID),
		MerchantId:   int32(doc.MerchantID),
		DocumentType: doc.DocumentType,
		DocumentUrl:  doc.DocumentUrl,
		Status:       doc.Status,
		Note:         mapSqlNullString(doc.Note),
		UploadedAt:   mapSqlNullTimePtr(doc.CreatedAt),
		UpdatedAt:    mapSqlNullTimePtr(doc.UpdatedAt),
		DeletedAt:    deletedAt,
	}
}

func mapResponsesGetTrashedMerchantDocumentsRow(docs []*repository.MerchantDocumentTrashedRow) []*pbmerchants_document.MerchantDocumentDeleteAt {
	var res []*pbmerchants_document.MerchantDocumentDeleteAt
	for _, doc := range docs {
		var deletedAt *wrapperspb.StringValue
		if doc.DeletedAt != nil {
			deletedAt = wrapperspb.String(doc.DeletedAt.Format("2006-01-02 15:04:05"))
		}
		res = append(res, &pbmerchants_document.MerchantDocumentDeleteAt{
			DocumentId:   int32(doc.DocumentID),
			MerchantId:   int32(doc.MerchantID),
			DocumentType: doc.DocumentType,
			DocumentUrl:  doc.DocumentUrl,
			Status:       doc.Status,
			Note:         mapSqlNullString(doc.Note),
			UploadedAt:   mapSqlNullTimePtr(doc.CreatedAt),
			UpdatedAt:    mapSqlNullTimePtr(doc.UpdatedAt),
			DeletedAt:    deletedAt,
		})
	}
	return res
}
