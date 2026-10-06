package tests

import (
	"context"
	"time"

	"reflect"

	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type BaseTestSuite struct {
	suite.Suite
	ts      *TestSuite
	Log     logger.LoggerInterface
	Obs     observability.TraceLoggerObservability
	Conns   map[string]*grpc.ClientConn
	Servers []*grpc.Server
	Ctx     context.Context
	Cancel  context.CancelFunc
}

func (s *BaseTestSuite) SetupSuite() {
	s.Ctx, s.Cancel = context.WithCancel(context.Background())

	ts, err := SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	s.Log, _ = logger.NewLogger("test", nil)

	if s.Log == nil || (reflect.ValueOf(s.Log).Kind() == reflect.Ptr && reflect.ValueOf(s.Log).IsNil()) {
		z, _ := zap.NewDevelopment()
		s.Log = &logger.Logger{Log: z}
	}

	s.Obs, err = observability.NewObservability("test", s.Log)
	s.Require().NoError(err)
	s.Require().NotNil(s.Obs)
	s.Conns = make(map[string]*grpc.ClientConn)
}

func (s *BaseTestSuite) TearDownSuite() {
	for _, conn := range s.Conns {
		conn.Close()
	}
	for _, server := range s.Servers {
		server.GracefulStop()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
	if s.Cancel != nil {
		s.Cancel()
	}
}

func (s *BaseTestSuite) MongoDB() (*mongo.Database, error) {
	return s.ts.MongoDB()
}

func (s *BaseTestSuite) RedisClient() *goredis.Client {
	return s.ts.RedisClient()
}

func (s *BaseTestSuite) RegisterServer(server *grpc.Server) string {
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Servers = append(s.Servers, server)
	return addr
}

func (s *BaseTestSuite) GetConnection(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	s.Require().NoError(err)
	return conn
}

func (s *BaseTestSuite) SeedUser(ctx context.Context) int {
	db, err := s.MongoDB()
	s.Require().NoError(err)

	now := time.Now().UTC()
	for _, roleName := range []string{"Admin Access 1", "ROLE_ADMIN"} {
		roleID, err := salemongo.NextID(ctx, db, "roles")
		s.Require().NoError(err)
		_, err = db.Collection("roles").UpdateOne(ctx,
			bson.M{"role_name": roleName},
			bson.M{"$setOnInsert": salemongo.Role{
				RoleID:    roleID,
				RoleName:  roleName,
				CreatedAt: &now,
				UpdatedAt: &now,
			}},
			options.Update().SetUpsert(true),
		)
		s.Require().NoError(err)
	}

	res, err := pbusers.NewUserCommandServiceClient(s.Conns["user"]).Create(ctx, &pbusers.CreateUserRequest{
		Firstname:       "Seed",
		Lastname:        "User",
		Email:           "seed.user@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedCategory(ctx context.Context) int {
	res, err := pbcategories.NewCategoryCommandServiceClient(s.Conns["category"]).Create(ctx, &pbcategories.CreateCategoryRequest{
		Name:        "Seed Category",
		Description: "Seed Description",
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedMerchant(ctx context.Context, userID int) int {
	res, err := pbmerchants.NewMerchantCommandServiceClient(s.Conns["merchant"]).Create(ctx, &pbmerchants.CreateMerchantRequest{
		UserId:       int32(userID),
		Name:         "Seed Merchant",
		Description:  "Seed Description",
		Address:      "Seed Address",
		ContactEmail: "merchant@example.com",
		ContactPhone: "08123456789",
		Status:       "active",
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedProduct(ctx context.Context, merchantID int, categoryID int) int {
	res, err := pbproducts.NewProductCommandServiceClient(s.Conns["product"]).Create(ctx, &pbproducts.CreateProductRequest{
		MerchantId:   int32(merchantID),
		CategoryId:   int32(categoryID),
		Name:         "Seed Product",
		Description:  "Seed Description",
		Price:        10000,
		CountInStock: 100,
		Brand:        "Seed Brand",
		Weight:       1000,
		ImageProduct: "seed.jpg",
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedOrder(ctx context.Context, userID int, merchID int, prodID int) int {
	db, err := s.MongoDB()
	s.Require().NoError(err)

	cashiersColl := db.Collection("cashiers")
	var cashier salemongo.Cashier
	err = cashiersColl.FindOne(ctx, bson.M{
		"user_id":     int32(userID),
		"merchant_id": int32(merchID),
		"deleted_at":  nil,
	}).Decode(&cashier)

	var cashierID int32
	if err != nil {
		cashierID, err = salemongo.NextID(ctx, db, "cashiers")
		s.Require().NoError(err)

		now := time.Now().UTC()
		_, err = cashiersColl.InsertOne(ctx, salemongo.Cashier{
			CashierID:  cashierID,
			MerchantID: int32(merchID),
			UserID:     int32(userID),
			Name:       "Seed Cashier",
			CreatedAt:  &now,
			UpdatedAt:  &now,
		})
		s.Require().NoError(err)
	} else {
		cashierID = cashier.CashierID
	}

	res, err := pborders.NewOrderCommandServiceClient(s.Conns["order"]).Create(ctx, &pborders.CreateOrderRequest{
		MerchantId: int32(merchID),
		CashierId:  int32(cashierID),
		Items: []*pborders.CreateOrderItemRequest{
			{
				ProductId: int32(prodID),
				Quantity:  1,
			},
		},
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedOrderItem(ctx context.Context, orderID int, productID int) int {
	db, err := s.MongoDB()
	s.Require().NoError(err)

	id, err := salemongo.NextID(ctx, db, "order_items")
	s.Require().NoError(err)

	now := time.Now().UTC()
	_, err = db.Collection("order_items").InsertOne(ctx, salemongo.OrderItem{
		OrderItemID: id,
		OrderID:     int32(orderID),
		ProductID:   int32(productID),
		Quantity:    1,
		Price:       1000,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	})
	s.Require().NoError(err)
	return int(id)
}
