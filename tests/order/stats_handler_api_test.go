package order_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apigw_cache "github.com/MamangRust/monolith-point-of-sale-apigateway/cache/gateway_cache"
	apigateway "github.com/MamangRust/monolith-point-of-sale-apigateway/handler"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	apierror "github.com/MamangRust/monolith-point-of-sale-shared/errors"
	response_api "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/api"
	tests "github.com/MamangRust/monolith-point-of-sale-test"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type OrderStatsApiTestSuite struct {
	tests.BaseTestSuite
	router     *echo.Echo
	merchantID int
	userID     int
}

func (s *OrderStatsApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupTransactionService()
	s.SetupOrderService()

	s.router = echo.New()

	orderConn := s.Conns["order"]

	apigateway.NewHandlerOrder(
		s.router,
		pborders.NewOrderQueryServiceClient(orderConn),
		pborders.NewOrderCommandServiceClient(orderConn),
		pborders.NewOrderStatsServiceClient(orderConn),
		pborders.NewOrderStatsByIdServiceClient(orderConn),
		pborders.NewOrderStatsByMerchantServiceClient(orderConn),
		s.Log,
		response_api.NewOrderResponseMapper(),
		apierror.NewApiHandler(s.Obs, s.Log),
		apigw_cache.NewGatewayCache(s.GetCacheStore()),
	)

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, catID)
	orderID := s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Ensure created_at is set to current time to be picked up by stats
	s.Require().NoError(s.Fixtures().SetOrderCreatedAt(ctx, int32(orderID), time.Now().UTC()))
}

func (s *OrderStatsApiTestSuite) TestFindMonthlyTotalRevenue() {
	now := time.Now()
	url := fmt.Sprintf("/api/order/monthly-total-revenue?year=%d&month=%d", now.Year(), int(now.Month()))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	s.Equal("success", res["status"])
	s.NotEmpty(res["data"])
}

func (s *OrderStatsApiTestSuite) TestFindYearlyTotalRevenue() {
	year := time.Now().Year()
	url := fmt.Sprintf("/api/order/yearly-total-revenue?year=%d", year)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	s.Equal("success", res["status"])
	s.NotEmpty(res["data"])
}

func (s *OrderStatsApiTestSuite) TestFindMonthlyOrder() {
	year := time.Now().Year()
	url := fmt.Sprintf("/api/order/monthly-revenue?year=%d", year)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	s.Equal("success", res["status"])
	s.NotEmpty(res["data"])
}

func (s *OrderStatsApiTestSuite) TestFindMonthlyTotalRevenueByMerchant() {
	now := time.Now()
	url := fmt.Sprintf("/api/order/merchant/monthly-total-revenue?year=%d&month=%d&merchant_id=%d",
		now.Year(), int(now.Month()), s.merchantID)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	s.Equal("success", res["status"])
	s.NotEmpty(res["data"])
}

func TestOrderStatsApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderStatsApiTestSuite))
}
