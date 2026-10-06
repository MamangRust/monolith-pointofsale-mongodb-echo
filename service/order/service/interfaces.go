package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-order/repository"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type OrderQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllOrders) ([]*repository.OrderResult, *int, error)
	FindById(ctx context.Context, orderID int) (*salemongo.Order, error)
	FindByActive(ctx context.Context, req *requests.FindAllOrders) ([]*repository.OrderResultDeleteAt, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllOrders) ([]*repository.OrderResultDeleteAt, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllOrderMerchant) ([]*repository.OrderResult, *int, error)
}

type OrderCommandService interface {
	CreateOrder(ctx context.Context, req *requests.CreateOrderRequest) (*salemongo.Order, error)
	UpdateOrder(ctx context.Context, req *requests.UpdateOrderRequest) (*salemongo.Order, error)
	TrashedOrder(ctx context.Context, orderID int) (*salemongo.Order, error)
	RestoreOrder(ctx context.Context, orderID int) (*salemongo.Order, error)
	DeleteOrderPermanent(ctx context.Context, orderID int) (bool, error)
	RestoreAllOrder(ctx context.Context) (bool, error)
	DeleteAllOrderPermanent(ctx context.Context) (bool, error)
}

type OrderStatsService interface {
	FindMonthlyTotalRevenue(ctx context.Context, req *requests.MonthTotalRevenue) ([]*salemongo.MonthlyTotalRevenueRow, error)
	FindYearlyTotalRevenue(ctx context.Context, year int) ([]*salemongo.YearlyTotalRevenueRow, error)
	FindMonthlyOrder(ctx context.Context, year int) ([]*salemongo.MonthlyOrderRow, error)
	FindYearlyOrder(ctx context.Context, year int) ([]*salemongo.YearlyOrderRow, error)
}

type OrderStatByMerchantQueryService interface {
	FindMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*salemongo.MonthlyTotalRevenueByMerchantRow, error)
	FindYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*salemongo.YearlyTotalRevenueByMerchantRow, error)
	FindMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*salemongo.MonthlyOrderByMerchantRow, error)
	FindYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*salemongo.YearlyOrderByMerchantRow, error)
}
