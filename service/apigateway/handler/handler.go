package handler

import (
	"fmt"
	"github.com/MamangRust/monolith-point-of-sale-pb"
	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pbmerchants_document "github.com/MamangRust/monolith-point-of-sale-pb/merchants_document"
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	pbtransactions "github.com/MamangRust/monolith-point-of-sale-pb/transactions"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"strconv"

	auth_cache "github.com/MamangRust/monolith-point-of-sale-apigateway/cache/auth"
	gateway_cache "github.com/MamangRust/monolith-point-of-sale-apigateway/cache/gateway_cache"
	"github.com/MamangRust/monolith-point-of-sale-pkg/auth"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-pkg/upload_image"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"

	response_api "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/api"
)

type ServiceConnections struct {
	Auth        *grpc.ClientConn
	Role        *grpc.ClientConn
	User        *grpc.ClientConn
	Cashier     *grpc.ClientConn
	Category    *grpc.ClientConn
	Merchant    *grpc.ClientConn
	OrderItem   *grpc.ClientConn
	Order       *grpc.ClientConn
	Product     *grpc.ClientConn
	Transaction *grpc.ClientConn
}

type Deps struct {
	Token              auth.TokenManager
	E                  *echo.Echo
	Logger             logger.LoggerInterface
	Mapping            *response_api.ResponseApiMapper
	ImageUpload        upload_image.ImageUploads
	ServiceConnections *ServiceConnections
	ApiHandler         errors.ApiHandler
	AuthCache          auth_cache.AuthMencache
	GatewayCache       *gateway_cache.GatewayCache
}

func NewHandler(deps *Deps) {
	clientAuth := pb.NewAuthServiceClient(deps.ServiceConnections.Auth)
	clientRole := pbroles.NewRoleQueryServiceClient(deps.ServiceConnections.Role)
	clientRoleCommand := pbroles.NewRoleCommandServiceClient(deps.ServiceConnections.Role)
	clientUser := pbusers.NewUserQueryServiceClient(deps.ServiceConnections.User)
	clientUserCommand := pbusers.NewUserCommandServiceClient(deps.ServiceConnections.User)
	clientCategory := pbcategories.NewCategoryQueryServiceClient(deps.ServiceConnections.Category)
	clientCategoryCommand := pbcategories.NewCategoryCommandServiceClient(deps.ServiceConnections.Category)
	clientCategoryStats := pbcategories.NewCategoryStatsServiceClient(deps.ServiceConnections.Category)
	clientCategoryStatsID := pbcategories.NewCategoryStatsByIdServiceClient(deps.ServiceConnections.Category)
	clientCategoryStatsM := pbcategories.NewCategoryStatsByMerchantServiceClient(deps.ServiceConnections.Category)
	clientCashier := pbcashiers.NewCashierQueryServiceClient(deps.ServiceConnections.Cashier)
	clientCashierCommand := pbcashiers.NewCashierCommandServiceClient(deps.ServiceConnections.Cashier)
	clientCashierStats := pbcashiers.NewCashierStatsServiceClient(deps.ServiceConnections.Cashier)
	clientCashierStatsID := pbcashiers.NewCashierStatsByIdServiceClient(deps.ServiceConnections.Cashier)
	clientCashierStatsM := pbcashiers.NewCashierStatsByMerchantServiceClient(deps.ServiceConnections.Cashier)
	clientMerchant := pbmerchants.NewMerchantQueryServiceClient(deps.ServiceConnections.Merchant)
	clientMerchantCommand := pbmerchants.NewMerchantCommandServiceClient(deps.ServiceConnections.Merchant)
	clientMerchantDocument := pbmerchants_document.NewMerchantDocumentQueryServiceClient(deps.ServiceConnections.Merchant)
	clientMerchantDocumentCommand := pbmerchants_document.NewMerchantDocumentCommandServiceClient(deps.ServiceConnections.Merchant)
	clientOrderItem := pborder_items.NewOrderItemQueryServiceClient(deps.ServiceConnections.OrderItem)
	clientOrder := pborders.NewOrderQueryServiceClient(deps.ServiceConnections.Order)
	clientOrderCommand := pborders.NewOrderCommandServiceClient(deps.ServiceConnections.Order)
	clientOrderStats := pborders.NewOrderStatsServiceClient(deps.ServiceConnections.Order)
	clientOrderStatsID := pborders.NewOrderStatsByIdServiceClient(deps.ServiceConnections.Order)
	clientOrderStatsM := pborders.NewOrderStatsByMerchantServiceClient(deps.ServiceConnections.Order)
	clientProduct := pbproducts.NewProductQueryServiceClient(deps.ServiceConnections.Product)
	clientProductCommand := pbproducts.NewProductCommandServiceClient(deps.ServiceConnections.Product)
	clientTransaction := pbtransactions.NewTransactionQueryServiceClient(deps.ServiceConnections.Transaction)
	clientTransactionCommand := pbtransactions.NewTransactionCommandServiceClient(deps.ServiceConnections.Transaction)
	clientTransactionStats := pbtransactions.NewTransactionStatsServiceClient(deps.ServiceConnections.Transaction)
	clientTransactionStatsM := pbtransactions.NewTransactionStatsByMerchantServiceClient(deps.ServiceConnections.Transaction)

	NewHandlerAuth(deps.E, clientAuth, deps.Logger, deps.Mapping.AuthResponseMapper, deps.ApiHandler, deps.AuthCache)
	NewHandlerRole(deps.E, clientRole, clientRoleCommand, deps.Logger, deps.Mapping.RoleResponseMapper, deps.ApiHandler, deps.GatewayCache)
	NewHandlerUser(deps.E, clientUser, clientUserCommand, deps.Logger, deps.Mapping.UserResponseMapper, deps.ApiHandler, deps.GatewayCache)
	NewHandlerCategory(deps.E, clientCategory, clientCategoryCommand, clientCategoryStats, clientCategoryStatsID, clientCategoryStatsM, deps.Logger, deps.Mapping.CategoryResponseMapper, deps.ApiHandler, deps.GatewayCache)
	NewHandlerCashier(deps.E, clientCashier, clientCashierCommand, clientCashierStats, clientCashierStatsID, clientCashierStatsM, deps.Logger, deps.Mapping.CashierResponseMapper, deps.ApiHandler, deps.GatewayCache)
	NewHandlerMerchant(deps.E, clientMerchant, clientMerchantCommand, deps.Logger, deps.Mapping.MerchantResponseMapper, deps.ApiHandler, deps.GatewayCache)
	NewHandlerMerchantDocument(deps.E, clientMerchantDocument, clientMerchantDocumentCommand, deps.Logger, deps.Mapping.MerchantDocumentProMapper, deps.ApiHandler, deps.GatewayCache)
	NewHandlerOrderItem(deps.E, clientOrderItem, deps.Logger, deps.Mapping.OrderItemResponseMapper, deps.ApiHandler, deps.GatewayCache)
	NewHandlerOrder(deps.E, clientOrder, clientOrderCommand, clientOrderStats, clientOrderStatsID, clientOrderStatsM, deps.Logger, deps.Mapping.OrderResponseMapper, deps.ApiHandler, deps.GatewayCache)
	NewHandlerProduct(deps.E, clientProduct, clientProductCommand, deps.Logger, deps.Mapping.ProductResponseMapper, deps.ImageUpload, deps.ApiHandler, deps.GatewayCache)
	NewHandlerTransaction(deps.E, clientTransaction, clientTransactionCommand, clientTransactionStats, clientTransactionStatsM, deps.Logger, deps.Mapping.TransactionResponseMapper, deps.ApiHandler, deps.GatewayCache)
}

func parseQueryInt(c echo.Context, key string, defaultValue int) int {
	val, err := strconv.Atoi(c.QueryParam(key))
	if err != nil || val <= 0 {
		return defaultValue
	}
	return val
}

func parseQueryIntWithValidation(c echo.Context, key string, min, max int) (int, error) {
	valStr := c.QueryParam(key)
	val, err := strconv.Atoi(valStr)
	if err != nil || val < min || val > max {
		return 0, fmt.Errorf("invalid %s: %s", key, valStr)
	}
	return val, nil
}
