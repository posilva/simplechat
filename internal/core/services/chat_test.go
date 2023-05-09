package services

import (
	"testing"
	"time"

	"github.com/posilva/simplechat/internal/adapters/output/moderator"
	"github.com/posilva/simplechat/internal/adapters/output/notifier"
	"github.com/posilva/simplechat/internal/adapters/output/repository"
	"github.com/posilva/simplechat/internal/core/domain"
	uuid "github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

const (
	// TODO make a public constant
	localTableName string = "local-dev-dev-simplechat"
	localURL       string = "amqp://guest:guest@localhost:5672/"
)

// TODO: come up with something to share helpers for tests
type testReceiver struct {
	ch chan domain.ModeratedMessage
	f  func()
}

func newTestReceiver(ch chan domain.ModeratedMessage, f func()) *testReceiver {
	return &testReceiver{
		ch,
		f,
	}
}
func (r *testReceiver) Receive(m domain.ModeratedMessage) {
	r.f()
	r.ch <- m
}
func (r *testReceiver) Recover() {
	close(r.ch)
}
func TestNewChatService(t *testing.T) {
	expType := &ChatService{}

	r, err := repository.NewDynamoDBRepository(repository.DefaultiLocalAWSClientConfig(), localTableName)
	assert.NoError(t, err)

	n, err := notifier.NewRabbitMQNotifierWithLocal(localURL)
	assert.NoError(t, err)

	m := moderator.NewIgnoreModerator()
	cs := NewChatService(r, n, m)

	assert.NotNil(t, cs)
	assert.IsType(t, expType, cs)
}

func TestChatService_Send(t *testing.T) {
	r, err := repository.NewDynamoDBRepository(repository.DefaultiLocalAWSClientConfig(), localTableName)
	assert.NoError(t, err)

	n, err := notifier.NewRabbitMQNotifierWithLocal(localURL)
	assert.NoError(t, err)

	m := moderator.NewIgnoreModerator()
	cs := NewChatService(r, n, m)

	assert.NotNil(t, cs)

	topic := "TestChatService_Send" + "_" + uuid.New().String()
	payload := "TestChatService_Send Message"

	c := make(chan domain.ModeratedMessage, 1)
	rc := newTestReceiver(c, func() {})

	id1 := uuid.New().String()
	err = n.Register(id1, topic, rc)
	assert.NoError(t, err)

	id2 := uuid.New().String()
	err = n.Register(id2, topic, rc)
	assert.NoError(t, err)
	msg := domain.Message{
		From:    id1,
		To:      topic,
		Payload: payload,
	}
	err = cs.Send(msg)
	assert.NoError(t, err)

	m1 := <-c
	assert.Equal(t, msg, m1.Message)
}

func TestChatService_History(t *testing.T) {
	r, err := repository.NewDynamoDBRepository(repository.DefaultiLocalAWSClientConfig(), localTableName)
	assert.NoError(t, err)

	n, err := notifier.NewRabbitMQNotifierWithLocal(localURL)
	assert.NoError(t, err)

	m := moderator.NewIgnoreModerator()
	cs := NewChatService(r, n, m)

	assert.NotNil(t, cs)

	id1 := uuid.New().String()

	topic := "TestChatService_History" + "_" + uuid.New().String()
	payload := "TestChatService_History Message"

	msg := domain.Message{
		From:    id1,
		To:      topic,
		Payload: payload,
	}
	err = cs.Send(msg)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)
	msgs, err := cs.History(topic, 3*time.Second)

	assert.NoError(t, err)
	assert.Len(t, msgs, 1)
}
