package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

// TransactionResult is the result type for paginated transaction queries.
type TransactionResult struct {
	TransactionID int32
	OrderID       int32
	MerchantID    int32
	PaymentMethod string
	Amount        int32
	ChangeAmount  *int32
	PaymentStatus string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
	TotalCount    int64
}

// TransactionResultDeleteAt is the result type with deleted_at.
type TransactionResultDeleteAt struct {
	TransactionID int32
	OrderID       int32
	MerchantID    int32
	PaymentMethod string
	Amount        int32
	ChangeAmount  *int32
	PaymentStatus string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
	DeletedAt     *time.Time
	TotalCount    int64
}

// TransactionByMerchantResult is the result type for merchant-filtered queries.
type TransactionByMerchantResult struct {
	TransactionID int32
	OrderID       int32
	MerchantID    int32
	PaymentMethod string
	Amount        int32
	ChangeAmount  *int32
	PaymentStatus string
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
	TotalCount    int64
}

type CashierQueryRepository interface {
	FindById(ctx context.Context, id int) (*salemongo.Cashier, error)
}

type MerchantQueryRepository interface {
	FindById(ctx context.Context, id int) (*salemongo.Merchant, error)
}

type OrderItemQueryRepository interface {
	FindOrderItemByOrder(ctx context.Context, order_id int) ([]*salemongo.OrderItem, error)
}

type OrderQueryRepository interface {
	FindById(ctx context.Context, id int) (*salemongo.Order, error)
}

type TransactionQueryRepository interface {
	FindAllTransactions(ctx context.Context, req *requests.FindAllTransaction) ([]*TransactionResult, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllTransaction) ([]*TransactionResultDeleteAt, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllTransaction) ([]*TransactionResultDeleteAt, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*TransactionByMerchantResult, *int, error)
	FindById(ctx context.Context, transaction_id int) (*salemongo.Transaction, error)
	FindByOrderId(ctx context.Context, order_id int) (*salemongo.Transaction, error)
}

type TransactionCommandRepository interface {
	CreateTransaction(ctx context.Context, request *requests.CreateTransactionRequest) (*salemongo.Transaction, error)
	UpdateTransaction(ctx context.Context, request *requests.UpdateTransactionRequest) (*salemongo.Transaction, error)
	TrashTransaction(ctx context.Context, transaction_id int) (*salemongo.Transaction, error)
	RestoreTransaction(ctx context.Context, transaction_id int) (*salemongo.Transaction, error)
	DeleteTransactionPermanently(ctx context.Context, transaction_id int) (bool, error)
	RestoreAllTransactions(ctx context.Context) (bool, error)
	DeleteAllTransactionPermanent(ctx context.Context) (bool, error)
}

// sqlc row type aliases for handler compatibility
type TransactionRow = TransactionResult
type TransactionActiveRow = TransactionResultDeleteAt
type TransactionTrashedRow = TransactionResultDeleteAt
type TransactionByMerchantRow = TransactionByMerchantResult
type TransactionMonthlyAmountSuccessRow = salemongo.TransactionMonthlyAmountSuccessRow
type TransactionMonthlyAmountSuccessByMerchantRow = salemongo.TransactionMonthlyAmountSuccessByMerchantRow
type TransactionMonthlyAmountFailedRow = salemongo.TransactionMonthlyAmountFailedRow
type TransactionMonthlyAmountFailedByMerchantRow = salemongo.TransactionMonthlyAmountFailedByMerchantRow
type TransactionMonthlyMethodSuccessRow = salemongo.TransactionMonthlyMethodSuccessRow
type TransactionMonthlyMethodSuccessByMerchantRow = salemongo.TransactionMonthlyMethodSuccessByMerchantRow
type TransactionMonthlyMethodFailedRow = salemongo.TransactionMonthlyMethodFailedRow
type TransactionMonthlyMethodFailedByMerchantRow = salemongo.TransactionMonthlyMethodFailedByMerchantRow
type TransactionYearlyAmountSuccessRow = salemongo.TransactionYearlyAmountSuccessRow
type TransactionYearlyAmountSuccessByMerchantRow = salemongo.TransactionYearlyAmountSuccessByMerchantRow
type TransactionYearlyAmountFailedRow = salemongo.TransactionYearlyAmountFailedRow
type TransactionYearlyAmountFailedByMerchantRow = salemongo.TransactionYearlyAmountFailedByMerchantRow
type TransactionYearlyMethodSuccessRow = salemongo.TransactionYearlyMethodSuccessRow
type TransactionYearlyMethodSuccessByMerchantRow = salemongo.TransactionYearlyMethodSuccessByMerchantRow
type TransactionYearlyMethodFailedRow = salemongo.TransactionYearlyMethodFailedRow
type TransactionYearlyMethodFailedByMerchantRow = salemongo.TransactionYearlyMethodFailedByMerchantRow

// TransactionStatsRepository provides aggregate transaction statistics.
type TransactionStatsRepository interface {
	GetMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*salemongo.TransactionMonthlyAmountSuccessRow, error)
	GetYearlyAmountSuccess(ctx context.Context, year int) ([]*salemongo.TransactionYearlyAmountSuccessRow, error)
	GetMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*salemongo.TransactionMonthlyAmountFailedRow, error)
	GetYearlyAmountFailed(ctx context.Context, year int) ([]*salemongo.TransactionYearlyAmountFailedRow, error)
	GetMonthlyTransactionMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*salemongo.TransactionMonthlyMethodSuccessRow, error)
	GetYearlyTransactionMethodSuccess(ctx context.Context, year int) ([]*salemongo.TransactionYearlyMethodSuccessRow, error)
	GetMonthlyTransactionMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*salemongo.TransactionMonthlyMethodFailedRow, error)
	GetYearlyTransactionMethodFailed(ctx context.Context, year int) ([]*salemongo.TransactionYearlyMethodFailedRow, error)
}

// TransactionStatsByMerchantRepository provides transaction statistics grouped by merchant.
type TransactionStatsByMerchantRepository interface {
	GetMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*salemongo.TransactionMonthlyAmountSuccessByMerchantRow, error)
	GetYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*salemongo.TransactionYearlyAmountSuccessByMerchantRow, error)
	GetMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*salemongo.TransactionMonthlyAmountFailedByMerchantRow, error)
	GetYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*salemongo.TransactionYearlyAmountFailedByMerchantRow, error)
	GetMonthlyTransactionMethodByMerchantSuccess(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*salemongo.TransactionMonthlyMethodByMerchantSuccessRow, error)
	GetYearlyTransactionMethodByMerchantSuccess(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*salemongo.TransactionYearlyMethodByMerchantSuccessRow, error)
	GetMonthlyTransactionMethodByMerchantFailed(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*salemongo.TransactionMonthlyMethodByMerchantFailedRow, error)
	GetYearlyTransactionMethodByMerchantFailed(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*salemongo.TransactionYearlyMethodByMerchantFailedRow, error)
}

type TransactionMonthlyMethodByMerchantSuccessRow = salemongo.TransactionMonthlyMethodByMerchantSuccessRow
type TransactionYearlyMethodByMerchantSuccessRow = salemongo.TransactionYearlyMethodByMerchantSuccessRow
type TransactionMonthlyMethodByMerchantFailedRow = salemongo.TransactionMonthlyMethodByMerchantFailedRow
type TransactionYearlyMethodByMerchantFailedRow = salemongo.TransactionYearlyMethodByMerchantFailedRow
