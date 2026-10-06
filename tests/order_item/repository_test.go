package order_item_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-point-of-sale-order-item/repository"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	tests "github.com/MamangRust/monolith-point-of-sale-test"

	"github.com/stretchr/testify/suite"
)

type OrderItemRepositoryTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	repo        *repository.Repositories
	orderItemID int
	orderID     int
	productID   int
}

func (s *OrderItemRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	mongoDB, err := s.ts.MongoDB()
	s.Require().NoError(err)

	s.repo = repository.NewRepositories(mongoDB)

	userID, err := s.ts.SeedUserRow(s.ts.Ctx, "OItem", "Repo", "oitem.repo@example.com", "password123")
	s.Require().NoError(err)

	merchantID, err := s.ts.SeedMerchantRow(s.ts.Ctx, userID, "OItem Merchant", "", "", "", "", "active")
	s.Require().NoError(err)

	categoryID, err := s.ts.SeedCategoryRow(s.ts.Ctx, "OItem Category", "")
	s.Require().NoError(err)

	productID, err := s.ts.SeedProductRow(s.ts.Ctx, merchantID, categoryID, "OItem Product", "", 1000, 10)
	s.Require().NoError(err)
	s.productID = int(productID)

	cashierID, err := s.ts.SeedCashierRow(s.ts.Ctx, merchantID, userID, "Test Cashier")
	s.Require().NoError(err)

	orderID, err := s.ts.SeedOrderRow(s.ts.Ctx, merchantID, cashierID, 5000)
	s.Require().NoError(err)
	s.orderID = int(orderID)

	orderItemID, err := s.ts.SeedOrderItemRow(s.ts.Ctx, orderID, productID, 2, 5000)
	s.Require().NoError(err)
	s.orderItemID = int(orderItemID)
}

func (s *OrderItemRepositoryTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *OrderItemRepositoryTestSuite) Test1_FindAll() {
	ctx := context.Background()

	results, total, err := s.repo.OrderItemQuery.FindAllOrderItems(ctx, &requests.FindAllOrderItems{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(total)
	s.GreaterOrEqual(*total, 1)
	s.GreaterOrEqual(len(results), 1)
}

func (s *OrderItemRepositoryTestSuite) Test2_FindOrderItemByOrder() {
	s.Require().NotZero(s.orderID)
	ctx := context.Background()

	found, err := s.repo.OrderItemQuery.FindOrderItemByOrder(ctx, s.orderID)
	s.NoError(err)
	s.NotEmpty(found)
	s.Equal(int32(s.orderID), found[0].OrderID)
}

func (s *OrderItemRepositoryTestSuite) Test3_FindByActive() {
	ctx := context.Background()

	results, total, err := s.repo.OrderItemQuery.FindByActive(ctx, &requests.FindAllOrderItems{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(total)
	s.GreaterOrEqual(*total, 1)
	s.GreaterOrEqual(len(results), 1)
}

func (s *OrderItemRepositoryTestSuite) Test4_FindByTrashed() {
	ctx := context.Background()

	s.Require().NoError(s.ts.TrashOrderItemRow(ctx, int32(s.orderItemID)))

	results, total, err := s.repo.OrderItemQuery.FindByTrashed(ctx, &requests.FindAllOrderItems{
		Search:   "",
		Page:     1,
		PageSize: 10,
	})
	s.NoError(err)
	s.NotNil(total)
	s.GreaterOrEqual(*total, 1)
	s.GreaterOrEqual(len(results), 1)

	s.Require().NoError(s.ts.RestoreOrderItemRow(ctx, int32(s.orderItemID)))
}

func TestOrderItemRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderItemRepositoryTestSuite))
}
