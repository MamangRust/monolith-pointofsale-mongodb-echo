package handler

import (
	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"strconv"

	"github.com/MamangRust/monolith-point-of-sale-cashier/repository"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/convert"
)

// Map helpers
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

func mapResponseCashier(cashier *salemongo.Cashier) *pbcashiers.CashierResponse {
	if cashier == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	if !cashier.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(cashier.CreatedAt)
	}
	if !cashier.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(cashier.UpdatedAt)
	}
	return &pbcashiers.CashierResponse{
		Id:         int32(cashier.CashierID),
		MerchantId: int32(cashier.MerchantID),
		Name:       cashier.Name,
		CreatedAt:  createdAtStr,
		UpdatedAt:  updatedAtStr,
	}
}

func mapResponsesCashier(cashiers []*repository.CashierRow) []*pbcashiers.CashierResponse {
	var mappedCashiers []*pbcashiers.CashierResponse
	for _, cashier := range cashiers {
		var createdAtStr, updatedAtStr string
		if !cashier.CreatedAt.IsZero() {
			createdAtStr = convert.FormatTimePtr(cashier.CreatedAt)
		}
		if !cashier.UpdatedAt.IsZero() {
			updatedAtStr = convert.FormatTimePtr(cashier.UpdatedAt)
		}
		mappedCashiers = append(mappedCashiers, &pbcashiers.CashierResponse{
			Id:         int32(cashier.CashierID),
			MerchantId: int32(cashier.MerchantID),
			Name:       cashier.Name,
			CreatedAt:  createdAtStr,
			UpdatedAt:  updatedAtStr,
		})
	}
	return mappedCashiers
}

func mapResponsesCashierByMerchant(cashiers []*repository.CashierByMerchantRow) []*pbcashiers.CashierResponse {
	var mappedCashiers []*pbcashiers.CashierResponse
	for _, cashier := range cashiers {
		var createdAtStr, updatedAtStr string
		if !cashier.CreatedAt.IsZero() {
			createdAtStr = convert.FormatTimePtr(cashier.CreatedAt)
		}
		if !cashier.UpdatedAt.IsZero() {
			updatedAtStr = convert.FormatTimePtr(cashier.UpdatedAt)
		}
		mappedCashiers = append(mappedCashiers, &pbcashiers.CashierResponse{
			Id:         int32(cashier.CashierID),
			MerchantId: int32(cashier.MerchantID),
			Name:       cashier.Name,
			CreatedAt:  createdAtStr,
			UpdatedAt:  updatedAtStr,
		})
	}
	return mappedCashiers
}

func mapResponseCashierDeleteAt(cashier *salemongo.Cashier) *pbcashiers.CashierResponseDeleteAt {
	if cashier == nil {
		return nil
	}
	var createdAtStr, updatedAtStr string
	if !cashier.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(cashier.CreatedAt)
	}
	if !cashier.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(cashier.UpdatedAt)
	}
	var deletedAt *wrapperspb.StringValue
	if cashier.DeletedAt != nil {
		deletedAt = convert.TimeToWrappers(cashier.DeletedAt)
	}

	return &pbcashiers.CashierResponseDeleteAt{
		Id:         int32(cashier.CashierID),
		MerchantId: int32(cashier.MerchantID),
		Name:       cashier.Name,
		CreatedAt:  createdAtStr,
		UpdatedAt:  updatedAtStr,
		DeletedAt:  deletedAt,
	}
}

func mapResponsesCashierActive(cashiers []*repository.CashierActiveRow) []*pbcashiers.CashierResponseDeleteAt {
	var mappedCashiers []*pbcashiers.CashierResponseDeleteAt
	for _, cashier := range cashiers {
		var createdAtStr, updatedAtStr string
		if !cashier.CreatedAt.IsZero() {
			createdAtStr = convert.FormatTimePtr(cashier.CreatedAt)
		}
		if !cashier.UpdatedAt.IsZero() {
			updatedAtStr = convert.FormatTimePtr(cashier.UpdatedAt)
		}
		var deletedAt *wrapperspb.StringValue
		if cashier.DeletedAt != nil {
			deletedAt = wrapperspb.String(cashier.DeletedAt.Format("2006-01-02 15:04:05"))
		}
		mappedCashiers = append(mappedCashiers, &pbcashiers.CashierResponseDeleteAt{
			Id:         int32(cashier.CashierID),
			MerchantId: int32(cashier.MerchantID),
			Name:       cashier.Name,
			CreatedAt:  createdAtStr,
			UpdatedAt:  updatedAtStr,
			DeletedAt:  deletedAt,
		})
	}
	return mappedCashiers
}

func mapResponsesCashierTrashed(cashiers []*repository.CashierTrashedRow) []*pbcashiers.CashierResponseDeleteAt {
	var mappedCashiers []*pbcashiers.CashierResponseDeleteAt
	for _, cashier := range cashiers {
		var createdAtStr, updatedAtStr string
		if !cashier.CreatedAt.IsZero() {
			createdAtStr = convert.FormatTimePtr(cashier.CreatedAt)
		}
		if !cashier.UpdatedAt.IsZero() {
			updatedAtStr = convert.FormatTimePtr(cashier.UpdatedAt)
		}
		var deletedAt *wrapperspb.StringValue
		if cashier.DeletedAt != nil {
			deletedAt = wrapperspb.String(cashier.DeletedAt.Format("2006-01-02 15:04:05"))
		}
		mappedCashiers = append(mappedCashiers, &pbcashiers.CashierResponseDeleteAt{
			Id:         int32(cashier.CashierID),
			MerchantId: int32(cashier.MerchantID),
			Name:       cashier.Name,
			CreatedAt:  createdAtStr,
			UpdatedAt:  updatedAtStr,
			DeletedAt:  deletedAt,
		})
	}
	return mappedCashiers
}

func mapResponseCashierMonthlySale(cashier *repository.CashierMonthlySaleRow) *pbcashiers.CashierResponseMonthSales {
	if cashier == nil {
		return nil
	}
	return &pbcashiers.CashierResponseMonthSales{
		Month:       strconv.Itoa(int(cashier.Month)),
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierMonthlySales(c []*repository.CashierMonthlySaleRow) []*pbcashiers.CashierResponseMonthSales {
	var cashierRecords []*pbcashiers.CashierResponseMonthSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierMonthlySale(cashier))
	}
	return cashierRecords
}

func mapResponseCashierMonthlySaleById(cashier *repository.CashierMonthlySaleByIdRow) *pbcashiers.CashierResponseMonthSales {
	if cashier == nil {
		return nil
	}
	return &pbcashiers.CashierResponseMonthSales{
		Month:       strconv.Itoa(int(cashier.Month)),
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierMonthlySalesById(c []*repository.CashierMonthlySaleByIdRow) []*pbcashiers.CashierResponseMonthSales {
	var cashierRecords []*pbcashiers.CashierResponseMonthSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierMonthlySaleById(cashier))
	}
	return cashierRecords
}

func mapResponseCashierMonthlySaleByMerchant(cashier *repository.CashierMonthlySaleByMerchantRow) *pbcashiers.CashierResponseMonthSales {
	if cashier == nil {
		return nil
	}
	return &pbcashiers.CashierResponseMonthSales{
		Month:       strconv.Itoa(int(cashier.Month)),
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierMonthlySalesByMerchant(c []*repository.CashierMonthlySaleByMerchantRow) []*pbcashiers.CashierResponseMonthSales {
	var cashierRecords []*pbcashiers.CashierResponseMonthSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierMonthlySaleByMerchant(cashier))
	}
	return cashierRecords
}

func mapResponseCashierYearlySale(cashier *repository.CashierYearlySaleRow) *pbcashiers.CashierResponseYearSales {
	if cashier == nil {
		return nil
	}
	return &pbcashiers.CashierResponseYearSales{
		Year:        strconv.Itoa(int(cashier.Year)),
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierYearlySales(c []*repository.CashierYearlySaleRow) []*pbcashiers.CashierResponseYearSales {
	var cashierRecords []*pbcashiers.CashierResponseYearSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierYearlySale(cashier))
	}
	return cashierRecords
}

func mapResponseCashierYearlySaleById(cashier *repository.CashierYearlySaleRow) *pbcashiers.CashierResponseYearSales {
	if cashier == nil {
		return nil
	}
	return &pbcashiers.CashierResponseYearSales{
		Year:        strconv.Itoa(int(cashier.Year)),
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierYearlySalesById(c []*repository.CashierYearlySaleRow) []*pbcashiers.CashierResponseYearSales {
	var cashierRecords []*pbcashiers.CashierResponseYearSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierYearlySaleById(cashier))
	}
	return cashierRecords
}

func mapResponseCashierYearlySaleByMerchant(cashier *repository.CashierYearlySaleRow) *pbcashiers.CashierResponseYearSales {
	if cashier == nil {
		return nil
	}
	return &pbcashiers.CashierResponseYearSales{
		Year:        strconv.Itoa(int(cashier.Year)),
		CashierId:   int32(cashier.CashierID),
		CashierName: cashier.CashierName,
		OrderCount:  int32(cashier.OrderCount),
		TotalSales:  int32(cashier.TotalSales),
	}
}

func mapResponsesCashierYearlySalesByMerchant(c []*repository.CashierYearlySaleRow) []*pbcashiers.CashierResponseYearSales {
	var cashierRecords []*pbcashiers.CashierResponseYearSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, mapResponseCashierYearlySaleByMerchant(cashier))
	}
	return cashierRecords
}

func mapResponseCashierMonthlyTotalSales(c []*repository.CashierMonthlyTotalSalesRow) []*pbcashiers.CashierResponseMonthTotalSales {
	var cashierRecords []*pbcashiers.CashierResponseMonthTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pbcashiers.CashierResponseMonthTotalSales{
			Year:       strconv.Itoa(int(cashier.Year)),
			Month:      strconv.Itoa(int(cashier.Month)),
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}

func mapResponseCashierMonthlyTotalSalesById(c []*repository.CashierMonthlyTotalSalesRow) []*pbcashiers.CashierResponseMonthTotalSales {
	var cashierRecords []*pbcashiers.CashierResponseMonthTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pbcashiers.CashierResponseMonthTotalSales{
			Year:       strconv.Itoa(int(cashier.Year)),
			Month:      strconv.Itoa(int(cashier.Month)),
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}

func mapResponseCashierMonthlyTotalSalesByMerchant(c []*repository.CashierMonthlyTotalSalesRow) []*pbcashiers.CashierResponseMonthTotalSales {
	var cashierRecords []*pbcashiers.CashierResponseMonthTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pbcashiers.CashierResponseMonthTotalSales{
			Year:       strconv.Itoa(int(cashier.Year)),
			Month:      strconv.Itoa(int(cashier.Month)),
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}

func mapResponseCashierYearlyTotalSales(c []*repository.CashierYearlyTotalSalesRow) []*pbcashiers.CashierResponseYearTotalSales {
	var cashierRecords []*pbcashiers.CashierResponseYearTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pbcashiers.CashierResponseYearTotalSales{
			Year:       strconv.Itoa(int(cashier.Year)),
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}

func mapResponseCashierYearlyTotalSalesById(c []*repository.CashierYearlyTotalSalesRow) []*pbcashiers.CashierResponseYearTotalSales {
	var cashierRecords []*pbcashiers.CashierResponseYearTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pbcashiers.CashierResponseYearTotalSales{
			Year:       strconv.Itoa(int(cashier.Year)),
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}

func mapResponseCashierYearlyTotalSalesByMerchant(c []*repository.CashierYearlyTotalSalesRow) []*pbcashiers.CashierResponseYearTotalSales {
	var cashierRecords []*pbcashiers.CashierResponseYearTotalSales
	for _, cashier := range c {
		cashierRecords = append(cashierRecords, &pbcashiers.CashierResponseYearTotalSales{
			Year:       strconv.Itoa(int(cashier.Year)),
			TotalSales: cashier.TotalSales,
		})
	}
	return cashierRecords
}
