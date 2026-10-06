package outbox

import (
	"context"
	"errors"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// offlineDB returns a *mongo.Database that is never dialed. It is only used for
// argument-validation tests, which must not reach the server.
func offlineDB(t *testing.T) *mongo.Database {
	t.Helper()
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("failed to build mongo client: %v", err)
	}
	return client.Database("test")
}

func TestReserveValidatesKeys(t *testing.T) {
	for name, tc := range map[string]struct {
		db           *mongo.Database
		consumerName string
		eventKey     string
	}{
		"nil db":          {db: nil, consumerName: "email-service-group", eventKey: "topic:evt-1"},
		"empty consumer":  {db: offlineDB(t), consumerName: "", eventKey: "topic:evt-1"},
		"empty event key": {db: offlineDB(t), consumerName: "email-service-group", eventKey: ""},
	} {
		t.Run(name, func(t *testing.T) {
			if _, _, _, err := Reserve(context.Background(), tc.db, tc.consumerName, tc.eventKey, "topic", 0, 1); !errors.Is(err, ErrInvalidInboxKey) {
				t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
			}
		})
	}
}

func TestMarkProcessedValidatesKeys(t *testing.T) {
	if err := MarkProcessed(context.Background(), nil, "email-service-group", "topic:evt-1", 1); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}
	if err := MarkProcessed(context.Background(), offlineDB(t), "", "topic:evt-1", 1); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}
}

func TestReleaseValidatesKeys(t *testing.T) {
	if err := Release(context.Background(), nil, "email-service-group", "topic:evt-1", 1, nil); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}
	if err := Release(context.Background(), offlineDB(t), "", "topic:evt-1", 1, nil); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}
}

func TestNewMongoInboxRejectsNilDatabase(t *testing.T) {
	if _, err := NewMongoInbox(nil); err == nil {
		t.Fatal("expected an error when the database is nil")
	}
}

func TestMongoInboxValidatesKeys(t *testing.T) {
	inbox, err := NewMongoInbox(offlineDB(t))
	if err != nil {
		t.Fatalf("failed to build inbox: %v", err)
	}

	if _, _, _, err := inbox.Reserve(context.Background(), "", "topic:evt-1", "topic", 0, 1); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}
	if err := inbox.MarkProcessed(context.Background(), "", "topic:evt-1", 1); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}
	if err := inbox.Release(context.Background(), "", "topic:evt-1", 1, nil); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}
}
