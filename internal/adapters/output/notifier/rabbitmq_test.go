package notifier

import (
	"crypto/tls"
	"testing"

	"github.com/posilva/simplechat/internal/core/domain"
	uuid "github.com/segmentio/ksuid"

	assert "github.com/stretchr/testify/assert"
)

const (
	localURL    string = "amqp://guest:guest@localhost:5672/"
	localURLSSL string = "amqps://guest:guest@localhost:5671/"
)

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

func TestNewRabbitMQNotifierWithLocal(t *testing.T) {
	got, err := NewRabbitMQNotifierWithLocal(localURL)

	expectedType := &RabbitMQNotifier{}
	assert.NoError(t, err, "expected not an error")
	assert.NotNil(t, got, "expected notifier to be not nil")
	assert.IsType(t, expectedType, got, "expected RabbitMQNotifier ")
	assert.NotNil(t, got.conn, "connection should be non nil")
	assert.NotNil(t, got.ch, "channel should be non nil")
}

func TestNewRabbitMQNotifierWithTLS(t *testing.T) {

	got, err := NewRabbitMQNotifierWithTLS(localURLSSL, &tls.Config{InsecureSkipVerify: true})

	expectedType := &RabbitMQNotifier{}
	assert.NoError(t, err, "expected not an error")
	assert.NotNil(t, got, "expected notifier to be not nil")
	assert.IsType(t, expectedType, got, "expected RabbitMQNotifier ")
	assert.NotNil(t, got.conn, "connection should be non nil")
	assert.NotNil(t, got.ch, "channel should be non nil")

}

func TestRabbitMQNotifier_Broadcast_ReceiveWihtPanic(t *testing.T) {
	got, err := NewRabbitMQNotifierWithLocal(localURL)

	assert.NoError(t, err, "expected not an error")
	topic := "TestRabbitMQNotifier_Broadcast_ReceiveWihtPanic"

	id1 := uuid.New().String()
	id2 := uuid.New().String()
	c := make(chan domain.ModeratedMessage, 1)
	r := newTestReceiver(c, func() { panic("receive with panic") })

	err = got.Register(id1, topic, r)
	assert.NoError(t, err, "expected not an error")

	err = got.Register(id2, topic, r)
	assert.NoError(t, err, "expected not an error")

	payload := "TestRabbitMQNotifier_Broadcast Message"
	m := domain.ModeratedMessage{
		Message: domain.Message{
			From:    id1,
			To:      topic,
			Payload: payload,
		},
		Level:           0,
		FilteredPayload: payload,
	}
	err = got.Broadcast(m)

	assert.NoError(t, err, "expected not an error")
	m1 := <-c
	assert.Equal(t, domain.ModeratedMessage{}, m1)
}

func TestRabbitMQNotifier_Broadcast(t *testing.T) {
	got, err := NewRabbitMQNotifierWithLocal(localURL)

	assert.NoError(t, err, "expected not an error")
	topic := "TestRabbitMQNotifier_Broadcast"

	id1 := uuid.New().String()
	id2 := uuid.New().String()
	c := make(chan domain.ModeratedMessage, 1)
	r := newTestReceiver(c, func() {})

	err = got.Register(id1, topic, r)
	assert.NoError(t, err, "expected not an error")

	err = got.Register(id2, topic, r)
	assert.NoError(t, err, "expected not an error")

	payload := "TestRabbitMQNotifier_Broadcast Message"
	m := domain.ModeratedMessage{
		Message: domain.Message{
			From:    id1,
			To:      topic,
			Payload: payload,
		},
		Level:           0,
		FilteredPayload: payload,
	}
	err = got.Broadcast(m)

	assert.NoError(t, err, "expected not an error")
	m1 := <-c
	assert.Equal(t, m, m1)
}

func TestRabbitMQNotifier_Register(t *testing.T) {
	got, err := NewRabbitMQNotifierWithLocal(localURL)

	assert.NoError(t, err, "expected not an error")
	topic := "TestRabbitMQNotifier_Register"

	id1 := uuid.New().String()
	id2 := uuid.New().String()
	id3 := uuid.New().String()
	id4 := uuid.New().String()

	r := &testReceiver{}
	err = got.Register(id1, topic, r)
	assert.NoError(t, err, "expected not an error")

	err = got.Register(id2, topic, r)
	assert.NoError(t, err, "expected not an error")

	err = got.Register(id3, topic, r)
	assert.NoError(t, err, "expected not an error")
	err = got.Register(id4, topic, r)
	assert.NoError(t, err, "expected not an error")
}

func TestRabbitMQNotifier_DeRegister(t *testing.T) {

	got, err := NewRabbitMQNotifierWithLocal(localURL)

	assert.NoError(t, err, "expected not an error")

	id1 := uuid.New().String()
	r := &testReceiver{}
	topic := "TestRabbitMQNotifier_DeRegister"
	err = got.Register(id1, topic, r)
	assert.NoError(t, err, "expected not an error")
	err = got.DeRegister(id1)
	assert.NoError(t, err, "expected not an error")
}
