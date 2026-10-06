package handler

import (
	"context"
	errorsstd "errors"

	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-product/service"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/product_errors"
)

type productCommandHandleGrpc struct {
	pbproducts.UnimplementedProductCommandServiceServer
	productCommandService service.ProductCommandService
	logger                logger.LoggerInterface
}

func NewProductCommandHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbproducts.ProductCommandServiceServer {
	return &productCommandHandleGrpc{
		productCommandService: service.ProductCommand,
		logger:                logger,
	}
}

func (s *productCommandHandleGrpc) Create(ctx context.Context, request *pbproducts.CreateProductRequest) (*pbproducts.ApiResponseProduct, error) {
	s.logger.Info("Create product called", zap.String("name", request.GetName()))

	req := &requests.CreateProductRequest{
		MerchantID:   int(request.GetMerchantId()),
		CategoryID:   int(request.GetCategoryId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Price:        int(request.GetPrice()),
		CountInStock: int(request.GetCountInStock()),
		Brand:        request.GetBrand(),
		Weight:       int(request.GetWeight()),
		ImageProduct: request.GetImageProduct(),
	}

	if err := req.Validate(); err != nil {
		return nil, product_errors.ErrGrpcValidateCreateProduct
	}

	product, err := s.productCommandService.CreateProduct(ctx, req)
	if err != nil {
		s.logger.Error("Create product failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Create product success")

	return &pbproducts.ApiResponseProduct{
		Status:  "success",
		Message: "Successfully created product",
		Data:    mapResponseProduct(product),
	}, nil
}

func (s *productCommandHandleGrpc) Update(ctx context.Context, request *pbproducts.UpdateProductRequest) (*pbproducts.ApiResponseProduct, error) {
	s.logger.Info("Update product called", zap.Int32("id", request.GetProductId()))

	id := int(request.GetProductId())
	if id <= 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	req := &requests.UpdateProductRequest{
		ProductID:    &id,
		MerchantID:   int(request.GetMerchantId()),
		CategoryID:   int(request.GetCategoryId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Price:        int(request.GetPrice()),
		CountInStock: int(request.GetCountInStock()),
		Brand:        request.GetBrand(),
		Weight:       int(request.GetWeight()),
		ImageProduct: request.GetImageProduct(),
	}

	if err := req.Validate(); err != nil {
		return nil, product_errors.ErrGrpcValidateUpdateProduct
	}

	product, err := s.productCommandService.UpdateProduct(ctx, req)
	if err != nil {
		s.logger.Error("Update product failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Update product success")

	return &pbproducts.ApiResponseProduct{
		Status:  "success",
		Message: "Successfully updated product",
		Data:    mapResponseProduct(product),
	}, nil
}

// DecrementStock reserves stock for an order. Insufficient stock is reported
// as FailedPrecondition so consumers can distinguish it from transport errors.
func (s *productCommandHandleGrpc) DecrementStock(ctx context.Context, request *pbproducts.AdjustProductStockRequest) (*pbproducts.ApiResponseProduct, error) {
	s.logger.Info("DecrementStock called", zap.Int32("product_id", request.GetProductId()))

	if request.GetProductId() <= 0 || request.GetQuantity() <= 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	product, err := s.productCommandService.DecrementProductCountStock(ctx, int(request.GetProductId()), int(request.GetQuantity()))
	if err != nil {
		s.logger.Error("DecrementStock failed", zap.Error(err))
		if errorsstd.Is(err, product_errors.ErrInsufficientProductStock) {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DecrementStock success")

	return &pbproducts.ApiResponseProduct{
		Status:  "success",
		Message: "Successfully decremented product stock",
		Data:    mapResponseProduct(product),
	}, nil
}

// IncrementStock restores stock after a compensated order operation.
func (s *productCommandHandleGrpc) IncrementStock(ctx context.Context, request *pbproducts.AdjustProductStockRequest) (*pbproducts.ApiResponseProduct, error) {
	s.logger.Info("IncrementStock called", zap.Int32("product_id", request.GetProductId()))

	if request.GetProductId() <= 0 || request.GetQuantity() <= 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	product, err := s.productCommandService.IncrementProductCountStock(ctx, int(request.GetProductId()), int(request.GetQuantity()))
	if err != nil {
		s.logger.Error("IncrementStock failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("IncrementStock success")

	return &pbproducts.ApiResponseProduct{
		Status:  "success",
		Message: "Successfully incremented product stock",
		Data:    mapResponseProduct(product),
	}, nil
}

func (s *productCommandHandleGrpc) TrashedProduct(ctx context.Context, request *pbproducts.FindByIdProductRequest) (*pbproducts.ApiResponseProductDeleteAt, error) {
	s.logger.Info("TrashedProduct called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	product, err := s.productCommandService.TrashProduct(ctx, id)
	if err != nil {
		s.logger.Error("TrashedProduct failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedProduct success")

	return &pbproducts.ApiResponseProductDeleteAt{
		Status:  "success",
		Message: "Successfully trashed product",
		Data:    mapResponseProductDeleteAt(product),
	}, nil
}

func (s *productCommandHandleGrpc) RestoreProduct(ctx context.Context, request *pbproducts.FindByIdProductRequest) (*pbproducts.ApiResponseProductDeleteAt, error) {
	s.logger.Info("RestoreProduct called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	product, err := s.productCommandService.RestoreProduct(ctx, id)
	if err != nil {
		s.logger.Error("RestoreProduct failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreProduct success")

	return &pbproducts.ApiResponseProductDeleteAt{
		Status:  "success",
		Message: "Successfully restored product",
		Data:    mapResponseProductDeleteAt(product),
	}, nil
}

func (s *productCommandHandleGrpc) DeleteProductPermanent(ctx context.Context, request *pbproducts.FindByIdProductRequest) (*pbproducts.ApiResponseProductDelete, error) {
	s.logger.Info("DeleteProductPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, product_errors.ErrGrpcInvalidID
	}

	_, err := s.productCommandService.DeleteProductPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteProductPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteProductPermanent success")

	return &pbproducts.ApiResponseProductDelete{
		Status:  "success",
		Message: "Successfully deleted Product permanently",
	}, nil
}

func (s *productCommandHandleGrpc) RestoreAllProduct(ctx context.Context, _ *emptypb.Empty) (*pbproducts.ApiResponseProductAll, error) {
	s.logger.Info("RestoreAllProduct called")

	_, err := s.productCommandService.RestoreAllProducts(ctx)
	if err != nil {
		s.logger.Error("RestoreAllProduct failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllProduct success")

	return &pbproducts.ApiResponseProductAll{
		Status:  "success",
		Message: "Successfully restore all Product",
	}, nil
}

func (s *productCommandHandleGrpc) DeleteAllProductPermanent(ctx context.Context, _ *emptypb.Empty) (*pbproducts.ApiResponseProductAll, error) {
	s.logger.Info("DeleteAllProductPermanent called")

	_, err := s.productCommandService.DeleteAllProductsPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllProductPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllProductPermanent success")

	return &pbproducts.ApiResponseProductAll{
		Status:  "success",
		Message: "Successfully delete Product permanen",
	}, nil
}

// Map helpers
