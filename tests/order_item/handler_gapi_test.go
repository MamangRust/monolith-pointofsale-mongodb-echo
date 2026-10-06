package order_item_test

import (
	"context"
	"testing"

	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"

	item_cache "github.com/MamangRust/monolith-point-of-sale-order-item/cache"
	item_handler "github.com/MamangRust/monolith-point-of-sale-order-item/handler"
	item_repo "github.com/MamangRust/monolith-point-of-sale-order-item/repository"
	item_service "github.com/MamangRust/monolith-point-of-sale-order-item/service"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	tests "github.com/MamangRust/monolith-point-of-sale-test"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type OrderItemGapiTestSuite struct {
	tests.BaseTestSuite
	client pborder_items.OrderItemQueryServiceClient
}

func (s *OrderItemGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupOrderService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	// Item dependencies
	mencache := item_cache.NewMencache(cacheStore)
	repos := item_repo.NewRepositories(mongoDB)
	svc := item_service.NewService(&item_service.Deps{
		Mencache:      mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := item_handler.NewHandler(&item_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pborder_items.RegisterOrderItemQueryServiceServer(server, handler.OrderItem)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.client = pborder_items.NewOrderItemQueryServiceClient(conn)
}

func (s *OrderItemGapiTestSuite) TestOrderItemGapiLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchantID := s.SeedMerchant(ctx, userID)
	productID := s.SeedProduct(ctx, merchantID, catID)
	orderID := s.SeedOrder(ctx, userID, merchantID, productID)

	// 2. Create an order item directly in DB for query testing
	orderItemID, err := s.Fixtures().SeedOrderItemRow(ctx, int32(orderID), int32(productID), 10, 700)
	s.Require().NoError(err)

	// 3. FindOrderItemByOrder
	findByOrderRes, err := s.client.FindOrderItemByOrder(ctx, &pborder_items.FindByIdOrderItemRequest{Id: int32(orderID)})
	s.Require().NoError(err)
	s.NotEmpty(findByOrderRes.Data)

	// 4. FindAll
	allRes, err := s.client.FindAll(ctx, &pborder_items.FindAllOrderItemRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. FindByActive
	activeRes, err := s.client.FindByActive(ctx, &pborder_items.FindAllOrderItemRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(activeRes.Data)

	// 6. Trash the order item directly in DB
	s.Require().NoError(s.Fixtures().TrashOrderItemRow(ctx, orderItemID))

	// 7. FindByTrashed
	trashedRes, err := s.client.FindByTrashed(ctx, &pborder_items.FindAllOrderItemRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 8. RestoreAll and DeleteAll (test these exist on the combined client)
	// These are empty operations on query-only interface
	_, err = s.client.FindAll(ctx, &pborder_items.FindAllOrderItemRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
}

func TestOrderItemGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderItemGapiTestSuite))
}
