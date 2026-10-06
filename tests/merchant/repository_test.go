package merchant_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-point-of-sale-merchant/repository"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	tests "github.com/MamangRust/monolith-point-of-sale-test"

	"github.com/stretchr/testify/suite"
)

type MerchantRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo       *repository.Repositories
	merchantID int
	userID     int
}

func (s *MerchantRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupUserService()

	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	userClient := pbusers.NewUserQueryServiceClient(s.Conns["user"])
	s.repo = repository.NewRepositories(mongoDB, userClient)

	userID, err := s.Fixtures().SeedUserRow(s.Ctx, "Merchant", "RepoTest", "merchant.repo@example.com", "password123")
	s.Require().NoError(err)
	s.userID = int(userID)
}

func (s *MerchantRepositoryTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *MerchantRepositoryTestSuite) Test1_CreateMerchant() {
	ctx := context.Background()

	req := &requests.CreateMerchantRequest{
		UserID:       s.userID,
		Name:         "Test Merchant",
		Description:  "Detailed description",
		Address:      "Merchant Street No. 1",
		ContactEmail: "merchant@example.com",
		ContactPhone: "08123456789",
		Status:       "active",
	}

	res, err := s.repo.MerchantCommand.CreateMerchant(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(req.Name, res.Name)
	s.Equal(int32(s.userID), res.UserID)
	s.merchantID = int(res.MerchantID)
}

func (s *MerchantRepositoryTestSuite) Test2_FindById() {
	s.Require().NotZero(s.merchantID)
	ctx := context.Background()

	found, err := s.repo.MerchantQuery.FindById(ctx, s.merchantID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(s.merchantID, int(found.MerchantID))
}

func (s *MerchantRepositoryTestSuite) Test3_UpdateMerchant() {
	s.Require().NotZero(s.merchantID)
	ctx := context.Background()

	req := &requests.UpdateMerchantRequest{
		MerchantID:   &s.merchantID,
		UserID:       s.userID,
		Name:         "Updated Merchant",
		Description:  "Updated description",
		Address:      "Updated Street",
		ContactEmail: "updated@example.com",
		ContactPhone: "08987654321",
		Status:       "active",
	}

	res, err := s.repo.MerchantCommand.UpdateMerchant(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal("Updated Merchant", res.Name)
}

func (s *MerchantRepositoryTestSuite) Test4_TrashAndRestore() {
	s.Require().NotZero(s.merchantID)
	ctx := context.Background()

	trashed, err := s.repo.MerchantCommand.TrashedMerchant(ctx, s.merchantID)
	s.NoError(err)
	s.NotNil(trashed)

	restored, err := s.repo.MerchantCommand.RestoreMerchant(ctx, s.merchantID)
	s.NoError(err)
	s.NotNil(restored)

	found, err := s.repo.MerchantQuery.FindById(ctx, s.merchantID)
	s.NoError(err)
	s.NotNil(found)
}

func (s *MerchantRepositoryTestSuite) Test5_DeletePermanent() {
	s.Require().NotZero(s.merchantID)
	ctx := context.Background()

	_, err := s.repo.MerchantCommand.TrashedMerchant(ctx, s.merchantID)
	s.NoError(err)

	success, err := s.repo.MerchantCommand.DeleteMerchantPermanent(ctx, s.merchantID)
	s.NoError(err)
	s.True(success)
}

func TestMerchantRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantRepositoryTestSuite))
}
