package transaction_test

import (
	"context"
	"testing"
	"time"

	pbtransactions "github.com/MamangRust/monolith-point-of-sale-pb/transactions"

	tests "github.com/MamangRust/monolith-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type TransactionStatsGapiTestSuite struct {
	tests.BaseTestSuite
	clientStats pbtransactions.TransactionStatsServiceClient
	userID      int
	merchantID  int
}

func (s *TransactionStatsGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupTransactionService()
	s.SetupOrderService()

	s.clientStats = pbtransactions.NewTransactionStatsServiceClient(s.Conns["transaction"])

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, catID)
	orderID := s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Create a successful transaction
	s.Require().NoError(s.Fixtures().SeedTransactionRow(ctx,
		int32(orderID), int32(s.merchantID), 0, 10000, "credit_card", "completed", time.Now().UTC()))

	// Create a failed transaction
	s.Require().NoError(s.Fixtures().SeedTransactionRow(ctx,
		int32(orderID), int32(s.merchantID), 0, 5000, "bank_transfer", "failed", time.Now().UTC()))
}

func (s *TransactionStatsGapiTestSuite) TestFindMonthStatusSuccess() {
	ctx := context.Background()
	now := time.Now()
	req := &pbtransactions.FindMonthlyTransactionStatus{
		Year:  int32(now.Year()),
		Month: int32(now.Month()),
	}

	res, err := s.clientStats.FindMonthStatusSuccess(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
}

func (s *TransactionStatsGapiTestSuite) TestFindYearStatusSuccess() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbtransactions.FindYearlyTransactionStatus{
		Year: int32(year),
	}

	res, err := s.clientStats.FindYearStatusSuccess(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
}

func (s *TransactionStatsGapiTestSuite) TestFindMonthStatusFailed() {
	ctx := context.Background()
	now := time.Now()
	req := &pbtransactions.FindMonthlyTransactionStatus{
		Year:  int32(now.Year()),
		Month: int32(now.Month()),
	}

	res, err := s.clientStats.FindMonthStatusFailed(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
}

func (s *TransactionStatsGapiTestSuite) TestFindYearStatusFailed() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbtransactions.FindYearlyTransactionStatus{
		Year: int32(year),
	}

	res, err := s.clientStats.FindYearStatusFailed(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
}

func TestTransactionStatsGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionStatsGapiTestSuite))
}
