package mencache

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-order/repository"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type OrderQueryCache interface {
	GetOrderAllCache(ctx context.Context, req *requests.FindAllOrders) ([]*repository.OrderResult, *int, bool)
	SetOrderAllCache(ctx context.Context, req *requests.FindAllOrders, data []*repository.OrderResult, total *int)

	GetCachedOrderCache(ctx context.Context, orderID int) (*salemongo.Order, bool)
	SetCachedOrderCache(ctx context.Context, data *salemongo.Order)

	GetCachedOrderMerchant(ctx context.Context, req *requests.FindAllOrderMerchant) ([]*repository.OrderResult, *int, bool)
	SetCachedOrderMerchant(ctx context.Context, req *requests.FindAllOrderMerchant, res []*repository.OrderResult, total *int)

	GetOrderActiveCache(ctx context.Context, req *requests.FindAllOrders) ([]*repository.OrderResultDeleteAt, *int, bool)
	SetOrderActiveCache(ctx context.Context, req *requests.FindAllOrders, data []*repository.OrderResultDeleteAt, total *int)

	GetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrders) ([]*repository.OrderResultDeleteAt, *int, bool)
	SetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrders, data []*repository.OrderResultDeleteAt, total *int)
}

type OrderCommandCache interface {
	DeleteOrderCache(ctx context.Context, id int)
	DeleteOrderAllCache(ctx context.Context)
}

type OrderStatsCache interface {
	GetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue) ([]*salemongo.MonthlyTotalRevenueRow, bool)
	SetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue, res []*salemongo.MonthlyTotalRevenueRow)

	GetYearlyTotalRevenueCache(ctx context.Context, year int) ([]*salemongo.YearlyTotalRevenueRow, bool)
	SetYearlyTotalRevenueCache(ctx context.Context, year int, res []*salemongo.YearlyTotalRevenueRow)

	GetMonthlyOrderCache(ctx context.Context, year int) ([]*salemongo.MonthlyOrderRow, bool)
	SetMonthlyOrderCache(ctx context.Context, year int, res []*salemongo.MonthlyOrderRow)

	GetYearlyOrderCache(ctx context.Context, year int) ([]*salemongo.YearlyOrderRow, bool)
	SetYearlyOrderCache(ctx context.Context, year int, res []*salemongo.YearlyOrderRow)
}

type OrderStatsByMerchantCache interface {
	GetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*salemongo.MonthlyTotalRevenueByMerchantRow, bool)
	SetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant, res []*salemongo.MonthlyTotalRevenueByMerchantRow)

	GetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*salemongo.YearlyTotalRevenueByMerchantRow, bool)
	SetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant, res []*salemongo.YearlyTotalRevenueByMerchantRow)

	GetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant) ([]*salemongo.MonthlyOrderByMerchantRow, bool)
	SetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant, res []*salemongo.MonthlyOrderByMerchantRow)

	GetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant) ([]*salemongo.YearlyOrderByMerchantRow, bool)
	SetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant, res []*salemongo.YearlyOrderByMerchantRow)
}
