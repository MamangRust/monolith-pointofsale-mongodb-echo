package product_test

import (
	"context"
	"net"
	"testing"

	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"

	"github.com/MamangRust/monolith-point-of-sale-product/repository"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	tests "github.com/MamangRust/monolith-point-of-sale-test"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductRepositoryTestSuite struct {
	suite.Suite
	ts        *tests.TestSuite
	repo      *repository.Repositories
	productID int
}

func (s *ProductRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	mongoDB, err := s.ts.MongoDB()
	s.Require().NoError(err)

	// category and merchant are reached through gRPC adapters, but this suite
	// only exercises the product-owned collections (the adapters are called by
	// the service layer, not the repositories), so the clients are never dialled.
	categoryLis, err := net.Listen("tcp", "localhost:0")
	s.Require().NoError(err)
	merchantLis, err := net.Listen("tcp", "localhost:0")
	s.Require().NoError(err)

	categoryConn, err := grpc.NewClient(categoryLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	merchantConn, err := grpc.NewClient(merchantLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.repo = repository.NewRepositories(
		mongoDB,
		pbcategories.NewCategoryQueryServiceClient(categoryConn),
		pbmerchants.NewMerchantQueryServiceClient(merchantConn),
	)

	// Seed a merchant and category for product tests
	userID, err := s.ts.SeedUserRow(s.ts.Ctx, "Prod", "Repo", "prod.repo@example.com", "password123")
	s.Require().NoError(err)

	_, err = s.ts.SeedCategoryRow(s.ts.Ctx, "Test Category", "Category for product tests")
	s.Require().NoError(err)

	_, err = s.ts.SeedMerchantRow(s.ts.Ctx, userID, "Test Merchant", "Desc", "Addr", "pm@example.com", "123", "active")
	s.Require().NoError(err)
}

func (s *ProductRepositoryTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *ProductRepositoryTestSuite) Test1_CreateProduct() {
	ctx := context.Background()

	req := &requests.CreateProductRequest{
		MerchantID:   1,
		CategoryID:   1,
		Name:         "Test Product",
		Description:  "Product description",
		Price:        100,
		CountInStock: 50,
		Brand:        "Test Brand",
		Weight:       1,
		ImageProduct: "test.jpg",
	}

	res, err := s.repo.ProductCommand.CreateProduct(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(req.Name, res.Name)
	s.productID = int(res.ProductID)
}

func (s *ProductRepositoryTestSuite) Test2_FindById() {
	s.Require().NotZero(s.productID)
	ctx := context.Background()

	found, err := s.repo.ProductQuery.FindById(ctx, s.productID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(s.productID, int(found.ProductID))
}

func (s *ProductRepositoryTestSuite) Test3_UpdateProduct() {
	s.Require().NotZero(s.productID)
	ctx := context.Background()

	req := &requests.UpdateProductRequest{
		ProductID:    &s.productID,
		MerchantID:   1,
		CategoryID:   1,
		Name:         "Updated Product",
		Description:  "Updated description",
		Price:        200,
		CountInStock: 100,
		Brand:        "Updated Brand",
		Weight:       2,
		ImageProduct: "updated.jpg",
	}

	res, err := s.repo.ProductCommand.UpdateProduct(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal("Updated Product", res.Name)
}

func (s *ProductRepositoryTestSuite) Test4_TrashAndRestore() {
	s.Require().NotZero(s.productID)
	ctx := context.Background()

	// Trash
	trashed, err := s.repo.ProductCommand.TrashedProduct(ctx, s.productID)
	s.NoError(err)
	s.NotNil(trashed)

	// Restore
	restored, err := s.repo.ProductCommand.RestoreProduct(ctx, s.productID)
	s.NoError(err)
	s.NotNil(restored)

	// Verify restored
	found, err := s.repo.ProductQuery.FindById(ctx, s.productID)
	s.NoError(err)
	s.NotNil(found)
}

func (s *ProductRepositoryTestSuite) Test5_DeletePermanent() {
	s.Require().NotZero(s.productID)
	ctx := context.Background()

	// Must be trashed first for permanent delete
	_, err := s.repo.ProductCommand.TrashedProduct(ctx, s.productID)
	s.NoError(err)

	success, err := s.repo.ProductCommand.DeleteProductPermanent(ctx, s.productID)
	s.NoError(err)
	s.True(success)
}

func TestProductRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ProductRepositoryTestSuite))
}
