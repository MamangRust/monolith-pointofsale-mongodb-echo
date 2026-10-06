package merchant_test

import (
	"context"
	"testing"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-merchant/repository"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	tests "github.com/MamangRust/monolith-point-of-sale-test"

	"github.com/stretchr/testify/suite"
)

type MerchantDocumentUpdateTestSuite struct {
	tests.BaseTestSuite
	repo       *repository.Repositories
	userID     int
	merchantID int
	docA       *salemongo.MerchantDocument
	docB       *salemongo.MerchantDocument
}

func (s *MerchantDocumentUpdateTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupUserService()

	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	userClient := pbusers.NewUserQueryServiceClient(s.Conns["user"])
	s.repo = repository.NewRepositories(mongoDB, userClient)

	userID, err := salemongo.NextID(s.Ctx, mongoDB, "users")
	s.Require().NoError(err)

	now := time.Now().UTC()
	verified := true
	_, err = mongoDB.Collection("users").InsertOne(s.Ctx, salemongo.User{
		UserID:           userID,
		Firstname:        "Document",
		Lastname:         "Owner",
		Email:            "merchant.document@example.com",
		Password:         "password123",
		VerificationCode: "doc-verify",
		IsVerified:       &verified,
		CreatedAt:        &now,
		UpdatedAt:        &now,
	})
	s.Require().NoError(err)
	s.userID = int(userID)

	ctx := context.Background()

	merchant, err := s.repo.MerchantCommand.CreateMerchant(ctx, &requests.CreateMerchantRequest{
		UserID:       s.userID,
		Name:         "Document Merchant",
		Description:  "Doc",
		Address:      "Addr",
		ContactEmail: "doc@m.com",
		ContactPhone: "123",
		Status:       "active",
	})
	s.Require().NoError(err)
	s.merchantID = int(merchant.MerchantID)

	s.docA, err = s.repo.MerchantDocumentCommand.CreateMerchantDocument(ctx, &requests.CreateMerchantDocumentRequest{
		MerchantID:   s.merchantID,
		DocumentType: "doc_a_type",
		DocumentUrl:  "https://example.com/a.pdf",
	})
	s.Require().NoError(err)

	s.docB, err = s.repo.MerchantDocumentCommand.CreateMerchantDocument(ctx, &requests.CreateMerchantDocumentRequest{
		MerchantID:   s.merchantID,
		DocumentType: "doc_b_type",
		DocumentUrl:  "https://example.com/b.pdf",
	})
	s.Require().NoError(err)

	s.Require().NotEqual(s.docA.DocumentID, s.docB.DocumentID)
	s.Require().NotEqual(int(s.docB.DocumentID), s.merchantID,
		"fixture requires the target document ID to differ from the merchant ID")
}

func (s *MerchantDocumentUpdateTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *MerchantDocumentUpdateTestSuite) Test1_UpdateMerchantDocument_UsesDocumentID() {
	ctx := context.Background()
	docBID := int(s.docB.DocumentID)

	updated, err := s.repo.MerchantDocumentCommand.UpdateMerchantDocument(ctx, &requests.UpdateMerchantDocumentRequest{
		DocumentID:   &docBID,
		MerchantID:   s.merchantID,
		DocumentType: "doc_b_updated",
		DocumentUrl:  "https://example.com/b-updated.pdf",
		Status:       "verified",
		Note:         "Approved",
	})
	s.Require().NoError(err)
	s.Require().NotNil(updated)

	s.Equal(s.docB.DocumentID, updated.DocumentID)
	s.Equal("doc_b_updated", updated.DocumentType)
	s.Equal("verified", updated.Status)

	docA, err := s.repo.MerchantDocumentQuery.FindById(ctx, int(s.docA.DocumentID))
	s.Require().NoError(err)
	s.Require().NotNil(docA)
	s.Equal("doc_a_type", docA.DocumentType)
	s.NotEqual("doc_b_updated", docA.DocumentType)
}

func (s *MerchantDocumentUpdateTestSuite) Test2_UpdateMerchantDocumentStatus_UsesDocumentID() {
	ctx := context.Background()
	docBID := int(s.docB.DocumentID)

	updated, err := s.repo.MerchantDocumentCommand.UpdateMerchantDocumentStatus(ctx, &requests.UpdateMerchantDocumentStatusRequest{
		DocumentID: &docBID,
		MerchantID: s.merchantID,
		Status:     "rejected",
		Note:       "Invalid document",
	})
	s.Require().NoError(err)
	s.Require().NotNil(updated)

	s.Equal(s.docB.DocumentID, updated.DocumentID)
	s.Equal("rejected", updated.Status)

	docA, err := s.repo.MerchantDocumentQuery.FindById(ctx, int(s.docA.DocumentID))
	s.Require().NoError(err)
	s.Require().NotNil(docA)
	s.NotEqual("rejected", docA.Status)
}

func TestMerchantDocumentUpdateSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(MerchantDocumentUpdateTestSuite))
}
