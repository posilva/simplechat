package services

import (
	"testing"
	"time"

	"github.com/posilva/simplechat/internal/adapters/output/logging"
	"github.com/posilva/simplechat/internal/adapters/output/moderator"
	"github.com/posilva/simplechat/internal/adapters/output/notifier"
	"github.com/posilva/simplechat/internal/adapters/output/notifier/codecs"
	"github.com/posilva/simplechat/internal/adapters/output/presence"
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

	room := testutils.NewUnique(testutils.Name(t))

	payload := "TestChatService_Send Message"

	ep1 := testutils.NewSimpleEndpoint(room)
	ep2 := testutils.NewSimpleEndpoint(room)

	err := cs.Register(ep1)
	assert.NoError(t, err)
	err = cs.Register(ep2)
	assert.NoError(t, err)

	msg := domain.Message{
		From:    ep1.ID(),
		To:      room,
		Payload: payload,
	}
	err = cs.Send(msg)
	assert.NoError(t, err)

	m1 := <-ep2.Channel()
	_ = m1.Payload.(domain.PresenceUpdate)

	m1 = <-ep2.Channel()
	_ = m1.Payload.(domain.PresenceUpdate)

	m2 := <-ep2.Channel()
	m := m2.Payload.(domain.ModeratedMessage)
	assert.Equal(t, msg, m.Message)
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
	room := testutils.NewUnique(testutils.Name(t))
	ep1 := testutils.NewSimpleEndpoint(room)

	err := cs.Register(ep1)
	assert.NoError(t, err)
}

func TestChatService_DeRegister(t *testing.T) {

	cs := newChatService(t)

	room := testutils.NewUnique(testutils.Name(t))
	ep1 := testutils.NewSimpleEndpoint(room)

	err := cs.DeRegister(ep1)
	assert.NoError(t, err)
}

func newChatService(t *testing.T) *ChatService {
	log := logging.NewSimpleLogger()
	r, err := repository.NewDynamoDBRepository(repository.DefaultLocalAWSClientConfig(), testutils.DynamoDBLocalTableName, log)
	assert.NoError(t, err)

	reg := registry.NewInMemoryRegistry(log)

	n, err := notifier.NewRabbitMQNotifierWithLocal[*codecs.JSONNotifierCodec](testutils.RabbitMQLocalURL, reg, log)
	assert.NoError(t, err)

	m := moderator.NewIgnoreModerator()
	ps, err := presence.NewRedisPresence(presence.DefaultLocalOpts(), n, log)

	cs := NewChatService(r, n, m, ps, log)
	assert.NoError(t, err)

	assert.NotNil(t, cs)

	return cs
}
