package notifier

import (
	"crypto/tls"
	"testing"

	"github.com/posilva/simplechat/internal/adapters/output/logging"
	"github.com/posilva/simplechat/internal/adapters/output/notifier/codecs"
	"github.com/posilva/simplechat/internal/adapters/output/registry"
	"github.com/posilva/simplechat/internal/core/domain"
	"github.com/posilva/simplechat/internal/core/ports"
	testutils "github.com/posilva/simplechat/internal/testutil"

	assert "github.com/stretchr/testify/assert"
)

func TestNewRabbitMQNotifierWithLocal(t *testing.T) {
	got, err := newRabbitMQNotifier[*codecs.JSONNotifierCodec](false)

	expectedType := &RabbitMQNotifier[*codecs.JSONNotifierCodec]{}
	assert.NoError(t, err, "expected not an error")
	assert.NotNil(t, got, "expected notifier to be not nil")
	assert.IsType(t, expectedType, got, "expected RabbitMQNotifier ")
	assert.NotNil(t, got.conn, "connection should be non nil")
	assert.NotNil(t, got.ch, "channel should be non nil")
}

func TestNewRabbitMQNotifierWithTLS(t *testing.T) {

	got, err := newRabbitMQNotifier[*codecs.JSONNotifierCodec](true)

	expectedType := &RabbitMQNotifier[*codecs.JSONNotifierCodec]{}
	assert.NoError(t, err, "expected not an error")
	assert.NotNil(t, got, "expected notifier to be not nil")
	assert.IsType(t, expectedType, got, "expected RabbitMQNotifier ")
	assert.NotNil(t, got.conn, "connection should be non nil")
	assert.NotNil(t, got.ch, "channel should be non nil")
}

func TestRabbitMQNotifier_Broadcast_ReceiveWithPanic(t *testing.T) {
	got, err := newRabbitMQNotifier[*codecs.JSONNotifierCodec](false)

	assert.NoError(t, err, "expected not an error")

	room := testutils.NewUnique(testutils.Name(t))

	r := testutils.NewTestReceiverWithFunc(func() { panic("receive with panic") })

	ep1 := testutils.NewTestEndpointWithReceiver(testutils.NewID(), room, r)

	err = got.Subscribe(ep1)
	assert.NoError(t, err, "expected not an error")

	payload := "TestRabbitMQNotifier_Broadcast Message"

	m := domain.Notication{
		Kind: domain.ModeratedMessageKind,
		To:   room,
		Payload: domain.ModeratedMessage{
			Message: domain.Message{
				From:    ep1.ID(),
				To:      room,
				Payload: payload,
			},
			Level:           0,
			FilteredPayload: payload,
		},
	}
	err = got.Broadcast(m)

	assert.NoError(t, err)
	t.Logf("Checking channel ")
	m1 := <-ep1.Channel()
	t.Logf("Checking channel emd")
	assert.Nil(t, m1.Payload)
}

func TestRabbitMQNotifier_Broadcast(t *testing.T) {
	got, err := newRabbitMQNotifier[*codecs.JSONNotifierCodec](false)

	assert.NoError(t, err, "expected not an error")
	topic := testutils.NewUnique(testutils.Name(t))

	ep1 := testutils.NewSimpleEndpoint(topic)
	ep2 := testutils.NewSimpleEndpoint(topic)

	err = got.Subscribe(ep1)
	assert.NoError(t, err, "expected not an error")

	err = got.Subscribe(ep2)
	assert.NoError(t, err, "expected not an error")

	payload := "TestRabbitMQNotifier_Broadcast Message"
	m := domain.Notication{
		To:   topic,
		Kind: domain.ModeratedMessageKind,
		Payload: domain.ModeratedMessage{
			Message: domain.Message{
				From:    ep1.ID(),
				To:      topic,
				Payload: payload,
			},
			Level:           0,
			FilteredPayload: payload,
		},
	}

	err = got.Broadcast(m)
	assert.NoError(t, err)

	m1 := <-ep1.Channel()
	assert.Equal(t, m1.Payload, m.Payload)
	m2 := <-ep2.Channel()
	assert.Equal(t, m2.Payload, m.Payload)
}

func TestRabbitMQNotifier_Subscriber(t *testing.T) {
	got, err := newRabbitMQNotifier[*codecs.JSONNotifierCodec](false)

	assert.NoError(t, err, "expected not an error")
	topic := testutils.NewUnique(testutils.Name(t))

	ep1 := testutils.NewSimpleEndpoint(topic)
	ep2 := testutils.NewSimpleEndpoint(topic)
	ep3 := testutils.NewSimpleEndpoint(topic)
	ep4 := testutils.NewSimpleEndpoint(topic)

	err = got.Subscribe(ep1)
	assert.NoError(t, err, "expected not an error")

	err = got.Subscribe(ep2)
	assert.NoError(t, err, "expected not an error")

	err = got.Subscribe(ep3)
	assert.NoError(t, err, "expected not an error")

	err = got.Subscribe(ep4)
	assert.NoError(t, err, "expected not an error")
}

func TestRabbitMQNotifier_Unsubscribe(t *testing.T) {
	got, err := newRabbitMQNotifier[*codecs.JSONNotifierCodec](false)

	assert.NoError(t, err, "expected not an error")
	topic := testutils.NewUnique(testutils.Name(t))

	ep1 := testutils.NewSimpleEndpoint(topic)

	err = got.Subscribe(ep1)
	assert.NoError(t, err, "expected not an error")

	err = got.Unsubscribe(ep1)
	assert.NoError(t, err, "expected not an error")
}

func newRabbitMQNotifier[T ports.NotifierCodec](secure bool) (*RabbitMQNotifier[T], error) {
	log := logging.NewSimpleLogger()
	reg := registry.NewInMemoryRegistry(log)
	if secure {
		return NewRabbitMQNotifierWithTLS[T](testutils.RabbitMQLocalURLSSL, &tls.Config{InsecureSkipVerify: true}, reg, log)
	}
	return NewRabbitMQNotifierWithLocal[T](testutils.RabbitMQLocalURL, reg, log)
}
