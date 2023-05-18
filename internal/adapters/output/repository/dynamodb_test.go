package repository

import (
	"testing"
	"time"

	"github.com/posilva/simplechat/internal/adapters/output/logging"
	"github.com/posilva/simplechat/internal/core/domain"
	testutils "github.com/posilva/simplechat/internal/testutil"
	"github.com/stretchr/testify/assert"
)

const (
	localTableName string = "local-dev-dev-simplechat"
)

func TestNewDynamoDBRepository(t *testing.T) {
	r := newDynamoDBRepository(t)

	expType := &DynamoDBRepository{}
	assert.IsType(t, expType, r)
}

func TestDynamoDBRepository_Store(t *testing.T) {
	r := newDynamoDBRepository(t)

	id1 := testutils.NewID()
	topic := testutils.NewUnique(testutils.Name(t))
	payload := "TestDynamoDBRepository_Store Message"

	m := domain.ModeratedMessage{
		Message: domain.Message{
			From:    id1,
			To:      topic,
			Payload: payload,
		},
		Level:           0,
		FilteredPayload: payload,
	}
	err := r.Store(m)
	assert.NoError(t, err)
}

func TestDynamoDBRepository_Fetch(t *testing.T) {
	r := newDynamoDBRepository(t)
	id1 := testutils.NewID()

	topic := testutils.NewUnique(testutils.Name(t))
	payload := "TestDynamoDBRepository_Fetch Message"

	m := domain.ModeratedMessage{
		Message: domain.Message{
			From:    id1,
			To:      topic,
			Payload: payload,
		},
		Level:           0,
		FilteredPayload: payload,
	}
	err := r.Store(m)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)
	msgs, err := r.Fetch(topic, 3*time.Second)
	assert.NoError(t, err)

	assert.Len(t, msgs, 1)
}

func newDynamoDBRepository(t *testing.T) *DynamoDBRepository {
	log := logging.NewSimpleLogger()
	r, err := NewDynamoDBRepository(DefaultLocalAWSClientConfig(), localTableName, log)
	assert.NoError(t, err)
	return r
}
