package order_test

import (
	"context"
	"testing"

	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"

	order_cache "github.com/MamangRust/monolith-point-of-sale-order/cache"
	"github.com/MamangRust/monolith-point-of-sale-order/repository"
	"github.com/MamangRust/monolith-point-of-sale-order/service"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	tests "github.com/MamangRust/monolith-point-of-sale-test"

	"github.com/stretchr/testify/suite"
)

type OrderServiceTestSuite struct {
	tests.BaseTestSuite
	svc       *service.Service
	orderID   int // merchant id
	cashierID int
	productID int
}

func (s *OrderServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	// gRPC service nyata: cashier(→user/merchant), merchant, product, order-item.
	s.SetupOrderService()

	// Seed user + merchant + cashier + category + product
	userID, err := s.Fixtures().SeedUserRow(s.Ctx, "Order", "Svc", "order.svc@example.com", "password123")
	s.Require().NoError(err)

	merchantID, err := s.Fixtures().SeedMerchantRow(s.Ctx, userID, "Order Svc Merchant", "Desc", "Addr", "os@example.com", "123", "active")
	s.Require().NoError(err)
	s.orderID = int(merchantID)

	cashierID, err := s.Fixtures().SeedCashierRow(s.Ctx, merchantID, userID, "Order Cashier")
	s.Require().NoError(err)
	s.cashierID = int(cashierID)

	categoryID, err := s.Fixtures().SeedCategoryRow(s.Ctx, "Order Category", "Desc")
	s.Require().NoError(err)

	productID, err := s.Fixtures().SeedProductRow(s.Ctx, merchantID, categoryID, "Order Product", "Desc", 10000, 100)
	s.Require().NoError(err)
	s.productID = int(productID)

	mencache := order_cache.NewMencache(s.GetCacheStore())
	repos := repository.NewRepositories(
		mongoDB,
		pbcashiers.NewCashierQueryServiceClient(s.Conns["cashier"]),
		pbmerchants.NewMerchantQueryServiceClient(s.Conns["merchant"]),
		pbproducts.NewProductQueryServiceClient(s.Conns["product"]),
		pbproducts.NewProductCommandServiceClient(s.Conns["product"]),
		pborder_items.NewOrderItemQueryServiceClient(s.Conns["order-item"]),
		pborder_items.NewOrderItemCommandServiceClient(s.Conns["order-item"]),
	)

	s.svc = service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        s.Log,
		Mencache:      mencache,
		Observability: s.Obs,
	})
}

func (s *OrderServiceTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *OrderServiceTestSuite) TestOrderLifecycle() {
	ctx := context.Background()

	// 1. Create
	req := &requests.CreateOrderRequest{
		MerchantID: s.orderID,
		CashierID:  s.cashierID,
		Items: []requests.CreateOrderItemRequest{
			{ProductID: s.productID, Quantity: 1},
		},
	}
	created, err := s.svc.OrderCommand.CreateOrder(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	orderID := int(created.OrderID)

	orderItemIDRaw, err := s.Fixtures().OrderItemIDForOrder(s.Ctx, int32(orderID))
	s.Require().NoError(err)
	orderItemID := int(orderItemIDRaw)

	// 2. FindByID
	found, err := s.svc.OrderQuery.FindById(ctx, orderID)
	s.Require().NoError(err)
	s.NotNil(found)

	// 3. Update
	updateReq := &requests.UpdateOrderRequest{
		OrderID: &orderID,
		Items: []requests.UpdateOrderItemRequest{
			{OrderItemID: orderItemID, ProductID: s.productID, Quantity: 2},
		},
	}
	updated, err := s.svc.OrderCommand.UpdateOrder(ctx, updateReq)
	s.Require().NoError(err)
	s.NotNil(updated)

	// 4. FindAll
	_, total, err := s.svc.OrderQuery.FindAll(ctx, &requests.FindAllOrders{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 5. Trash
	_, err = s.svc.OrderCommand.TrashedOrder(ctx, orderID)
	s.Require().NoError(err)

	// 6. FindTrashed
	_, totalTrashed, err := s.svc.OrderQuery.FindByTrashed(ctx, &requests.FindAllOrders{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 7. FindActive
	active, _, err := s.svc.OrderQuery.FindByActive(ctx, &requests.FindAllOrders{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	for _, o := range active {
		s.NotEqual(orderID, int(o.OrderID))
	}

	// 8. Restore
	_, err = s.svc.OrderCommand.RestoreOrder(ctx, orderID)
	s.Require().NoError(err)

	// 9. DeletePermanent
	_, err = s.svc.OrderCommand.TrashedOrder(ctx, orderID)
	s.Require().NoError(err)
	success, err := s.svc.OrderCommand.DeleteOrderPermanent(ctx, orderID)
	s.Require().NoError(err)
	s.True(success)

	// 10. RestoreAll & DeleteAll
	o1, _ := s.svc.OrderCommand.CreateOrder(ctx, &requests.CreateOrderRequest{
		MerchantID: s.orderID,
		CashierID:  s.cashierID,
		Items: []requests.CreateOrderItemRequest{
			{ProductID: s.productID, Quantity: 1},
		},
	})
	o2, _ := s.svc.OrderCommand.CreateOrder(ctx, &requests.CreateOrderRequest{
		MerchantID: s.orderID,
		CashierID:  s.cashierID,
		Items: []requests.CreateOrderItemRequest{
			{ProductID: s.productID, Quantity: 2},
		},
	})

	s.svc.OrderCommand.TrashedOrder(ctx, int(o1.OrderID))
	s.svc.OrderCommand.TrashedOrder(ctx, int(o2.OrderID))

	resRestoreAll, err := s.svc.OrderCommand.RestoreAllOrder(ctx)
	s.Require().NoError(err)
	s.True(resRestoreAll)

	s.svc.OrderCommand.TrashedOrder(ctx, int(o1.OrderID))
	s.svc.OrderCommand.TrashedOrder(ctx, int(o2.OrderID))

	resDeleteAll, err := s.svc.OrderCommand.DeleteAllOrderPermanent(ctx)
	s.Require().NoError(err)
	s.True(resDeleteAll)
}

func TestOrderServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderServiceTestSuite))
}
