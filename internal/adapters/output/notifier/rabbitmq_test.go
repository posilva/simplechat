package notifier

import (
	"crypto/tls"
	"testing"

	"github.com/posilva/simplechat/internal/adapters/output/registry"
	"github.com/posilva/simplechat/internal/core/domain"
	testutils "github.com/posilva/simplechat/internal/testutil"

	assert "github.com/stretchr/testify/assert"
)

const (
	localURL    string = "amqp://guest:guest@localhost:5672/"
	localURLSSL string = "amqps://guest:guest@localhost:5671/"
)

func TestNewRabbitMQNotifierWithLocal(t *testing.T) {
	got, err := NewRabbitMQNotifierWithLocal(localURL, registry.NewInMemoryRegistry())

	expectedType := &RabbitMQNotifier{}
	assert.NoError(t, err, "expected not an error")
	assert.NotNil(t, got, "expected notifier to be not nil")
	assert.IsType(t, expectedType, got, "expected RabbitMQNotifier ")
	assert.NotNil(t, got.conn, "connection should be non nil")
	assert.NotNil(t, got.ch, "channel should be non nil")
}

func TestNewRabbitMQNotifierWithTLS(t *testing.T) {

	got, err := NewRabbitMQNotifierWithTLS(localURLSSL, &tls.Config{InsecureSkipVerify: true}, registry.NewInMemoryRegistry())

	expectedType := &RabbitMQNotifier{}
	assert.NoError(t, err, "expected not an error")
	assert.NotNil(t, got, "expected notifier to be not nil")
	assert.IsType(t, expectedType, got, "expected RabbitMQNotifier ")
	assert.NotNil(t, got.conn, "connection should be non nil")
	assert.NotNil(t, got.ch, "channel should be non nil")
}

func TestRabbitMQNotifier_Broadcast_ReceiveWihtPanic(t *testing.T) {
	got, err := NewRabbitMQNotifierWithLocal(localURL, registry.NewInMemoryRegistry())

	assert.NoError(t, err, "expected not an error")

	topic := testutils.NewUnique(testutils.Name(t))

	r := testutils.NewTestReceiverWithFunc(func() { panic("receive with panic") })

	ep1 := testutils.NewTestEndpoint(testutils.NewID(), topic, r)
	ep2 := testutils.NewTestEndpoint(testutils.NewID(), topic, r)

	err = got.Subscribe(ep1)
	assert.NoError(t, err, "expected not an error")

	err = got.Subscribe(ep2)
	assert.NoError(t, err, "expected not an error")
	payload := "TestRabbitMQNotifier_Broadcast Message"
	m := domain.ModeratedMessage{
		Message: domain.Message{
			From:    ep1.ID(),
			To:      topic,
			Payload: payload,
		},
		Level:           0,
		FilteredPayload: payload,
	}
	err = got.Broadcast(m)

	assert.NoError(t, err, "expected not an error")
	m1 := <-r.Channel()
	assert.Equal(t, domain.ModeratedMessage{}, m1)
}

func TestRabbitMQNotifier_Broadcast(t *testing.T) {
	got, err := NewRabbitMQNotifierWithLocal(localURL, registry.NewInMemoryRegistry())

	assert.NoError(t, err, "expected not an error")
	topic := testutils.NewUnique(testutils.Name(t))

	r := testutils.NewTestReceiver()

	ep1 := testutils.NewTestEndpoint(testutils.NewID(), topic, r)
	ep2 := testutils.NewTestEndpoint(testutils.NewID(), topic, r)

	err = got.Subscribe(ep1)
	assert.NoError(t, err, "expected not an error")

	err = got.Subscribe(ep2)
	assert.NoError(t, err, "expected not an error")

	payload := "TestRabbitMQNotifier_Broadcast Message"
	m := domain.ModeratedMessage{
		Message: domain.Message{
			From:    ep1.ID(),
			To:      topic,
			Payload: payload,
		},
		Level:           0,
		FilteredPayload: payload,
	}

	err = got.Broadcast(m)

	assert.NoError(t, err, "expected not an error")
	m1 := <-r.Channel()
	assert.Equal(t, m, m1)
}

func TestRabbitMQNotifier_Subscriber(t *testing.T) {
	got, err := NewRabbitMQNotifierWithLocal(localURL, registry.NewInMemoryRegistry())

	assert.NoError(t, err, "expected not an error")
	topic := testutils.NewUnique(testutils.Name(t))

	r := testutils.NewTestReceiver()
	ep1 := testutils.NewTestEndpoint(testutils.NewID(), topic, r)
	ep2 := testutils.NewTestEndpoint(testutils.NewID(), topic, r)
	ep3 := testutils.NewTestEndpoint(testutils.NewID(), topic, r)
	ep4 := testutils.NewTestEndpoint(testutils.NewID(), topic, r)

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

	got, err := NewRabbitMQNotifierWithLocal(localURL, registry.NewInMemoryRegistry())

	assert.NoError(t, err, "expected not an error")
	topic := testutils.NewUnique(testutils.Name(t))

	r := testutils.NewTestReceiver()
	ep1 := testutils.NewTestEndpoint(testutils.NewID(), topic, r)

	err = got.Subscribe(ep1)
	assert.NoError(t, err, "expected not an error")

	err = got.Unsubscribe(ep1)
	assert.NoError(t, err, "expected not an error")
}
