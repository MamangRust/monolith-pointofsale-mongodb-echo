package repository

import (
	"context"
	"time"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

// OrderResult is the result type for paginated order queries.
type OrderResult struct {
	OrderID    int32
	MerchantID int32
	CashierID  int32
	TotalPrice int64
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	TotalCount int64
}

// OrderResultDeleteAt is the result type for paginated order queries with deleted_at.
type OrderResultDeleteAt struct {
	OrderID    int32
	MerchantID int32
	CashierID  int32
	TotalPrice int64
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	DeletedAt  *time.Time
	TotalCount int64
}

type OrderQueryRepository interface {
	FindAllOrders(ctx context.Context, req *requests.FindAllOrders) ([]*OrderResult, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllOrders) ([]*OrderResultDeleteAt, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllOrders) ([]*OrderResultDeleteAt, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllOrderMerchant) ([]*OrderResult, *int, error)
	FindById(ctx context.Context, orderID int) (*salemongo.Order, error)
	FindByTrashedId(ctx context.Context, orderID int) (*salemongo.Order, error)
}

type OrderCommandRepository interface {
	DeleteOrder(ctx context.Context, orderID int) error
	CreateOrder(ctx context.Context, request *requests.CreateOrderRecordRequest) (*salemongo.Order, error)
	UpdateOrder(ctx context.Context, request *requests.UpdateOrderRecordRequest) (*salemongo.Order, error)
	FindAllTrashed(ctx context.Context) ([]*salemongo.Order, error)
	TrashedOrder(ctx context.Context, orderID int) (*salemongo.Order, error)
	RestoreOrder(ctx context.Context, orderID int) (*salemongo.Order, error)
	DeleteOrderPermanent(ctx context.Context, orderID int) (bool, error)
	DeleteAllOrderPermanent(ctx context.Context) (bool, error)
}

// The cashier/merchant/product/order_item lookups and writes order depends on
// live in other bounded contexts and are reached through the shared gRPC
// adapters in pkg/adapter/{cashier,merchant,product,orderitem} (see
// repository.go).

// OrderStatsRepository provides aggregate order statistics.
type OrderStatsRepository interface {
	GetMonthlyTotalRevenue(ctx context.Context, req *requests.MonthTotalRevenue) ([]*salemongo.MonthlyTotalRevenueRow, error)
	GetYearlyTotalRevenue(ctx context.Context, year int) ([]*salemongo.YearlyTotalRevenueRow, error)
	GetMonthlyOrder(ctx context.Context, year int) ([]*salemongo.MonthlyOrderRow, error)
	GetYearlyOrder(ctx context.Context, year int) ([]*salemongo.YearlyOrderRow, error)
}

// OrderStatByMerchantRepository provides order statistics grouped by merchant.
type OrderStatByMerchantRepository interface {
	GetMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*salemongo.MonthlyTotalRevenueByMerchantRow, error)
	GetYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*salemongo.YearlyTotalRevenueByMerchantRow, error)
	GetMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*salemongo.MonthlyOrderByMerchantRow, error)
	GetYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*salemongo.YearlyOrderByMerchantRow, error)
}

// sqlc row type aliases for handler compatibility
type OrderRow = OrderResult
type OrderActiveRow = OrderResultDeleteAt
type OrderTrashedRow = OrderResultDeleteAt
