package mencache

import (
	"context"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-transacton/repository"
)

type TransactionQueryCache interface {
	GetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction) ([]*repository.TransactionResult, *int, bool)
	SetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction, data []*repository.TransactionResult, total *int)

	GetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*repository.TransactionByMerchantResult, *int, bool)
	SetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant, data []*repository.TransactionByMerchantResult, total *int)

	GetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction) ([]*repository.TransactionResultDeleteAt, *int, bool)
	SetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction, data []*repository.TransactionResultDeleteAt, total *int)

	GetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction) ([]*repository.TransactionResultDeleteAt, *int, bool)
	SetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction, data []*repository.TransactionResultDeleteAt, total *int)

	GetCachedTransactionCache(ctx context.Context, id int) (*salemongo.Transaction, bool)
	SetCachedTransactionCache(ctx context.Context, data *salemongo.Transaction)

	GetCachedTransactionByOrderId(ctx context.Context, orderID int) (*salemongo.Transaction, bool)
	SetCachedTransactionByOrderId(ctx context.Context, orderID int, data *salemongo.Transaction)
}

type TransactionCommandCache interface {
	DeleteTransactionCache(ctx context.Context, transactionID int)
	DeleteTransactionAllCache(ctx context.Context)
}

// TransactionStatsCache caches aggregate transaction statistics.
type TransactionStatsCache interface {
	GetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransaction) ([]*salemongo.TransactionMonthlyAmountSuccessRow, bool)
	SetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransaction, res []*salemongo.TransactionMonthlyAmountSuccessRow)
	GetCachedYearAmountSuccessCached(ctx context.Context, year int) ([]*salemongo.TransactionYearlyAmountSuccessRow, bool)
	SetCachedYearAmountSuccessCached(ctx context.Context, year int, res []*salemongo.TransactionYearlyAmountSuccessRow)
	GetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransaction) ([]*salemongo.TransactionMonthlyAmountFailedRow, bool)
	SetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransaction, res []*salemongo.TransactionMonthlyAmountFailedRow)
	GetCachedYearAmountFailedCached(ctx context.Context, year int) ([]*salemongo.TransactionYearlyAmountFailedRow, bool)
	SetCachedYearAmountFailedCached(ctx context.Context, year int, res []*salemongo.TransactionYearlyAmountFailedRow)
	GetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransaction) ([]*salemongo.TransactionMonthlyMethodSuccessRow, bool)
	SetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransaction, res []*salemongo.TransactionMonthlyMethodSuccessRow)
	GetCachedYearMethodSuccessCached(ctx context.Context, year int) ([]*salemongo.TransactionYearlyMethodSuccessRow, bool)
	SetCachedYearMethodSuccessCached(ctx context.Context, year int, res []*salemongo.TransactionYearlyMethodSuccessRow)
	GetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransaction) ([]*salemongo.TransactionMonthlyMethodFailedRow, bool)
	SetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransaction, res []*salemongo.TransactionMonthlyMethodFailedRow)
	GetCachedYearMethodFailedCached(ctx context.Context, year int) ([]*salemongo.TransactionYearlyMethodFailedRow, bool)
	SetCachedYearMethodFailedCached(ctx context.Context, year int, res []*salemongo.TransactionYearlyMethodFailedRow)
}

// TransactionStatsByMerchantCache caches per-merchant transaction statistics.
type TransactionStatsByMerchantCache interface {
	GetCachedMonthAmountSuccessByMerchantCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*salemongo.TransactionMonthlyAmountSuccessByMerchantRow, bool)
	SetCachedMonthAmountSuccessByMerchantCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*salemongo.TransactionMonthlyAmountSuccessByMerchantRow)
	GetCachedMonthAmountFailedByMerchantCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*salemongo.TransactionMonthlyAmountFailedByMerchantRow, bool)
	SetCachedMonthAmountFailedByMerchantCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*salemongo.TransactionMonthlyAmountFailedByMerchantRow)
	GetCachedYearAmountFailedByMerchantCached(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*salemongo.TransactionYearlyAmountFailedByMerchantRow, bool)
	SetCachedYearAmountFailedByMerchantCached(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*salemongo.TransactionYearlyAmountFailedByMerchantRow)
	GetCachedYearAmountSuccessByMerchantCached(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*salemongo.TransactionYearlyAmountSuccessByMerchantRow, bool)
	SetCachedYearAmountSuccessByMerchantCached(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*salemongo.TransactionYearlyAmountSuccessByMerchantRow)
	GetCachedMonthMethodSuccessByMerchantCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*salemongo.TransactionMonthlyMethodByMerchantSuccessRow, bool)
	SetCachedMonthMethodSuccessByMerchantCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*salemongo.TransactionMonthlyMethodByMerchantSuccessRow)
	GetCachedYearMethodSuccessByMerchantCached(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*salemongo.TransactionYearlyMethodByMerchantSuccessRow, bool)
	SetCachedYearMethodSuccessByMerchantCached(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*salemongo.TransactionYearlyMethodByMerchantSuccessRow)
	GetCachedMonthMethodFailedByMerchantCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*salemongo.TransactionMonthlyMethodByMerchantFailedRow, bool)
	SetCachedMonthMethodFailedByMerchantCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*salemongo.TransactionMonthlyMethodByMerchantFailedRow)
	GetCachedYearMethodFailedByMerchantCached(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*salemongo.TransactionYearlyMethodByMerchantFailedRow, bool)
	SetCachedYearMethodFailedByMerchantCached(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*salemongo.TransactionYearlyMethodByMerchantFailedRow)
}
