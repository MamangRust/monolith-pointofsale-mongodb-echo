package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-point-of-sale-order-item/cache"
	"github.com/MamangRust/monolith-point-of-sale-order-item/repository"
	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-point-of-sale-shared/errorhandler"
	orderitem_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/order_item_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type orderItemCommandDeps struct {
	Cache         mencache.OrderItemCommandCache
	Repo          repository.OrderItemCommandRepository
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

type orderItemCommandService struct {
	mencache      mencache.OrderItemCommandCache
	repo          repository.OrderItemCommandRepository
	logger        logger.LoggerInterface
	observability observability.TraceLoggerObservability
}

func NewOrderItemCommandService(params *orderItemCommandDeps) OrderItemCommandService {
	return &orderItemCommandService{
		mencache:      params.Cache,
		repo:          params.Repo,
		logger:        params.Logger,
		observability: params.Observability,
	}
}

func (s *orderItemCommandService) CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*salemongo.OrderItem, error) {
	const method = "CreateOrderItem"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("orderID", req.OrderID),
		attribute.Int("productID", req.ProductID),
		attribute.Int("quantity", req.Quantity),
	)
	defer func() { end(status) }()

	orderItem, err := s.repo.CreateOrderItem(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*salemongo.OrderItem](
			s.logger,
			orderitem_errors.ErrFailedCreateOrderItem.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCachedOrderItems(ctx, req.OrderID)
	s.mencache.DeleteCachedOrderItemAllCache(ctx)

	logSuccess("Successfully created order item", zap.Int32("orderItemID", orderItem.OrderItemID))
	return orderItem, nil
}

func (s *orderItemCommandService) UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*salemongo.OrderItem, error) {
	const method = "UpdateOrderItem"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("orderItemID", req.OrderItemID),
		attribute.Int("orderID", req.OrderID),
	)
	defer func() { end(status) }()

	orderItem, err := s.repo.UpdateOrderItem(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*salemongo.OrderItem](
			s.logger,
			orderitem_errors.ErrFailedUpdateOrderItem.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCachedOrderItems(ctx, req.OrderID)
	s.mencache.DeleteCachedOrderItemAllCache(ctx)

	logSuccess("Successfully updated order item", zap.Int32("orderItemID", orderItem.OrderItemID))
	return orderItem, nil
}

// DeleteOrderItem drops one order item. order-service uses it as the
// compensating action when an order cannot be completed. The cache key for the
// owning order is unknown here (only the item id is known), so the broad
// order_item wildcard is invalidated instead.
func (s *orderItemCommandService) DeleteOrderItem(ctx context.Context, orderItemID int) error {
	const method = "DeleteOrderItem"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("orderItemID", orderItemID),
	)
	defer func() { end(status) }()

	if err := s.repo.DeleteOrderItem(ctx, orderItemID); err != nil {
		status = "error"
		_, handleErr := sharederrorhandler.HandleError[bool](
			s.logger,
			orderitem_errors.ErrFailedDeleteOrderItem.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
		return handleErr
	}

	s.mencache.DeleteCachedOrderItemAllCache(ctx)

	logSuccess("Successfully deleted order item", zap.Int("orderItemID", orderItemID))
	return nil
}
