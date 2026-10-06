package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	inboxStatusProcessing = "processing"
	inboxStatusProcessed  = "processed"
	inboxStatusPending    = "pending"

	// inboxLeaseDuration extends a reservation lease by one minute, matching
	// the original sqlc consumer_inbox queries (current_timestamp + interval
	// '1 minute').
	inboxLeaseDuration = time.Minute

	// inboxReleaseDefaultError is stored when Release receives a nil error,
	// matching the original outbox.Release behavior.
	inboxReleaseDefaultError = "consumer processing failed"

	inboxCollection = "consumer_inbox"
)

// inboxDoc mirrors the durable consumer reservation document.
type inboxDoc struct {
	ConsumerName       string     `bson:"consumer_name"`
	EventKey           string     `bson:"event_key"`
	Topic              string     `bson:"topic"`
	PartitionID        int32      `bson:"partition_id"`
	MessageOffset      int64      `bson:"message_offset"`
	Status             string     `bson:"status"`
	Attempts           int64      `bson:"attempts"`
	ReservationVersion int64      `bson:"reservation_version"`
	LeaseUntil         *time.Time `bson:"lease_until"`
	LastError          string     `bson:"last_error"`
	ProcessedAt        *time.Time `bson:"processed_at"`
	CreatedAt          *time.Time `bson:"created_at"`
	UpdatedAt          *time.Time `bson:"updated_at"`
}

type consumerInboxRepository struct {
	coll *mongodriver.Collection
}

// NewConsumerInboxRepository builds the MongoDB-backed ConsumerInbox.
func NewConsumerInboxRepository(db *mongodriver.Database) ConsumerInbox {
	if db == nil {
		return &consumerInboxRepository{}
	}
	return &consumerInboxRepository{coll: db.Collection(inboxCollection)}
}

// Reserve claims an event for a consumer. It mirrors the original sqlc
// ReserveConsumerInbox statement: the guarded upsert only updates the
// conflicting document when it is neither processed nor under an active
// lease; otherwise the unique (consumer_name, event_key) index rejects the
// insert and the existing state is reported so callers dedupe or redeliver
// instead of double-sending.
func (r *consumerInboxRepository) Reserve(ctx context.Context, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error) {
	if r.coll == nil || consumerName == "" || eventKey == "" {
		return false, false, 0, ErrInvalidInboxKey
	}

	now := time.Now().UTC()

	filter := bson.M{
		"consumer_name": consumerName,
		"event_key":     eventKey,
		"status":        bson.M{"$ne": inboxStatusProcessed},
		"lease_until":   bson.M{"$lte": now},
	}

	update := bson.M{
		"$set": bson.M{
			"consumer_name":  consumerName,
			"event_key":      eventKey,
			"status":         inboxStatusProcessing,
			"topic":          topic,
			"partition_id":   partition,
			"message_offset": offset,
			"lease_until":    now.Add(inboxLeaseDuration),
			"last_error":     "",
			"processed_at":   nil,
			"updated_at":     now,
		},
		"$inc":         bson.M{"attempts": 1, "reservation_version": 1},
		"$setOnInsert": bson.M{"created_at": now},
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var claimed inboxDoc
	err := r.coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&claimed)
	if err == nil {
		return true, false, claimed.ReservationVersion, nil
	}
	if !mongodriver.IsDuplicateKeyError(err) && !errors.Is(err, mongodriver.ErrNoDocuments) {
		return false, false, 0, err
	}

	// The conflicting row rejected the guard: it is either already processed
	// or still under an active lease.
	var existing inboxDoc
	if err := r.coll.FindOne(ctx, bson.M{"consumer_name": consumerName, "event_key": eventKey}).Decode(&existing); err != nil {
		if errors.Is(err, mongodriver.ErrNoDocuments) {
			// Nothing to dedupe against yet; the next delivery retries cleanly.
			return false, false, 0, nil
		}
		return false, false, 0, err
	}
	return false, existing.Status == inboxStatusProcessed, existing.ReservationVersion, nil
}

// MarkProcessed completes only the active reservation. Like the original sqlc
// :exec query, a version mismatch is a silent no-op rather than an error.
func (r *consumerInboxRepository) MarkProcessed(ctx context.Context, consumerName, eventKey string, reservationVersion int64) error {
	if r.coll == nil || consumerName == "" || eventKey == "" {
		return ErrInvalidInboxKey
	}

	now := time.Now().UTC()
	_, err := r.coll.UpdateOne(ctx,
		bson.M{
			"consumer_name":       consumerName,
			"event_key":           eventKey,
			"status":              inboxStatusProcessing,
			"reservation_version": reservationVersion,
		},
		bson.M{"$set": bson.M{
			"status":       inboxStatusProcessed,
			"processed_at": now,
			"lease_until":  now,
			"last_error":   "",
			"updated_at":   now,
		}},
	)
	return err
}

// Release releases only the active reservation so the retry path can re-claim
// the event. A stale version is a no-op, matching the original query.
func (r *consumerInboxRepository) Release(ctx context.Context, consumerName, eventKey string, reservationVersion int64, processingErr error) error {
	if r.coll == nil || consumerName == "" || eventKey == "" {
		return ErrInvalidInboxKey
	}

	lastError := inboxReleaseDefaultError
	if processingErr != nil {
		lastError = processingErr.Error()
	}

	now := time.Now().UTC()
	_, err := r.coll.UpdateOne(ctx,
		bson.M{
			"consumer_name":       consumerName,
			"event_key":           eventKey,
			"status":              inboxStatusProcessing,
			"reservation_version": reservationVersion,
		},
		bson.M{"$set": bson.M{
			"status":      inboxStatusPending,
			"lease_until": now,
			"last_error":  lastError,
			"updated_at":  now,
		}},
	)
	return err
}
