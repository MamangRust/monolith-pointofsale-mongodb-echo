package handler

import (
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pbtransactions "github.com/MamangRust/monolith-point-of-sale-pb/transactions"
	"google.golang.org/protobuf/types/known/wrapperspb"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/convert"
	"github.com/MamangRust/monolith-point-of-sale-transacton/repository"
)

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

func mapResponseTransaction(transaction *salemongo.Transaction) *pbtransactions.TransactionResponse {
	if transaction == nil {
		return nil
	}
	var createdAtStr string
	if transaction.CreatedAt != nil && !transaction.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(transaction.CreatedAt)
	}
	var updatedAtStr string
	if transaction.UpdatedAt != nil && !transaction.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(transaction.UpdatedAt)
	}
	var changeAmount int32
	if transaction.ChangeAmount != nil {
		changeAmount = *transaction.ChangeAmount
	}
	return &pbtransactions.TransactionResponse{
		Id:            transaction.TransactionID,
		OrderId:       transaction.OrderID,
		MerchantId:    transaction.MerchantID,
		PaymentMethod: transaction.PaymentMethod,
		Amount:        transaction.Amount,
		ChangeAmount:  changeAmount,
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     createdAtStr,
		UpdatedAt:     updatedAtStr,
	}
}

func mapResponsesTransaction(transactions []*repository.TransactionRow) []*pbtransactions.TransactionResponse {
	var mappedTransactions []*pbtransactions.TransactionResponse
	for _, t := range transactions {
		if t == nil {
			continue
		}
		var createdAtStr string
		if t.CreatedAt != nil && !t.CreatedAt.IsZero() {
			createdAtStr = convert.FormatTimePtr(t.CreatedAt)
		}
		var updatedAtStr string
		if t.UpdatedAt != nil && !t.UpdatedAt.IsZero() {
			updatedAtStr = convert.FormatTimePtr(t.UpdatedAt)
		}
		var changeAmount int32
		if t.ChangeAmount != nil {
			changeAmount = *t.ChangeAmount
		}
		mappedTransactions = append(mappedTransactions, &pbtransactions.TransactionResponse{
			Id:            t.TransactionID,
			OrderId:       t.OrderID,
			MerchantId:    t.MerchantID,
			PaymentMethod: t.PaymentMethod,
			Amount:        t.Amount,
			ChangeAmount:  changeAmount,
			PaymentStatus: t.PaymentStatus,
			CreatedAt:     createdAtStr,
			UpdatedAt:     updatedAtStr,
		})
	}
	return mappedTransactions
}

func mapResponsesTransactionByMerchant(transactions []*repository.TransactionByMerchantRow) []*pbtransactions.TransactionResponse {
	var mappedTransactions []*pbtransactions.TransactionResponse
	for _, t := range transactions {
		if t == nil {
			continue
		}
		var createdAtStr string
		if t.CreatedAt != nil && !t.CreatedAt.IsZero() {
			createdAtStr = convert.FormatTimePtr(t.CreatedAt)
		}
		var updatedAtStr string
		if t.UpdatedAt != nil && !t.UpdatedAt.IsZero() {
			updatedAtStr = convert.FormatTimePtr(t.UpdatedAt)
		}
		var changeAmount int32
		if t.ChangeAmount != nil {
			changeAmount = *t.ChangeAmount
		}
		mappedTransactions = append(mappedTransactions, &pbtransactions.TransactionResponse{
			Id:            t.TransactionID,
			OrderId:       t.OrderID,
			MerchantId:    t.MerchantID,
			PaymentMethod: t.PaymentMethod,
			Amount:        t.Amount,
			ChangeAmount:  changeAmount,
			PaymentStatus: t.PaymentStatus,
			CreatedAt:     createdAtStr,
			UpdatedAt:     updatedAtStr,
		})
	}
	return mappedTransactions
}

func mapResponseTransactionDeleteAt(transaction *salemongo.Transaction) *pbtransactions.TransactionResponseDeleteAt {
	if transaction == nil {
		return nil
	}
	var createdAtStr string
	if transaction.CreatedAt != nil && !transaction.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(transaction.CreatedAt)
	}
	var updatedAtStr string
	if transaction.UpdatedAt != nil && !transaction.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(transaction.UpdatedAt)
	}
	var deletedAt *wrapperspb.StringValue
	if transaction.DeletedAt != nil {
		deletedAt = convert.TimeToWrappers(transaction.DeletedAt)
	}
	var changeAmount int32
	if transaction.ChangeAmount != nil {
		changeAmount = *transaction.ChangeAmount
	}

	return &pbtransactions.TransactionResponseDeleteAt{
		Id:            transaction.TransactionID,
		OrderId:       transaction.OrderID,
		MerchantId:    transaction.MerchantID,
		PaymentMethod: transaction.PaymentMethod,
		Amount:        transaction.Amount,
		ChangeAmount:  changeAmount,
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     createdAtStr,
		UpdatedAt:     updatedAtStr,
		DeletedAt:     deletedAt,
	}
}

func mapResponsesTransactionActive(transactions []*repository.TransactionActiveRow) []*pbtransactions.TransactionResponseDeleteAt {
	var mappedTransactions []*pbtransactions.TransactionResponseDeleteAt
	for _, t := range transactions {
		if t == nil {
			continue
		}
		var createdAtStr string
		if t.CreatedAt != nil && !t.CreatedAt.IsZero() {
			createdAtStr = convert.FormatTimePtr(t.CreatedAt)
		}
		var updatedAtStr string
		if t.UpdatedAt != nil && !t.UpdatedAt.IsZero() {
			updatedAtStr = convert.FormatTimePtr(t.UpdatedAt)
		}
		var deletedAt *wrapperspb.StringValue
		if t.DeletedAt != nil {
			deletedAt = wrapperspb.String(t.DeletedAt.Format("2006-01-02 15:04:05"))
		}
		var changeAmount int32
		if t.ChangeAmount != nil {
			changeAmount = *t.ChangeAmount
		}
		mappedTransactions = append(mappedTransactions, &pbtransactions.TransactionResponseDeleteAt{
			Id:            t.TransactionID,
			OrderId:       t.OrderID,
			MerchantId:    t.MerchantID,
			PaymentMethod: t.PaymentMethod,
			Amount:        t.Amount,
			ChangeAmount:  changeAmount,
			PaymentStatus: t.PaymentStatus,
			CreatedAt:     createdAtStr,
			UpdatedAt:     updatedAtStr,
			DeletedAt:     deletedAt,
		})
	}
	return mappedTransactions
}

func mapResponsesTransactionTrashed(transactions []*repository.TransactionTrashedRow) []*pbtransactions.TransactionResponseDeleteAt {
	var mappedTransactions []*pbtransactions.TransactionResponseDeleteAt
	for _, t := range transactions {
		if t == nil {
			continue
		}
		var createdAtStr string
		if t.CreatedAt != nil && !t.CreatedAt.IsZero() {
			createdAtStr = convert.FormatTimePtr(t.CreatedAt)
		}
		var updatedAtStr string
		if t.UpdatedAt != nil && !t.UpdatedAt.IsZero() {
			updatedAtStr = convert.FormatTimePtr(t.UpdatedAt)
		}
		var deletedAt *wrapperspb.StringValue
		if t.DeletedAt != nil {
			deletedAt = wrapperspb.String(t.DeletedAt.Format("2006-01-02 15:04:05"))
		}
		var changeAmount int32
		if t.ChangeAmount != nil {
			changeAmount = *t.ChangeAmount
		}
		mappedTransactions = append(mappedTransactions, &pbtransactions.TransactionResponseDeleteAt{
			Id:            t.TransactionID,
			OrderId:       t.OrderID,
			MerchantId:    t.MerchantID,
			PaymentMethod: t.PaymentMethod,
			Amount:        t.Amount,
			ChangeAmount:  changeAmount,
			PaymentStatus: t.PaymentStatus,
			CreatedAt:     createdAtStr,
			UpdatedAt:     updatedAtStr,
			DeletedAt:     deletedAt,
		})
	}
	return mappedTransactions
}

func mapResponseTransactionMonthAmountSuccess(row *repository.TransactionMonthlyAmountSuccessRow) *pbtransactions.TransactionMonthlyAmountSuccess {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionMonthlyAmountSuccess{
		Year:         row.Year,
		Month:        row.Month,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyAmountSuccess(rows []*repository.TransactionMonthlyAmountSuccessRow) []*pbtransactions.TransactionMonthlyAmountSuccess {
	var transaction []*pbtransactions.TransactionMonthlyAmountSuccess
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthAmountSuccess(row))
	}
	return transaction
}

func mapResponseTransactionYearAmountSuccess(row *repository.TransactionYearlyAmountSuccessRow) *pbtransactions.TransactionYearlyAmountSuccess {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionYearlyAmountSuccess{
		Year:         row.Year,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyAmountSuccess(rows []*repository.TransactionYearlyAmountSuccessRow) []*pbtransactions.TransactionYearlyAmountSuccess {
	var transaction []*pbtransactions.TransactionYearlyAmountSuccess
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearAmountSuccess(row))
	}
	return transaction
}

func mapResponseTransactionMonthAmountFailed(row *repository.TransactionMonthlyAmountFailedRow) *pbtransactions.TransactionMonthlyAmountFailed {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionMonthlyAmountFailed{
		Year:        row.Year,
		Month:       row.Month,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyAmountFailed(rows []*repository.TransactionMonthlyAmountFailedRow) []*pbtransactions.TransactionMonthlyAmountFailed {
	var transaction []*pbtransactions.TransactionMonthlyAmountFailed
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthAmountFailed(row))
	}
	return transaction
}

func mapResponseTransactionYearAmountFailed(row *repository.TransactionYearlyAmountFailedRow) *pbtransactions.TransactionYearlyAmountFailed {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionYearlyAmountFailed{
		Year:        row.Year,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyAmountFailed(rows []*repository.TransactionYearlyAmountFailedRow) []*pbtransactions.TransactionYearlyAmountFailed {
	var transaction []*pbtransactions.TransactionYearlyAmountFailed
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearAmountFailed(row))
	}
	return transaction
}

func mapResponseTransactionMonthAmountSuccessByMerchant(row *repository.TransactionMonthlyAmountSuccessByMerchantRow) *pbtransactions.TransactionMonthlyAmountSuccess {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionMonthlyAmountSuccess{
		Year:         row.Year,
		Month:        row.Month,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyAmountSuccessByMerchant(rows []*repository.TransactionMonthlyAmountSuccessByMerchantRow) []*pbtransactions.TransactionMonthlyAmountSuccess {
	var transaction []*pbtransactions.TransactionMonthlyAmountSuccess
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthAmountSuccessByMerchant(row))
	}
	return transaction
}

func mapResponseTransactionYearAmountSuccessByMerchant(row *repository.TransactionYearlyAmountSuccessByMerchantRow) *pbtransactions.TransactionYearlyAmountSuccess {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionYearlyAmountSuccess{
		Year:         row.Year,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyAmountSuccessByMerchant(rows []*repository.TransactionYearlyAmountSuccessByMerchantRow) []*pbtransactions.TransactionYearlyAmountSuccess {
	var transaction []*pbtransactions.TransactionYearlyAmountSuccess
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearAmountSuccessByMerchant(row))
	}
	return transaction
}

func mapResponseTransactionMonthAmountFailedByMerchant(row *repository.TransactionMonthlyAmountFailedByMerchantRow) *pbtransactions.TransactionMonthlyAmountFailed {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionMonthlyAmountFailed{
		Year:        row.Year,
		Month:       row.Month,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyAmountFailedByMerchant(rows []*repository.TransactionMonthlyAmountFailedByMerchantRow) []*pbtransactions.TransactionMonthlyAmountFailed {
	var transaction []*pbtransactions.TransactionMonthlyAmountFailed
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthAmountFailedByMerchant(row))
	}
	return transaction
}

func mapResponseTransactionYearAmountFailedByMerchant(row *repository.TransactionYearlyAmountFailedByMerchantRow) *pbtransactions.TransactionYearlyAmountFailed {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionYearlyAmountFailed{
		Year:        row.Year,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyAmountFailedByMerchant(rows []*repository.TransactionYearlyAmountFailedByMerchantRow) []*pbtransactions.TransactionYearlyAmountFailed {
	var transaction []*pbtransactions.TransactionYearlyAmountFailed
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearAmountFailedByMerchant(row))
	}
	return transaction
}

func mapResponseTransactionMonthMethodSuccess(row *repository.TransactionMonthlyMethodSuccessRow) *pbtransactions.TransactionMonthlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionMonthlyMethod{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyMethodSuccess(rows []*repository.TransactionMonthlyMethodSuccessRow) []*pbtransactions.TransactionMonthlyMethod {
	var transaction []*pbtransactions.TransactionMonthlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthMethodSuccess(row))
	}
	return transaction
}

func mapResponseTransactionYearlyMethodSuccess(row *repository.TransactionYearlyMethodSuccessRow) *pbtransactions.TransactionYearlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionYearlyMethod{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyMethodSuccess(rows []*repository.TransactionYearlyMethodSuccessRow) []*pbtransactions.TransactionYearlyMethod {
	var transaction []*pbtransactions.TransactionYearlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearlyMethodSuccess(row))
	}
	return transaction
}

func mapResponseTransactionMonthMethodFailed(row *repository.TransactionMonthlyMethodFailedRow) *pbtransactions.TransactionMonthlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionMonthlyMethod{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyMethodFailed(rows []*repository.TransactionMonthlyMethodFailedRow) []*pbtransactions.TransactionMonthlyMethod {
	var transaction []*pbtransactions.TransactionMonthlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthMethodFailed(row))
	}
	return transaction
}

func mapResponseTransactionYearlyMethodFailed(row *repository.TransactionYearlyMethodFailedRow) *pbtransactions.TransactionYearlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionYearlyMethod{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyMethodFailed(rows []*repository.TransactionYearlyMethodFailedRow) []*pbtransactions.TransactionYearlyMethod {
	var transaction []*pbtransactions.TransactionYearlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearlyMethodFailed(row))
	}
	return transaction
}

func mapResponseTransactionMonthMethodByMerchantSuccess(row *repository.TransactionMonthlyMethodByMerchantSuccessRow) *pbtransactions.TransactionMonthlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionMonthlyMethod{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyMethodByMerchantSuccess(rows []*repository.TransactionMonthlyMethodByMerchantSuccessRow) []*pbtransactions.TransactionMonthlyMethod {
	var transaction []*pbtransactions.TransactionMonthlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthMethodByMerchantSuccess(row))
	}
	return transaction
}

func mapResponseTransactionYearlyMethodByMerchantSuccess(row *repository.TransactionYearlyMethodByMerchantSuccessRow) *pbtransactions.TransactionYearlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionYearlyMethod{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyMethodByMerchantSuccess(rows []*repository.TransactionYearlyMethodByMerchantSuccessRow) []*pbtransactions.TransactionYearlyMethod {
	var transaction []*pbtransactions.TransactionYearlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearlyMethodByMerchantSuccess(row))
	}
	return transaction
}

func mapResponseTransactionMonthMethodByMerchantFailed(row *repository.TransactionMonthlyMethodByMerchantFailedRow) *pbtransactions.TransactionMonthlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionMonthlyMethod{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyMethodByMerchantFailed(rows []*repository.TransactionMonthlyMethodByMerchantFailedRow) []*pbtransactions.TransactionMonthlyMethod {
	var transaction []*pbtransactions.TransactionMonthlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthMethodByMerchantFailed(row))
	}
	return transaction
}

func mapResponseTransactionYearlyMethodByMerchantFailed(row *repository.TransactionYearlyMethodByMerchantFailedRow) *pbtransactions.TransactionYearlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionYearlyMethod{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyMethodByMerchantFailed(rows []*repository.TransactionYearlyMethodByMerchantFailedRow) []*pbtransactions.TransactionYearlyMethod {
	var transaction []*pbtransactions.TransactionYearlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearlyMethodByMerchantFailed(row))
	}
	return transaction
}
