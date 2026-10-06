package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	pbuserroles "github.com/MamangRust/monolith-point-of-sale-pb/user_roles"
	"net/http"
	"net/http/httptest"
	"testing"

	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"

	apigw_cache "github.com/MamangRust/monolith-point-of-sale-apigateway/cache/gateway_cache"
	apigateway "github.com/MamangRust/monolith-point-of-sale-apigateway/handler"
	"github.com/MamangRust/monolith-point-of-sale-pkg/hash"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	apierror "github.com/MamangRust/monolith-point-of-sale-shared/errors"
	response_api "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/api"
	tests "github.com/MamangRust/monolith-point-of-sale-test"
	user_cache "github.com/MamangRust/monolith-point-of-sale-user/cache"
	gapi "github.com/MamangRust/monolith-point-of-sale-user/handler"
	"github.com/MamangRust/monolith-point-of-sale-user/repository"
	"github.com/MamangRust/monolith-point-of-sale-user/service"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type UserHandlerTestSuite struct {
	tests.BaseTestSuite
	client      pbusers.UserCommandServiceClient
	queryClient pbusers.UserQueryServiceClient
	router      *echo.Echo
	userID      int
	userEmail   string
}

func (s *UserHandlerTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	s.SetupRoleService()

	mongoDB, err := s.MongoDB()
	s.Require().NoError(err)

	// Seed default role required by user service CreateUser
	s.Require().NoError(s.Fixtures().SeedRoleRow(s.Ctx, "ROLE_ADMIN"))

	repos := repository.NewRepositories(mongoDB, pbroles.NewRoleQueryServiceClient(s.Conns["role"]), pbuserroles.NewUserRoleServiceClient(s.Conns["role"]))

	log, _ := logger.NewLogger("test", nil)
	hasher := hash.NewHashingPassword()
	cacheStore := s.GetCacheStore()
	mencache := user_cache.NewMencache(cacheStore)

	userService := service.NewService(&service.Deps{
		Mencache:      mencache,
		Repositories:  repos,
		Hash:          hasher,
		Logger:        log,
		Observability: s.Obs,
	})

	// Start gRPC Server
	userHandler := gapi.NewHandler(&gapi.Deps{
		Service: userService,
		Logger:  log,
	})
	server := grpc.NewServer()
	pbusers.RegisterUserQueryServiceServer(server, userHandler.User)
	pbusers.RegisterUserCommandServiceServer(server, userHandler.UserCommand)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)
	s.client = pbusers.NewUserCommandServiceClient(conn)
	s.queryClient = pbusers.NewUserQueryServiceClient(conn)

	// Setup Echo
	s.router = echo.New()
	apigateway.NewHandlerUser(
		s.router,
		pbusers.NewUserQueryServiceClient(conn),
		pbusers.NewUserCommandServiceClient(conn),
		log,
		response_api.NewUserResponseMapper(),
		apierror.NewApiHandler(s.Obs, log),
		apigw_cache.NewGatewayCache(cacheStore),
	)
}

func (s *UserHandlerTestSuite) TestUserApiLifecycle() {
	// 1. Create
	s.userEmail = "handler.user@example.com"
	createReq := requests.CreateUserRequest{
		FirstName:       "Handler",
		LastName:        "User",
		Email:           s.userEmail,
		Password:        "password123",
		ConfirmPassword: "password123",
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/user/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Require().Equal(http.StatusOK, rec.Code, rec.Body.String())
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})
	s.userID = int(data["id"].(float64))

	// 2. FindAll
	req = httptest.NewRequest(http.MethodGet, "/api/user", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. FindById
	s.Require().NotZero(s.userID)
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/user/%d", s.userID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. FindByActive
	req = httptest.NewRequest(http.MethodGet, "/api/user/active", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. Update
	updateReq := requests.UpdateUserRequest{
		FirstName:       "Updated",
		LastName:        "User",
		Email:           s.userEmail,
		Password:        "password123",
		ConfirmPassword: "password123",
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/user/update/%d", s.userID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 6. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/user/trashed/%d", s.userID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 7. FindByTrashed
	req = httptest.NewRequest(http.MethodGet, "/api/user/trashed", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 8. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/user/restore/%d", s.userID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 9. DeletePermanent
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/user/trashed/%d", s.userID), nil)
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/user/permanent/%d", s.userID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 10. RestoreAll
	req = httptest.NewRequest(http.MethodPost, "/api/user/restore/all", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 11. DeleteAll
	req = httptest.NewRequest(http.MethodPost, "/api/user/permanent/all", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestUserHandlerSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(UserHandlerTestSuite))
}
