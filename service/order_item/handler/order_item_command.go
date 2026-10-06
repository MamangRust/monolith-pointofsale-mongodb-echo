package handler

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-order-item/service"
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	orderitem_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/order_item_errors"
	"go.uber.org/zap"
)

type orderItemCommandHandleGrpc struct {
	pborder_items.UnimplementedOrderItemCommandServiceServer
	orderItemCommandService service.OrderItemCommandService
	logger                  logger.LoggerInterface
}

func NewOrderItemCommandHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pborder_items.OrderItemCommandServiceServer {
	return &orderItemCommandHandleGrpc{
		orderItemCommandService: service.OrderItemCommand,
		logger:                  logger,
	}
}

func (s *orderItemCommandHandleGrpc) CreateOrderItem(ctx context.Context, request *pborder_items.CreateOrderItemRecordRequest) (*pborder_items.ApiResponseOrderItem, error) {
	s.logger.Info("CreateOrderItem called", zap.Int32("order_id", request.GetOrderId()))

	if request.GetOrderId() <= 0 || request.GetProductId() <= 0 {
		return nil, orderitem_errors.ErrGrpcInvalidID
	}

	req := &requests.CreateOrderItemRecordRequest{
		OrderID:   int(request.GetOrderId()),
		ProductID: int(request.GetProductId()),
		Quantity:  int(request.GetQuantity()),
		Price:     int(request.GetPrice()),
	}

	if req.Quantity <= 0 {
		return nil, orderitem_errors.ErrFailedInvalidQuantity
	}

	orderItem, err := s.orderItemCommandService.CreateOrderItem(ctx, req)
	if err != nil {
		s.logger.Error("CreateOrderItem failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("CreateOrderItem success")

	return &pborder_items.ApiResponseOrderItem{
		Status:  "success",
		Message: "Successfully created order item",
		Data:    mapOrderItemToProto(orderItem),
	}, nil
}

func (s *orderItemCommandHandleGrpc) UpdateOrderItem(ctx context.Context, request *pborder_items.UpdateOrderItemRecordRequest) (*pborder_items.ApiResponseOrderItem, error) {
	s.logger.Info("UpdateOrderItem called", zap.Int32("order_item_id", request.GetOrderItemId()))

	if request.GetOrderItemId() <= 0 || request.GetOrderId() <= 0 || request.GetProductId() <= 0 {
		return nil, orderitem_errors.ErrGrpcInvalidID
	}

	req := &requests.UpdateOrderItemRecordRequest{
		OrderItemID: int(request.GetOrderItemId()),
		OrderID:     int(request.GetOrderId()),
		ProductID:   int(request.GetProductId()),
		Quantity:    int(request.GetQuantity()),
		Price:       int(request.GetPrice()),
	}

	if req.Quantity <= 0 {
		return nil, orderitem_errors.ErrFailedInvalidQuantity
	}

	orderItem, err := s.orderItemCommandService.UpdateOrderItem(ctx, req)
	if err != nil {
		s.logger.Error("UpdateOrderItem failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("UpdateOrderItem success")

	return &pborder_items.ApiResponseOrderItem{
		Status:  "success",
		Message: "Successfully updated order item",
		Data:    mapOrderItemToProto(orderItem),
	}, nil
}

func (s *orderItemCommandHandleGrpc) DeleteOrderItem(ctx context.Context, request *pborder_items.DeleteOrderItemRecordRequest) (*pborder_items.ApiResponseOrderItemDelete, error) {
	s.logger.Info("DeleteOrderItem called", zap.Int32("order_item_id", request.GetOrderItemId()))

	if request.GetOrderItemId() <= 0 {
		return nil, orderitem_errors.ErrGrpcInvalidID
	}

	if err := s.orderItemCommandService.DeleteOrderItem(ctx, int(request.GetOrderItemId())); err != nil {
		s.logger.Error("DeleteOrderItem failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteOrderItem success")

	return &pborder_items.ApiResponseOrderItemDelete{
		Status:  "success",
		Message: "Successfully deleted order item",
	}, nil
}
