package services

import (
	"testing"
	"time"

	"github.com/posilva/simplechat/internal/adapters/output/moderator"
	"github.com/posilva/simplechat/internal/adapters/output/notifier"
	"github.com/posilva/simplechat/internal/adapters/output/registry"
	"github.com/posilva/simplechat/internal/adapters/output/repository"
	"github.com/posilva/simplechat/internal/core/domain"
	testutils "github.com/posilva/simplechat/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNewChatService(t *testing.T) {
	expType := &ChatService{}
	cs := newChatService(t)
	assert.IsType(t, expType, cs)
}

func TestChatService_Send(t *testing.T) {
	cs := newChatService(t)

	topic := testutils.NewUnique(testutils.Name(t))

	payload := "TestChatService_Send Message"

	rc := testutils.NewTestReceiver()

	ep1 := testutils.NewTestEndpoint(testutils.NewID(), topic, rc)
	ep2 := testutils.NewTestEndpoint(testutils.NewID(), topic, rc)

	err := cs.Register(ep1)
	assert.NoError(t, err)
	err = cs.Register(ep2)
	assert.NoError(t, err)

	msg := domain.Message{
		From:    ep1.ID(),
		To:      topic,
		Payload: payload,
	}
	err = cs.Send(msg)
	assert.NoError(t, err)

	m1 := <-rc.Channel()
	assert.Equal(t, msg, m1.Message)
}

func TestChatService_History(t *testing.T) {
	cs := newChatService(t)

	id1 := testutils.NewID()

	topic := testutils.NewUnique(testutils.Name(t))

	payload := "TestChatService_History Message"

	msg := domain.Message{
		From:    id1,
		To:      topic,
		Payload: payload,
	}
	err := cs.Send(msg)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)
	msgs, err := cs.History(topic, 3*time.Second)

	assert.NoError(t, err)
	assert.Len(t, msgs, 1)
}
func TestChatService_Register(t *testing.T) {
	cs := newChatService(t)

	topic := testutils.NewUnique(testutils.Name(t))

	rc := testutils.NewTestReceiver()
	ep1 := testutils.NewTestEndpoint(testutils.NewID(), topic, rc)

	err := cs.Register(ep1)
	assert.NoError(t, err)

}

func TestChatService_DeRegister(t *testing.T) {

	cs := newChatService(t)

	topic := testutils.NewUnique(testutils.Name(t))

	rc := testutils.NewTestReceiver()
	ep1 := testutils.NewTestEndpoint(testutils.NewID(), topic, rc)

	err := cs.DeRegister(ep1)
	assert.NoError(t, err)
}

func newChatService(t *testing.T) *ChatService {
	r, err := repository.NewDynamoDBRepository(repository.DefaultiLocalAWSClientConfig(), testutils.DynamoDBLocalTableName)
	assert.NoError(t, err)

	reg := registry.NewInMemoryRegistry()

	n, err := notifier.NewRabbitMQNotifierWithLocal(testutils.RabbitMQLocalURL, reg)
	assert.NoError(t, err)

	m := moderator.NewIgnoreModerator()

	cs := NewChatService(r, n, m)

	assert.NotNil(t, cs)

	return cs
}
