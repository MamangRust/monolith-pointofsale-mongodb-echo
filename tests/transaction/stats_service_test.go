package transaction_test

import (
	"context"
	"testing"
	"time"

	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	tests "github.com/MamangRust/monolith-point-of-sale-test"
	transaction_cache "github.com/MamangRust/monolith-point-of-sale-transacton/cache"
	"github.com/MamangRust/monolith-point-of-sale-transacton/repository"
	"github.com/MamangRust/monolith-point-of-sale-transacton/service"
	"github.com/stretchr/testify/suite"
)

type TransactionStatsServiceTestSuite struct {
	tests.BaseTestSuite
	svc           service.TransactionStatsService
	svcByMerchant service.TransactionStatsByMerchantQueryService
	merchantID    int
	userID        int
}

func (s *TransactionStatsServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupTransactionService()
	s.SetupOrderService()

	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)
	cashierClient := pbcashiers.NewCashierQueryServiceClient(s.Conns["cashier"])
	merchantClient := pbmerchants.NewMerchantQueryServiceClient(s.Conns["merchant"])
	orderClient := pborders.NewOrderQueryServiceClient(s.Conns["order"])
	orderItemClient := pborder_items.NewOrderItemQueryServiceClient(s.Conns["order-item"])
	repos := repository.NewRepositories(mongoDB, cashierClient, merchantClient, orderClient, orderItemClient)
	cacheStore := s.GetCacheStore()
	mencache := transaction_cache.NewMencache(cacheStore)

	s.svc = service.NewTransactionStatsService(
		mencache,
		repos.TransactionStats,
		s.Log,
		s.Obs,
	)
	s.svcByMerchant = service.NewTransactionStatsByMerchantQueryService(
		mencache,
		repos.TransactionStatsByMerchant,
		s.Log,
		s.Obs,
	)

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, catID)
	orderID := s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Create a successful transaction for this order
	s.Require().NoError(s.Fixtures().SeedTransactionRow(ctx,
		int32(orderID), int32(s.merchantID), 0, 10000, "credit_card", "completed", time.Now().UTC()))

	// Create a failed transaction for this order
	s.Require().NoError(s.Fixtures().SeedTransactionRow(ctx,
		int32(orderID), int32(s.merchantID), 0, 10000, "credit_card", "failed", time.Now().UTC()))
}

func (s *TransactionStatsServiceTestSuite) TestFindMonthlyAmountSuccess() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.MonthAmountTransaction{
		Year:  now.Year(),
		Month: int(now.Month()),
	}

	res, err := s.svc.FindMonthlyAmountSuccess(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsServiceTestSuite) TestFindYearlyAmountSuccess() {
	ctx := context.Background()
	year := time.Now().Year()

	res, err := s.svc.FindYearlyAmountSuccess(ctx, year)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsServiceTestSuite) TestFindMonthlyAmountFailed() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.MonthAmountTransaction{
		Year:  now.Year(),
		Month: int(now.Month()),
	}

	res, err := s.svc.FindMonthlyAmountFailed(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsServiceTestSuite) TestFindYearlyAmountFailed() {
	ctx := context.Background()
	year := time.Now().Year()

	res, err := s.svc.FindYearlyAmountFailed(ctx, year)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsServiceTestSuite) TestFindMonthlyMethodSuccess() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.MonthMethodTransaction{
		Year:  now.Year(),
		Month: int(now.Month()),
	}

	res, err := s.svc.FindMonthlyMethodSuccess(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsServiceTestSuite) TestFindYearlyMethodSuccess() {
	ctx := context.Background()
	year := time.Now().Year()

	res, err := s.svc.FindYearlyMethodSuccess(ctx, year)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsServiceTestSuite) TestFindMonthlyAmountSuccessByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.MonthAmountTransactionMerchant{
		Year:       now.Year(),
		Month:      int(now.Month()),
		MerchantID: s.merchantID,
	}

	res, err := s.svcByMerchant.FindMonthlyAmountSuccessByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func (s *TransactionStatsServiceTestSuite) TestFindYearlyAmountSuccessByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &requests.YearAmountTransactionMerchant{
		Year:       now.Year(),
		MerchantID: s.merchantID,
	}

	res, err := s.svcByMerchant.FindYearlyAmountSuccessByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res)
}

func TestTransactionStatsServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionStatsServiceTestSuite))
}
