package outbox_test

import (
	"context"
	"errors"
	"testing"
	"time"

	tests "github.com/MamangRust/monolith-point-of-sale-test"

	"github.com/MamangRust/monolith-point-of-sale-pkg/outbox"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	consumerName   = "email-service-group"
	inboxCollName  = "consumer_inbox"
	statusProcess  = "processing"
	statusFinished = "processed"
	statusPending  = "pending"
)

type InboxTestSuite struct {
	suite.Suite
	ts *tests.TestSuite
}

func (s *InboxTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts
}

func (s *InboxTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *InboxTestSuite) inboxRow(ctx context.Context, eventKey string) bson.M {
	db, err := s.ts.MongoDB()
	s.Require().NoError(err)

	var row bson.M
	err = db.Collection(inboxCollName).FindOne(ctx, bson.M{
		"consumer_name": consumerName,
		"event_key":     eventKey,
	}).Decode(&row)
	s.Require().NoError(err)
	return row
}

func (s *InboxTestSuite) TestReserveClaimsAndFencesEvents() {
	ctx := context.Background()
	db, err := s.ts.MongoDB()
	s.Require().NoError(err)

	eventKey := "order-topic:evt-1"

	inserted, processed, version, err := outbox.Reserve(ctx, db, consumerName, eventKey, "order-topic", 0, 1)
	s.Require().NoError(err)
	s.True(inserted, "first reserve must insert a new inbox row")
	s.False(processed)
	s.Equal(int64(1), version)

	// A second worker cannot insert, but it can see the row is still in flight.
	inserted, processed, version, err = outbox.Reserve(ctx, db, consumerName, eventKey, "order-topic", 0, 1)
	s.Require().NoError(err)
	s.False(inserted)
	s.False(processed)
	s.Equal(int64(2), version, "reserving an in-flight event must bump the reservation version")

	row := s.inboxRow(ctx, eventKey)
	s.Equal(statusProcess, row["status"])
	s.Equal(int32(2), row["attempts"])

	// Marking processed with a stale reservation must be a no-op.
	s.Require().NoError(outbox.MarkProcessed(ctx, db, consumerName, eventKey, version-1))
	row = s.inboxRow(ctx, eventKey)
	s.Equal(statusProcess, row["status"], "stale reservations must not mark the event processed")

	// The active reservation wins.
	s.Require().NoError(outbox.MarkProcessed(ctx, db, consumerName, eventKey, version))
	row = s.inboxRow(ctx, eventKey)
	s.Equal(statusFinished, row["status"])
	s.NotNil(row["processed_at"])

	// Once processed the event is reported as a duplicate.
	inserted, processed, _, err = outbox.Reserve(ctx, db, consumerName, eventKey, "order-topic", 0, 1)
	s.NoError(err)
	s.False(inserted)
	s.True(processed)
}

func (s *InboxTestSuite) TestReleaseAllowsRetryOnlyForActiveReservation() {
	ctx := context.Background()
	db, err := s.ts.MongoDB()
	s.Require().NoError(err)

	eventKey := "order-topic:evt-2"

	_, _, version, err := outbox.Reserve(ctx, db, consumerName, eventKey, "order-topic", 0, 7)
	s.Require().NoError(err)

	// A stale release must not touch the active reservation.
	s.Require().NoError(outbox.Release(ctx, db, consumerName, eventKey, version+1, errors.New("stale worker")))
	row := s.inboxRow(ctx, eventKey)
	s.Equal(statusProcess, row["status"])

	// The active release hands the event back for retry and records the error.
	processingErr := errors.New("smtp down")
	s.Require().NoError(outbox.Release(ctx, db, consumerName, eventKey, version, processingErr))

	row = s.inboxRow(ctx, eventKey)
	s.Equal(statusPending, row["status"])
	s.Equal(processingErr.Error(), row["last_error"])

	leaseUntil, ok := row["lease_until"].(interface{ Time() time.Time })
	if ok {
		s.False(leaseUntil.Time().After(time.Now().Add(time.Minute)), "release must expire the lease for immediate retry")
	}
}

func TestInboxSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(InboxTestSuite))
}
