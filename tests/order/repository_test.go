package order_test

import (
	"context"
	"net"
	"testing"

	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"

	"github.com/MamangRust/monolith-point-of-sale-order/repository"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	tests "github.com/MamangRust/monolith-point-of-sale-test"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	ts         *tests.TestSuite
	repo       *repository.Repositories
	orderID    int
	merchantID int
	cashierID  int
}

func (s *OrderRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	mongoDB, err := s.ts.MongoDB()
	s.Require().NoError(err)

	// cashier/merchant/product/order-item are reached through gRPC adapters,
	// but this suite only exercises the order-owned collections (the adapters
	// are called by the service layer, not the repositories), so the clients
	// are never dialled.
	cashierLis, err := net.Listen("tcp", "localhost:0")
	s.Require().NoError(err)
	merchantLis, err := net.Listen("tcp", "localhost:0")
	s.Require().NoError(err)
	productLis, err := net.Listen("tcp", "localhost:0")
	s.Require().NoError(err)
	orderItemLis, err := net.Listen("tcp", "localhost:0")
	s.Require().NoError(err)

	cashierConn, err := grpc.NewClient(cashierLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	merchantConn, err := grpc.NewClient(merchantLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	productConn, err := grpc.NewClient(productLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	orderItemConn, err := grpc.NewClient(orderItemLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.repo = repository.NewRepositories(
		mongoDB,
		pbcashiers.NewCashierQueryServiceClient(cashierConn),
		pbmerchants.NewMerchantQueryServiceClient(merchantConn),
		pbproducts.NewProductQueryServiceClient(productConn),
		pbproducts.NewProductCommandServiceClient(productConn),
		pborder_items.NewOrderItemQueryServiceClient(orderItemConn),
		pborder_items.NewOrderItemCommandServiceClient(orderItemConn),
	)

	userID, err := s.ts.SeedUserRow(s.ts.Ctx, "Order", "Repo", "order.repo@example.com", "password123")
	s.Require().NoError(err)

	merchantID, err := s.ts.SeedMerchantRow(s.ts.Ctx, userID, "Order Test Merchant", "Desc", "Addr", "o@example.com", "123", "active")
	s.Require().NoError(err)
	s.merchantID = int(merchantID)

	cashierID, err := s.ts.SeedCashierRow(s.ts.Ctx, merchantID, userID, "Order Test Cashier")
	s.Require().NoError(err)
	s.cashierID = int(cashierID)

	orderID, err := s.ts.SeedOrderRow(s.ts.Ctx, merchantID, cashierID, 5000)
	s.Require().NoError(err)
	s.orderID = int(orderID)
}

func (s *OrderRepositoryTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *OrderRepositoryTestSuite) Test1_FindById() {
	s.Require().NotZero(s.orderID)
	ctx := context.Background()

	found, err := s.repo.OrderQuery.FindById(ctx, s.orderID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(int32(s.orderID), found.OrderID)
}

func (s *OrderRepositoryTestSuite) Test2_CreateOrder() {
	ctx := context.Background()

	req := &requests.CreateOrderRecordRequest{
		MerchantID: s.merchantID,
		CashierID:  s.cashierID,
		TotalPrice: 10000,
	}

	res, err := s.repo.OrderCommand.CreateOrder(ctx, req)
	s.NoError(err)
	s.NotNil(res)
}

func (s *OrderRepositoryTestSuite) Test3_UpdateOrder() {
	s.Require().NotZero(s.orderID)
	ctx := context.Background()

	req := &requests.UpdateOrderRecordRequest{
		OrderID:    s.orderID,
		TotalPrice: 15000,
	}

	res, err := s.repo.OrderCommand.UpdateOrder(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(int64(15000), res.TotalPrice)
}

func (s *OrderRepositoryTestSuite) Test4_TrashAndRestore() {
	s.Require().NotZero(s.orderID)
	ctx := context.Background()

	trashed, err := s.repo.OrderCommand.TrashedOrder(ctx, s.orderID)
	s.NoError(err)
	s.NotNil(trashed)

	restored, err := s.repo.OrderCommand.RestoreOrder(ctx, s.orderID)
	s.NoError(err)
	s.NotNil(restored)

	found, err := s.repo.OrderQuery.FindById(ctx, s.orderID)
	s.NoError(err)
	s.NotNil(found)
}

func (s *OrderRepositoryTestSuite) Test5_DeletePermanent() {
	s.Require().NotZero(s.orderID)
	ctx := context.Background()

	_, err := s.repo.OrderCommand.TrashedOrder(ctx, s.orderID)
	s.NoError(err)

	success, err := s.repo.OrderCommand.DeleteOrderPermanent(ctx, s.orderID)
	s.NoError(err)
	s.True(success)
}

func TestOrderRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderRepositoryTestSuite))
}
