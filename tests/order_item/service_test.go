package order_item_test

import (
	"context"
	"testing"

	item_cache "github.com/MamangRust/monolith-point-of-sale-order-item/cache"
	"github.com/MamangRust/monolith-point-of-sale-order-item/repository"
	"github.com/MamangRust/monolith-point-of-sale-order-item/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	tests "github.com/MamangRust/monolith-point-of-sale-test"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type OrderItemServiceTestSuite struct {
	suite.Suite
	ts  *tests.TestSuite
	svc *service.Service
}

func (s *OrderItemServiceTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	mongoDB, err := s.ts.MongoDB()
	s.Require().NoError(err)

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	redisClient := redis.NewClient(opts)

	log, _ := logger.NewLogger("test", nil)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(redisClient, log, cacheMetrics)
	mencache := item_cache.NewMencache(cacheStore)

	obs, _ := observability.NewObservability("test", log)

	repos := repository.NewRepositories(mongoDB)

	s.svc = service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        log,
		Mencache:      mencache,
		Observability: obs,
	})
}

func (s *OrderItemServiceTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *OrderItemServiceTestSuite) TestOrderItemLifecycle() {
	ctx := context.Background()

	// Seed data directly (user → merchant → cashier → category → product → order)
	userID, err := s.ts.SeedUserRow(ctx, "OItem", "Svc", "oitem.svc@example.com", "password123")
	s.Require().NoError(err)

	merchantID, err := s.ts.SeedMerchantRow(ctx, userID, "OI Merchant", "Desc", "Addr", "oi@example.com", "123", "active")
	s.Require().NoError(err)

	cashierID, err := s.ts.SeedCashierRow(ctx, merchantID, userID, "OI Cashier")
	s.Require().NoError(err)

	categoryID, err := s.ts.SeedCategoryRow(ctx, "Category", "Desc")
	s.Require().NoError(err)

	productID, err := s.ts.SeedProductRow(ctx, merchantID, categoryID, "Product", "Desc", 1000, 10)
	s.Require().NoError(err)

	orderID, err := s.ts.SeedOrderRow(ctx, merchantID, cashierID, 5000)
	s.Require().NoError(err)

	// 1. Create OrderItem (direct DB since service is query-only and doesn't have command)
	// The order_item service is query-only (OrderItemQuery), items are created via order service
	// Test the query interface

	// Create an order item directly for query testing
	oItemID, err := s.ts.SeedOrderItemRow(ctx, orderID, productID, 2, 5000)
	s.Require().NoError(err)

	// 2. FindByOrder
	items, err := s.svc.OrderItemQuery.FindOrderItemByOrder(ctx, int(orderID))
	s.Require().NoError(err)
	s.NotEmpty(items)

	// 3. FindAll
	_, total, err := s.svc.OrderItemQuery.FindAllOrderItems(ctx, &requests.FindAllOrderItems{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 4. FindByActive
	active, _, err := s.svc.OrderItemQuery.FindByActive(ctx, &requests.FindAllOrderItems{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(active)

	// 5. Trash the order item directly
	s.Require().NoError(s.ts.TrashOrderItemRow(ctx, oItemID))

	// 6. FindByTrashed
	_, totalTrashed, err := s.svc.OrderItemQuery.FindByTrashed(ctx, &requests.FindAllOrderItems{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)
}

func TestOrderItemServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderItemServiceTestSuite))
}
