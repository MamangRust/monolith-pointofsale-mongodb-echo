package repository

import (
	"context"
	"time"

	sharedErrors "github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type outboxRepository struct{}

func NewOutboxRepository() OutboxRepository {
	return &outboxRepository{}
}

func (r *outboxRepository) EnqueueInTx(ctx context.Context, tx *mongo.Database, topic, key string, payload []byte) error {
	coll := tx.Collection("outbox_events")
	now := time.Now()
	doc := bson.M{
		"topic":           topic,
		"event_key":       key,
		"payload":         payload,
		"status":          "pending",
		"attempts":        int32(0),
		"next_attempt_at": now,
		"created_at":      now,
		"updated_at":      now,
	}
	if _, err := coll.InsertOne(ctx, doc); err != nil {
		return sharedErrors.ErrInternal.WithInternal(err)
	}
	return nil
}
