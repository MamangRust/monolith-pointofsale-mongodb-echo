package handler

import (
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-merchant/repository"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
)

func mapSqlNullString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func mapSqlNullTimeGorm(t *time.Time) string {
	if t != nil {
		return t.Format("2006-01-02 15:04:05")
	}
	return ""
}

func mapSqlNullTimePtr(t *time.Time) string {
	if t != nil && !t.IsZero() {
		return t.Format("2006-01-02 15:04:05")
	}
	return ""
}

func mapSqlNullTime(t time.Time) string {
	if !t.IsZero() {
		return t.Format("2006-01-02 15:04:05")
	}
	return ""
}

func mapPaginationMeta(meta *pbcommon.PaginationMeta) *pbcommon.PaginationMeta {
	if meta == nil {
		return nil
	}
	return &pbcommon.PaginationMeta{
		CurrentPage:  meta.CurrentPage,
		PageSize:     meta.PageSize,
		TotalPages:   meta.TotalPages,
		TotalRecords: meta.TotalRecords,
	}
}

func mapResponseMerchant(merchant *salemongo.Merchant) *pbmerchants.MerchantResponse {
	if merchant == nil {
		return nil
	}
	return &pbmerchants.MerchantResponse{
		Id:           int32(merchant.MerchantID),
		UserId:       int32(merchant.UserID),
		Name:         merchant.Name,
		Description:  mapSqlNullString(merchant.Description),
		Address:      mapSqlNullString(merchant.Address),
		ContactEmail: mapSqlNullString(merchant.ContactEmail),
		ContactPhone: mapSqlNullString(merchant.ContactPhone),
		Status:       merchant.Status,
		CreatedAt:    mapSqlNullTimePtr(merchant.CreatedAt),
		UpdatedAt:    mapSqlNullTimePtr(merchant.UpdatedAt),
	}
}

func mapResponsesGetMerchantsRow(merchants []*repository.MerchantRow) []*pbmerchants.MerchantResponse {
	var mapped []*pbmerchants.MerchantResponse
	for _, m := range merchants {
		mapped = append(mapped, &pbmerchants.MerchantResponse{
			Id:           int32(m.MerchantID),
			UserId:       int32(m.UserID),
			Name:         m.Name,
			Description:  mapSqlNullString(m.Description),
			Address:      mapSqlNullString(m.Address),
			ContactEmail: mapSqlNullString(m.ContactEmail),
			ContactPhone: mapSqlNullString(m.ContactPhone),
			Status:       m.Status,
			CreatedAt:    mapSqlNullTimePtr(m.CreatedAt),
			UpdatedAt:    mapSqlNullTimePtr(m.UpdatedAt),
		})
	}
	return mapped
}

func mapResponseMerchantDeleteAt(merchant *salemongo.Merchant) *pbmerchants.MerchantResponseDeleteAt {
	if merchant == nil {
		return nil
	}
	return &pbmerchants.MerchantResponseDeleteAt{
		Id:           int32(merchant.MerchantID),
		UserId:       int32(merchant.UserID),
		Name:         merchant.Name,
		Description:  mapSqlNullString(merchant.Description),
		Address:      mapSqlNullString(merchant.Address),
		ContactEmail: mapSqlNullString(merchant.ContactEmail),
		ContactPhone: mapSqlNullString(merchant.ContactPhone),
		Status:       merchant.Status,
		CreatedAt:    mapSqlNullTimePtr(merchant.CreatedAt),
		UpdatedAt:    mapSqlNullTimePtr(merchant.UpdatedAt),
		DeletedAt:    mapSqlNullTimeGorm(merchant.DeletedAt),
	}
}

func mapResponsesGetMerchantsActiveRow(merchants []*repository.MerchantActiveRow) []*pbmerchants.MerchantResponseDeleteAt {
	var mapped []*pbmerchants.MerchantResponseDeleteAt
	for _, m := range merchants {
		mapped = append(mapped, &pbmerchants.MerchantResponseDeleteAt{
			Id:           int32(m.MerchantID),
			UserId:       int32(m.UserID),
			Name:         m.Name,
			Description:  mapSqlNullString(m.Description),
			Address:      mapSqlNullString(m.Address),
			ContactEmail: mapSqlNullString(m.ContactEmail),
			ContactPhone: mapSqlNullString(m.ContactPhone),
			Status:       m.Status,
			CreatedAt:    mapSqlNullTimePtr(m.CreatedAt),
			UpdatedAt:    mapSqlNullTimePtr(m.UpdatedAt),
			DeletedAt:    mapSqlNullTimePtr(m.DeletedAt),
		})
	}
	return mapped
}

func mapResponsesGetMerchantsTrashedRow(merchants []*repository.MerchantTrashedRow) []*pbmerchants.MerchantResponseDeleteAt {
	var mapped []*pbmerchants.MerchantResponseDeleteAt
	for _, m := range merchants {
		mapped = append(mapped, &pbmerchants.MerchantResponseDeleteAt{
			Id:           int32(m.MerchantID),
			UserId:       int32(m.UserID),
			Name:         m.Name,
			Description:  mapSqlNullString(m.Description),
			Address:      mapSqlNullString(m.Address),
			ContactEmail: mapSqlNullString(m.ContactEmail),
			ContactPhone: mapSqlNullString(m.ContactPhone),
			Status:       m.Status,
			CreatedAt:    mapSqlNullTimePtr(m.CreatedAt),
			UpdatedAt:    mapSqlNullTimePtr(m.UpdatedAt),
			DeletedAt:    mapSqlNullTimePtr(m.DeletedAt),
		})
	}
	return mapped
}
