package tests

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"

	salemongo "github.com/MamangRust/monolith-point-of-sale-pkg/database/mongo"
	goredis "github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type TestSuite struct {
	MongoContainer *mongodb.MongoDBContainer
	RedisContainer *redis.RedisContainer
	MongoURI       string
	RedisURL       string
	Ctx            context.Context

	client *mongo.Client
	db     *mongo.Database
}

func SetupTestSuite() (*TestSuite, error) {
	ctx := context.Background()

	// Setup MongoDB
	mongoContainer, err := mongodb.Run(ctx,
		"mongo:7.0",
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("27017").WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start mongo container: %w", err)
	}

	host, err := mongoContainer.Host(ctx)
	if err != nil {
		mongoContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to get mongo host: %w", err)
	}

	port, err := mongoContainer.MappedPort(ctx, "27017")
	if err != nil {
		mongoContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to get mongo port: %w", err)
	}

	mongoURI := fmt.Sprintf("mongodb://%s:%s", host, port.Port())

	// Setup Redis
	redisContainer, err := redis.Run(ctx, "redis:7-alpine")
	if err != nil {
		mongoContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to start redis container: %w", err)
	}

	redisURL, err := redisContainer.ConnectionString(ctx)
	if err != nil {
		mongoContainer.Terminate(ctx)
		redisContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to get redis connection string: %w", err)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		mongoContainer.Terminate(ctx)
		redisContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	db := client.Database("testdb")
	if err := salemongo.CreateIndexes(ctx, db); err != nil {
		client.Disconnect(ctx)
		mongoContainer.Terminate(ctx)
		redisContainer.Terminate(ctx)
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	ts := &TestSuite{
		MongoContainer: mongoContainer,
		RedisContainer: redisContainer,
		MongoURI:       mongoURI,
		RedisURL:       redisURL,
		Ctx:            ctx,
		client:         client,
		db:             db,
	}

	if err := ts.SeedMinimalRoles(); err != nil {
		ts.Teardown()
		return nil, fmt.Errorf("failed to seed minimal roles: %w", err)
	}

	return ts, nil
}

// SeedMinimalRoles inserts the default roles required by the auth service
// (register assigns ROLE_ADMIN by default). Seeders create no schema, so
// fixture roles must exist for integration tests.
func (ts *TestSuite) SeedMinimalRoles() error {
	db, err := ts.MongoDB()
	if err != nil {
		return fmt.Errorf("failed to open mongo connection: %w", err)
	}

	now := time.Now().UTC()
	for _, roleName := range []string{"ROLE_ADMIN", "ROLE_CASHIER"} {
		roleID, err := salemongo.NextID(ts.Ctx, db, "roles")
		if err != nil {
			return fmt.Errorf("failed to generate role id: %w", err)
		}

		if _, err := db.Collection("roles").UpdateOne(ts.Ctx,
			map[string]any{"role_name": roleName},
			map[string]any{"$setOnInsert": salemongo.Role{
				RoleID:    roleID,
				RoleName:  roleName,
				CreatedAt: &now,
				UpdatedAt: &now,
			}},
			options.Update().SetUpsert(true),
		); err != nil {
			return fmt.Errorf("failed to seed role %s: %w", roleName, err)
		}
	}

	return nil
}

func (ts *TestSuite) MongoDB() (*mongo.Database, error) {
	if ts.db == nil {
		return nil, fmt.Errorf("mongo database is not initialized")
	}
	return ts.db, nil
}

func (ts *TestSuite) RedisClient() *goredis.Client {
	opts, _ := goredis.ParseURL(ts.RedisURL)
	return goredis.NewClient(opts)
}

func (ts *TestSuite) Teardown() {
	if ts.client != nil {
		if err := ts.client.Disconnect(ts.Ctx); err != nil {
			log.Printf("failed to disconnect mongo client: %v", err)
		}
	}
	if ts.MongoContainer != nil {
		if err := ts.MongoContainer.Terminate(ts.Ctx); err != nil {
			log.Printf("failed to terminate mongo container: %v", err)
		}
	}
	if ts.RedisContainer != nil {
		if err := ts.RedisContainer.Terminate(ts.Ctx); err != nil {
			log.Printf("failed to terminate redis container: %v", err)
		}
	}
}

type GRPCServerRunner interface {
	Serve(lis net.Listener) error
	Stop()
	GracefulStop()
}

func RunGRPCServer(server *grpc.Server) (string, error) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", err
	}
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Printf("grpc server error: %v", err)
		}
	}()
	return lis.Addr().String(), nil
}
