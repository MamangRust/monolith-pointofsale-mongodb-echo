package tests

import (
	"bytes"
	"context"
	"mime/multipart"
	"net"

	pb "github.com/MamangRust/monolith-point-of-sale-pb"
	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pbmerchants_document "github.com/MamangRust/monolith-point-of-sale-pb/merchants_document"
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	pbtransactions "github.com/MamangRust/monolith-point-of-sale-pb/transactions"
	pbuserroles "github.com/MamangRust/monolith-point-of-sale-pb/user_roles"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"

	"github.com/MamangRust/monolith-point-of-sale-pkg/auth"
	"github.com/MamangRust/monolith-point-of-sale-pkg/hash"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Role
	role_cache "github.com/MamangRust/monolith-point-of-sale-role/cache"
	role_handler "github.com/MamangRust/monolith-point-of-sale-role/handler"
	role_repo "github.com/MamangRust/monolith-point-of-sale-role/repository"
	role_service "github.com/MamangRust/monolith-point-of-sale-role/service"

	// User
	user_cache "github.com/MamangRust/monolith-point-of-sale-user/cache"
	user_handler "github.com/MamangRust/monolith-point-of-sale-user/handler"
	user_repo "github.com/MamangRust/monolith-point-of-sale-user/repository"
	user_service "github.com/MamangRust/monolith-point-of-sale-user/service"

	// Auth
	auth_cache "github.com/MamangRust/monolith-point-of-sale-auth/cache"
	auth_handler "github.com/MamangRust/monolith-point-of-sale-auth/handler"
	auth_repo "github.com/MamangRust/monolith-point-of-sale-auth/repository"
	auth_service "github.com/MamangRust/monolith-point-of-sale-auth/service"

	// Category
	category_cache "github.com/MamangRust/monolith-point-of-sale-category/cache"
	category_handler "github.com/MamangRust/monolith-point-of-sale-category/handler"
	category_repo "github.com/MamangRust/monolith-point-of-sale-category/repository"
	category_service "github.com/MamangRust/monolith-point-of-sale-category/service"

	// Product
	product_cache "github.com/MamangRust/monolith-point-of-sale-product/cache"
	product_handler "github.com/MamangRust/monolith-point-of-sale-product/handler"
	product_repo "github.com/MamangRust/monolith-point-of-sale-product/repository"
	product_service "github.com/MamangRust/monolith-point-of-sale-product/service"

	// Merchant
	merchant_cache "github.com/MamangRust/monolith-point-of-sale-merchant/cache"
	merchant_handler "github.com/MamangRust/monolith-point-of-sale-merchant/handler"
	merchant_repo "github.com/MamangRust/monolith-point-of-sale-merchant/repository"
	merchant_service "github.com/MamangRust/monolith-point-of-sale-merchant/service"

	// Order
	order_cache "github.com/MamangRust/monolith-point-of-sale-order/cache"
	order_handler "github.com/MamangRust/monolith-point-of-sale-order/handler"
	order_repo "github.com/MamangRust/monolith-point-of-sale-order/repository"
	order_service "github.com/MamangRust/monolith-point-of-sale-order/service"

	// Transaction
	transaction_cache "github.com/MamangRust/monolith-point-of-sale-transacton/cache"
	transaction_handler "github.com/MamangRust/monolith-point-of-sale-transacton/handler"
	transaction_repo "github.com/MamangRust/monolith-point-of-sale-transacton/repository"
	transaction_service "github.com/MamangRust/monolith-point-of-sale-transacton/service"

	// Order Item
	order_item_cache "github.com/MamangRust/monolith-point-of-sale-order-item/cache"
	order_item_handler "github.com/MamangRust/monolith-point-of-sale-order-item/handler"
	order_item_repo "github.com/MamangRust/monolith-point-of-sale-order-item/repository"
	order_item_service "github.com/MamangRust/monolith-point-of-sale-order-item/service"

	// Cashier
	mencache "github.com/MamangRust/monolith-point-of-sale-cashier/cache"
	cashier_handler "github.com/MamangRust/monolith-point-of-sale-cashier/handler"
	cashier_repo "github.com/MamangRust/monolith-point-of-sale-cashier/repository"
	cashier_service "github.com/MamangRust/monolith-point-of-sale-cashier/service"

	"go.mongodb.org/mongo-driver/mongo"
)

func (s *BaseTestSuite) SetupRoleService() {
	cacheStore := s.GetCacheStore()
	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	roleMencache := role_cache.NewMencache(cacheStore)
	roleRepos := role_repo.NewRepositories(mongoDB)
	roleSvc := role_service.NewService(&role_service.Deps{
		Repositories:  roleRepos,
		Logger:        s.Log,
		Mencache:      roleMencache,
		Observability: s.Obs,
	})
	roleGapi := role_handler.NewHandler(&role_handler.Deps{
		Service: roleSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbroles.RegisterRoleQueryServiceServer(server, roleGapi.Role)
	pbroles.RegisterRoleCommandServiceServer(server, roleGapi.RoleCommand)
	pbuserroles.RegisterUserRoleServiceServer(server, roleGapi.UserRole)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["role"] = conn
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupUserService() {
	cacheStore := s.GetCacheStore()
	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)
	hasher := hash.NewHashingPassword()

	userMencache := user_cache.NewMencache(cacheStore)
	if _, ok := s.Conns["role"]; !ok {
		s.SetupRoleService()
	}
	roleQueryClient := pbroles.NewRoleQueryServiceClient(s.Conns["role"])
	userRoleClient := pbuserroles.NewUserRoleServiceClient(s.Conns["role"])
	userRepos := user_repo.NewRepositories(mongoDB, roleQueryClient, userRoleClient)
	userSvc := user_service.NewService(&user_service.Deps{
		Repositories:  userRepos,
		Logger:        s.Log,
		Hash:          hasher,
		Mencache:      userMencache,
		Observability: s.Obs,
	})
	userGapi := user_handler.NewHandler(&user_handler.Deps{
		Service: userSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbusers.RegisterUserQueryServiceServer(server, userGapi.User)
	pbusers.RegisterUserCommandServiceServer(server, userGapi.UserCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["user"] = conn
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupAuthService() {
	if _, ok := s.Conns["user"]; !ok {
		s.SetupUserService()
	}
	if _, ok := s.Conns["role"]; !ok {
		s.SetupRoleService()
	}

	cacheStore := s.GetCacheStore()
	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)
	hasher := hash.NewHashingPassword()
	tokenManager, _ := auth.NewManager("mysecret")

	userQueryClient := pbusers.NewUserQueryServiceClient(s.Conns["user"])
	userCommandClient := pbusers.NewUserCommandServiceClient(s.Conns["user"])
	roleQueryClient := pbroles.NewRoleQueryServiceClient(s.Conns["role"])
	userRoleClient := pbuserroles.NewUserRoleServiceClient(s.Conns["role"])
	authRepos := auth_repo.NewRepositories(mongoDB, userQueryClient, userCommandClient, roleQueryClient, userRoleClient)
	authMencache := auth_cache.NewMencache(cacheStore)
	authSvc := auth_service.NewService(&auth_service.Deps{
		Repositories:  authRepos,
		Logger:        s.Log,
		Mencache:      authMencache,
		Token:         tokenManager,
		Hash:          hasher,
		Kafka:         nil,
		Observability: s.Obs,
	})
	authGapi := auth_handler.NewAuthHandleGrpc(authSvc, s.Log)
	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, authGapi)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["auth"] = conn
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupCategoryService() {
	cacheStore := s.GetCacheStore()
	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	catMencache := category_cache.NewMencache(cacheStore)
	catRepos := category_repo.NewRepositories(mongoDB)
	catSvc := category_service.NewService(&category_service.Deps{
		Mencache:      catMencache,
		CacheStore:    cacheStore,
		Repositories:  catRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	catGapi := category_handler.NewHandler(&category_handler.Deps{
		Service: catSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbcategories.RegisterCategoryQueryServiceServer(server, catGapi.Category)
	pbcategories.RegisterCategoryCommandServiceServer(server, catGapi.CategoryCommand)
	pbcategories.RegisterCategoryStatsServiceServer(server, catGapi.CategoryStats)
	pbcategories.RegisterCategoryStatsByIdServiceServer(server, catGapi.CategoryStatsById)
	pbcategories.RegisterCategoryStatsByMerchantServiceServer(server, catGapi.CategoryStatsByMerchant)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["category"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupProductService() {
	if _, ok := s.Conns["category"]; !ok {
		s.SetupCategoryService()
	}
	if _, ok := s.Conns["merchant"]; !ok {
		s.SetupMerchantService()
	}

	cacheStore := s.GetCacheStore()
	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	prodMencache := product_cache.NewMencache(cacheStore)
	categoryClient := pbcategories.NewCategoryQueryServiceClient(s.Conns["category"])
	merchantClient := pbmerchants.NewMerchantQueryServiceClient(s.Conns["merchant"])
	prodRepos := product_repo.NewRepositories(mongoDB, categoryClient, merchantClient)
	prodSvc := product_service.NewService(&product_service.Deps{
		Mencache:      prodMencache,
		Repositories:  prodRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	prodGapi := product_handler.NewHandler(&product_handler.Deps{
		Service: prodSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbproducts.RegisterProductQueryServiceServer(server, prodGapi.Product)
	pbproducts.RegisterProductCommandServiceServer(server, prodGapi.ProductCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["product"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupMerchantService() {
	if _, ok := s.Conns["user"]; !ok {
		s.SetupUserService()
	}

	cacheStore := s.GetCacheStore()
	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	merchantMencache := merchant_cache.NewMencache(cacheStore)
	userClient := pbusers.NewUserQueryServiceClient(s.Conns["user"])
	merchantRepos := merchant_repo.NewRepositories(mongoDB, userClient)
	merchantSvc := merchant_service.NewService(&merchant_service.Deps{
		Mencache:      merchantMencache,
		Repositories:  merchantRepos,
		Logger:        s.Log,
		Observability: s.Obs,
		Kafka:         nil,
	})
	merchantGapi := merchant_handler.NewHandler(&merchant_handler.Deps{
		Service: merchantSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbmerchants.RegisterMerchantQueryServiceServer(server, merchantGapi.Merchant)
	pbmerchants.RegisterMerchantCommandServiceServer(server, merchantGapi.MerchantCommand)
	pbmerchants_document.RegisterMerchantDocumentQueryServiceServer(server, merchantGapi.MerchantDocument)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["merchant"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupOrderService() {
	if _, ok := s.Conns["cashier"]; !ok {
		s.SetupCashierService()
	}
	if _, ok := s.Conns["merchant"]; !ok {
		s.SetupMerchantService()
	}
	if _, ok := s.Conns["product"]; !ok {
		s.SetupProductService()
	}
	if _, ok := s.Conns["order-item"]; !ok {
		s.SetupOrderItemService()
	}

	cacheStore := s.GetCacheStore()
	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	orderMencache := order_cache.NewMencache(cacheStore)
	cashierClient := pbcashiers.NewCashierQueryServiceClient(s.Conns["cashier"])
	merchantClient := pbmerchants.NewMerchantQueryServiceClient(s.Conns["merchant"])
	productClient := pbproducts.NewProductQueryServiceClient(s.Conns["product"])
	productCommandClient := pbproducts.NewProductCommandServiceClient(s.Conns["product"])
	orderItemClient := pborder_items.NewOrderItemQueryServiceClient(s.Conns["order-item"])
	orderItemCommandClient := pborder_items.NewOrderItemCommandServiceClient(s.Conns["order-item"])
	orderRepos := order_repo.NewRepositories(mongoDB, cashierClient, merchantClient, productClient, productCommandClient, orderItemClient, orderItemCommandClient)
	orderSvc := order_service.NewService(&order_service.Deps{
		Mencache:      orderMencache,
		Repositories:  orderRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	orderGapi := order_handler.NewHandler(&order_handler.Deps{
		Service: orderSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pborders.RegisterOrderQueryServiceServer(server, orderGapi.Order)
	pborders.RegisterOrderCommandServiceServer(server, orderGapi.OrderCommand)
	pborders.RegisterOrderStatsServiceServer(server, orderGapi.OrderStats)
	pborders.RegisterOrderStatsByIdServiceServer(server, orderGapi.OrderStatsById)
	pborders.RegisterOrderStatsByMerchantServiceServer(server, orderGapi.OrderStatsByMerchant)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["order"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupTransactionService() {
	if _, ok := s.Conns["cashier"]; !ok {
		s.SetupCashierService()
	}
	if _, ok := s.Conns["merchant"]; !ok {
		s.SetupMerchantService()
	}
	if _, ok := s.Conns["order"]; !ok {
		s.SetupOrderService()
	}
	if _, ok := s.Conns["order-item"]; !ok {
		s.SetupOrderItemService()
	}

	cacheStore := s.GetCacheStore()
	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	transactionMencache := transaction_cache.NewMencache(cacheStore)
	cashierClient := pbcashiers.NewCashierQueryServiceClient(s.Conns["cashier"])
	merchantClient := pbmerchants.NewMerchantQueryServiceClient(s.Conns["merchant"])
	orderClient := pborders.NewOrderQueryServiceClient(s.Conns["order"])
	orderItemClient := pborder_items.NewOrderItemQueryServiceClient(s.Conns["order-item"])
	transactionRepos := transaction_repo.NewRepositories(mongoDB, cashierClient, merchantClient, orderClient, orderItemClient)
	transactionSvc := transaction_service.NewService(&transaction_service.Deps{
		Mencache:      transactionMencache,
		CacheStore:    cacheStore,
		Repositories:  transactionRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	transactionGapi := transaction_handler.NewHandler(&transaction_handler.Deps{
		Service: transactionSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbtransactions.RegisterTransactionQueryServiceServer(server, transactionGapi.Transaction)
	pbtransactions.RegisterTransactionCommandServiceServer(server, transactionGapi.TransactionCommand)
	pbtransactions.RegisterTransactionStatsServiceServer(server, transactionGapi.TransactionStats)
	pbtransactions.RegisterTransactionStatsByMerchantServiceServer(server, transactionGapi.TransactionStatsByMerchant)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["transaction"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupOrderItemService() {
	cacheStore := s.GetCacheStore()
	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	itemMencache := order_item_cache.NewMencache(cacheStore)
	itemRepos := order_item_repo.NewRepositories(mongoDB)
	itemSvc := order_item_service.NewService(&order_item_service.Deps{
		Mencache:      itemMencache,
		Repositories:  itemRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	itemGapi := order_item_handler.NewHandler(&order_item_handler.Deps{
		Service: itemSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pborder_items.RegisterOrderItemQueryServiceServer(server, itemGapi.OrderItem)
	pborder_items.RegisterOrderItemCommandServiceServer(server, itemGapi.OrderItemCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["order-item"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupCashierService() {
	if _, ok := s.Conns["user"]; !ok {
		s.SetupUserService()
	}
	if _, ok := s.Conns["merchant"]; !ok {
		s.SetupMerchantService()
	}

	cacheStore := s.GetCacheStore()
	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	cashierMencache := mencache.NewMencache(cacheStore)
	merchantClient := pbmerchants.NewMerchantQueryServiceClient(s.Conns["merchant"])
	userClient := pbusers.NewUserQueryServiceClient(s.Conns["user"])
	cashierRepos := cashier_repo.NewRepositories(mongoDB, merchantClient, userClient)
	cashierSvc := cashier_service.NewService(&cashier_service.Deps{
		Ctx:           context.Background(),
		Mencache:      cashierMencache,
		Repositories:  cashierRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	cashierGapi := cashier_handler.NewHandler(&cashier_handler.Deps{
		Service: cashierSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbcashiers.RegisterCashierQueryServiceServer(server, cashierGapi.Cashier)
	pbcashiers.RegisterCashierCommandServiceServer(server, cashierGapi.CashierCommand)
	pbcashiers.RegisterCashierStatsServiceServer(server, cashierGapi.CashierStats)
	pbcashiers.RegisterCashierStatsByIdServiceServer(server, cashierGapi.CashierStatsById)
	pbcashiers.RegisterCashierStatsByMerchantServiceServer(server, cashierGapi.CashierStatsByMerchant)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["cashier"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) dial(addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	return conn
}

func (s *BaseTestSuite) GetCacheStore() *cache.CacheStore {
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	return cache.NewCacheStore(s.ts.RedisClient(), s.Log, cacheMetrics)
}

func (s *BaseTestSuite) BuildMultipartRequestBody(fields map[string]string, fieldName, fileName string) ([]byte, string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range fields {
		fw, _ := w.CreateFormField(key)
		fw.Write([]byte(r))
	}
	fw, _ := w.CreateFormFile(fieldName, fileName)
	fw.Write([]byte("dummy image content"))
	w.Close()
	return b.Bytes(), w.FormDataContentType()
}

// IdentityClients holds the gRPC clients auth needs for users and roles, plus
// the shutdown hook for the temporary servers that back them.
type IdentityClients struct {
	UserQuery   pbusers.UserQueryServiceClient
	UserCommand pbusers.UserCommandServiceClient
	RoleQuery   pbroles.RoleQueryServiceClient
	UserRole    pbuserroles.UserRoleServiceClient
	stop        func()
}

// Stop shuts the temporary role and user servers down and closes their
// connections. It is safe to call on a nil receiver.
func (c *IdentityClients) Stop() {
	if c != nil && c.stop != nil {
		c.stop()
	}
}

// StartIdentityServices boots the role and user gRPC servers on ephemeral ports
// and returns clients for them.
//
// auth reaches users, roles and user_roles exclusively over gRPC, but the
// password-reset, auth-service and auth-repository suites own their own
// fixtures instead of embedding BaseTestSuite, so they cannot use its
// Setup*Service helpers. The user service needs the role service, hence both
// are started here in order.
func StartIdentityServices(
	mongoDB *mongo.Database,
	cacheStore *cache.CacheStore,
	log logger.LoggerInterface,
	obs observability.TraceLoggerObservability,
) (*IdentityClients, error) {
	roleSvc := role_service.NewService(&role_service.Deps{
		Repositories:  role_repo.NewRepositories(mongoDB),
		Logger:        log,
		Mencache:      role_cache.NewMencache(cacheStore),
		Observability: obs,
	})
	roleGapi := role_handler.NewHandler(&role_handler.Deps{
		Service: roleSvc,
		Logger:  log,
	})
	roleServer := grpc.NewServer()
	pbroles.RegisterRoleQueryServiceServer(roleServer, roleGapi.Role)
	pbroles.RegisterRoleCommandServiceServer(roleServer, roleGapi.RoleCommand)
	pbuserroles.RegisterUserRoleServiceServer(roleServer, roleGapi.UserRole)

	roleLis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}
	go roleServer.Serve(roleLis)

	roleConn, err := grpc.NewClient(roleLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		roleServer.Stop()
		return nil, err
	}

	roleQueryClient := pbroles.NewRoleQueryServiceClient(roleConn)
	userRoleClient := pbuserroles.NewUserRoleServiceClient(roleConn)

	userSvc := user_service.NewService(&user_service.Deps{
		Repositories:  user_repo.NewRepositories(mongoDB, roleQueryClient, userRoleClient),
		Logger:        log,
		Hash:          hash.NewHashingPassword(),
		Mencache:      user_cache.NewMencache(cacheStore),
		Observability: obs,
	})
	userGapi := user_handler.NewHandler(&user_handler.Deps{
		Service: userSvc,
		Logger:  log,
	})
	userServer := grpc.NewServer()
	pbusers.RegisterUserQueryServiceServer(userServer, userGapi.User)
	pbusers.RegisterUserCommandServiceServer(userServer, userGapi.UserCommand)

	userLis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		roleConn.Close()
		roleServer.Stop()
		return nil, err
	}
	go userServer.Serve(userLis)

	userConn, err := grpc.NewClient(userLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		userServer.Stop()
		roleConn.Close()
		roleServer.Stop()
		return nil, err
	}

	return &IdentityClients{
		UserQuery:   pbusers.NewUserQueryServiceClient(userConn),
		UserCommand: pbusers.NewUserCommandServiceClient(userConn),
		RoleQuery:   roleQueryClient,
		UserRole:    userRoleClient,
		stop: func() {
			userConn.Close()
			userServer.Stop()
			roleConn.Close()
			roleServer.Stop()
		},
	}, nil
}
