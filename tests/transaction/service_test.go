package transaction_test

import (
	"context"
	"testing"

	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	tests "github.com/MamangRust/monolith-point-of-sale-test"
	trans_cache "github.com/MamangRust/monolith-point-of-sale-transacton/cache"
	"github.com/MamangRust/monolith-point-of-sale-transacton/repository"
	"github.com/MamangRust/monolith-point-of-sale-transacton/service"
	"github.com/stretchr/testify/suite"
)

type TransactionServiceTestSuite struct {
	tests.BaseTestSuite
	svc *service.Service
}

func (s *TransactionServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// gRPC service nyata: cashier(→user/merchant), merchant, order(→product/order-item), category.
	s.SetupTransactionService()
	s.SetupCategoryService()

	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	mencache := trans_cache.NewMencache(s.GetCacheStore())
	cashierClient := pbcashiers.NewCashierQueryServiceClient(s.Conns["cashier"])
	merchantClient := pbmerchants.NewMerchantQueryServiceClient(s.Conns["merchant"])
	orderClient := pborders.NewOrderQueryServiceClient(s.Conns["order"])
	orderItemClient := pborder_items.NewOrderItemQueryServiceClient(s.Conns["order-item"])
	repos := repository.NewRepositories(mongoDB, cashierClient, merchantClient, orderClient, orderItemClient)

	s.svc = service.NewService(&service.Deps{
		Kafka:         nil,
		Mencache:      mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
}

func (s *TransactionServiceTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *TransactionServiceTestSuite) TestTransactionLifecycle() {
	ctx := context.Background()

	userID := s.SeedUser(ctx)
	merchantID := s.SeedMerchant(ctx, userID)
	categoryID := s.SeedCategory(ctx)
	productID := s.SeedProduct(ctx, merchantID, categoryID)
	orderID := s.SeedOrder(ctx, userID, merchantID, productID)
	s.SeedOrderItem(ctx, orderID, productID)

	// cashier_id as seeded by SeedOrder (cashiers table)
	cashierIDRaw, err := s.Fixtures().CashierIDForUser(ctx, int32(userID), int32(merchantID))
	s.Require().NoError(err)
	cashierID := int(cashierIDRaw)

	// 2. Create Transaction
	req := &requests.CreateTransactionRequest{
		OrderID:       orderID,
		CashierID:     cashierID,
		MerchantID:    merchantID,
		PaymentMethod: "Transfer Bank",
		Amount:        1000000,
	}
	created, err := s.svc.TransactionCommand.CreateTransaction(ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(created)
	transactionID := int(created.TransactionID)

	// 3. FindByID
	found, err := s.svc.TransactionQuery.FindById(ctx, transactionID)
	s.Require().NoError(err)
	s.Equal("success", found.PaymentStatus)

	// 4. Update
	newPaymentMethod := "GOPAY"
	updateReq := &requests.UpdateTransactionRequest{
		TransactionID: &transactionID,
		OrderID:       orderID,
		CashierID:     cashierID,
		PaymentMethod: newPaymentMethod,
		Amount:        req.Amount,
	}
	updated, err := s.svc.TransactionCommand.UpdateTransaction(ctx, updateReq)
	s.Require().NoError(err)
	s.Equal(newPaymentMethod, updated.PaymentMethod)

	// 5. FindAll
	_, total, err := s.svc.TransactionQuery.FindAllTransactions(ctx, &requests.FindAllTransaction{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*total, 1)

	// 6. Trash
	_, err = s.svc.TransactionCommand.TrashedTransaction(ctx, transactionID)
	s.Require().NoError(err)

	// 7. FindTrashed
	_, totalTrashed, err := s.svc.TransactionQuery.FindByTrashed(ctx, &requests.FindAllTransaction{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.GreaterOrEqual(*totalTrashed, 1)

	// 8. FindActive
	active, _, err := s.svc.TransactionQuery.FindByActive(ctx, &requests.FindAllTransaction{Search: "", Page: 1, PageSize: 10})
	s.Require().NoError(err)
	for _, tx := range active {
		s.NotEqual(transactionID, int(tx.TransactionID))
	}

	// 9. Restore
	_, err = s.svc.TransactionCommand.RestoreTransaction(ctx, transactionID)
	s.Require().NoError(err)

	// 10. DeletePermanent
	_, err = s.svc.TransactionCommand.TrashedTransaction(ctx, transactionID)
	s.Require().NoError(err)
	success, err := s.svc.TransactionCommand.DeleteTransactionPermanently(ctx, transactionID)
	s.Require().NoError(err)
	s.True(success)

	// 11. RestoreAll & DeleteAll — create transactions and trash them
	txOrder1 := s.SeedOrder(ctx, userID, merchantID, productID)
	s.SeedOrderItem(ctx, txOrder1, productID)
	tx1, _ := s.svc.TransactionCommand.CreateTransaction(ctx, &requests.CreateTransactionRequest{OrderID: txOrder1, CashierID: cashierID, MerchantID: merchantID, PaymentMethod: "Cash", Amount: 100000})

	txOrder2 := s.SeedOrder(ctx, userID, merchantID, productID)
	s.SeedOrderItem(ctx, txOrder2, productID)
	tx2, _ := s.svc.TransactionCommand.CreateTransaction(ctx, &requests.CreateTransactionRequest{OrderID: txOrder2, CashierID: cashierID, MerchantID: merchantID, PaymentMethod: "Cash", Amount: 100000})

	s.svc.TransactionCommand.TrashedTransaction(ctx, int(tx1.TransactionID))
	s.svc.TransactionCommand.TrashedTransaction(ctx, int(tx2.TransactionID))

	resRestoreAll, err := s.svc.TransactionCommand.RestoreAllTransactions(ctx)
	s.Require().NoError(err)
	s.True(resRestoreAll)

	s.svc.TransactionCommand.TrashedTransaction(ctx, int(tx1.TransactionID))
	s.svc.TransactionCommand.TrashedTransaction(ctx, int(tx2.TransactionID))

	resDeleteAll, err := s.svc.TransactionCommand.DeleteAllTransactionPermanent(ctx)
	s.Require().NoError(err)
	s.True(resDeleteAll)
}

func TestTransactionServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionServiceTestSuite))
}
