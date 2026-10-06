package service

import (
	"context"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-transacton/repository"
)

// EmailEventPublisher publishes email notification events to Kafka topics.
type EmailEventPublisher interface {
	SendMessage(ctx context.Context, topic string, key string, value []byte) error
}

type TransactionQueryService interface {
	FindAllTransactions(ctx context.Context, req *requests.FindAllTransaction) ([]*repository.TransactionResult, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*repository.TransactionByMerchantResult, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllTransaction) ([]*repository.TransactionResultDeleteAt, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllTransaction) ([]*repository.TransactionResultDeleteAt, *int, error)
	FindById(ctx context.Context, transactionID int) (*salemongo.Transaction, error)
	FindByOrderId(ctx context.Context, orderID int) (*salemongo.Transaction, error)
}

type TransactionCommandService interface {
	CreateTransaction(ctx context.Context, req *requests.CreateTransactionRequest) (*salemongo.Transaction, error)
	UpdateTransaction(ctx context.Context, req *requests.UpdateTransactionRequest) (*salemongo.Transaction, error)
	TrashedTransaction(ctx context.Context, transaction_id int) (*salemongo.Transaction, error)
	RestoreTransaction(ctx context.Context, transaction_id int) (*salemongo.Transaction, error)
	DeleteTransactionPermanently(ctx context.Context, transactionID int) (bool, error)
	RestoreAllTransactions(ctx context.Context) (bool, error)
	DeleteAllTransactionPermanent(ctx context.Context) (bool, error)
}

// TransactionStatsService exposes aggregate transaction statistics.
type TransactionStatsService interface {
	FindMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*salemongo.TransactionMonthlyAmountSuccessRow, error)
	FindYearlyAmountSuccess(ctx context.Context, year int) ([]*salemongo.TransactionYearlyAmountSuccessRow, error)
	FindMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*salemongo.TransactionMonthlyAmountFailedRow, error)
	FindYearlyAmountFailed(ctx context.Context, year int) ([]*salemongo.TransactionYearlyAmountFailedRow, error)
	FindMonthlyMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*salemongo.TransactionMonthlyMethodSuccessRow, error)
	FindYearlyMethodSuccess(ctx context.Context, year int) ([]*salemongo.TransactionYearlyMethodSuccessRow, error)
	FindMonthlyMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*salemongo.TransactionMonthlyMethodFailedRow, error)
	FindYearlyMethodFailed(ctx context.Context, year int) ([]*salemongo.TransactionYearlyMethodFailedRow, error)
}

// TransactionStatsByMerchantQueryService exposes per-merchant transaction statistics.
type TransactionStatsByMerchantQueryService interface {
	FindMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*salemongo.TransactionMonthlyAmountSuccessByMerchantRow, error)
	FindYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*salemongo.TransactionYearlyAmountSuccessByMerchantRow, error)
	FindMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*salemongo.TransactionMonthlyAmountFailedByMerchantRow, error)
	FindYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*salemongo.TransactionYearlyAmountFailedByMerchantRow, error)
	FindMonthlyMethodByMerchantSuccess(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*salemongo.TransactionMonthlyMethodByMerchantSuccessRow, error)
	FindYearlyMethodByMerchantSuccess(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*salemongo.TransactionYearlyMethodByMerchantSuccessRow, error)
	FindMonthlyMethodByMerchantFailed(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*salemongo.TransactionMonthlyMethodByMerchantFailedRow, error)
	FindYearlyMethodByMerchantFailed(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*salemongo.TransactionYearlyMethodByMerchantFailedRow, error)
}
