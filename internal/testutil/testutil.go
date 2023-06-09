// Package testutil is used to share test utilities
package testutil

import (
	"testing"

	"github.com/posilva/simplechat/internal/core/domain"
	"github.com/posilva/simplechat/internal/core/ports"
	uuid "github.com/segmentio/ksuid"
)

const (
	// DynamoDBLocalTableName defines the DynamoDB table name for local development with LocalStack
	DynamoDBLocalTableName string = "local-dev-dev-simplechat"
	// RabbitMQLocalURL defines the local url to connect to Rabbit MQ
	RabbitMQLocalURL string = "amqp://guest:guest@localhost:5672/"
	// RabbitMQLocalURLSSL defines the local url to connect to Rabbit MQ using SSL
	RabbitMQLocalURLSSL string = "amqps://guest:guest@localhost:5671/"
)

// Name returns the name of the test
func Name(t *testing.T) string {
	return t.Name()
}

// NewID returns an ID for tests using kuid package
func NewID() string {
	return uuid.New().String()
}

// NewUnique appends to a string a UUID to allow for uniqueness
func NewUnique(prefix string) string {
	return prefix + NewID()
}

// TestEndpoint is used for testing when am Endpoint is needed
type TestEndpoint struct {
	receiver ports.Receiver
	id       string
	room     string
}

// NewTestEndpointWithReceiver creates a test endpoint with a passed receiver
func NewTestEndpointWithReceiver(id string, room string, r ports.Receiver) *TestEndpoint {
	return &TestEndpoint{
		receiver: r,
		id:       id,
		room:     room,
	}
}

// NewTestEndpoint creates a new test endpoint with default test receiver
func NewTestEndpoint(id string, room string) *TestEndpoint {
	return &TestEndpoint{
		receiver: NewTestReceiver(),
		id:       id,
		room:     room,
	}
}

// Channel returns the channel from the internal receiver
func (ep *TestEndpoint) Channel() chan domain.Notication {
	r := ep.receiver
	return r.(*TestReceiver).ch
}

// ID returns the id of the endpoint
func (ep *TestEndpoint) ID() string {
	return ep.id
}

// Room returns the room id of the endpoint
func (ep *TestEndpoint) Room() string {
	return ep.room
}

// Receive implements the Receiver interface
func (ep *TestEndpoint) Receive(m domain.Notication) {
	ep.receiver.Receive(m)
}

// Recover implements the Receiver interface
func (ep *TestEndpoint) Recover() {
	ep.receiver.Recover()
}

// TestReceiver represents a Receiver interface used for Tests
type TestReceiver struct {
	ch chan domain.Notication
	f  func()
}

// NewTestReceiverWithFunc creates a test receiver that can intercept the receive with a function
// used to generate panic
func NewTestReceiverWithFunc(f func()) *TestReceiver {
	return &TestReceiver{
		ch: make(chan domain.Notication, 1),
		f:  f,
	}
}

// NewTestReceiver creates a test receiver with a channel for tests
func NewTestReceiver() *TestReceiver {
	return &TestReceiver{
		ch: make(chan domain.Notication, 1),
		f:  func() {},
	}
}

// Receive is called every time a message should be delivered
func (r *TestReceiver) Receive(m domain.Notication) {
	r.f()
	r.ch <- m
}

// Recover is called when there is an panic in go routine receiving the message
func (r *TestReceiver) Recover() {
	close(r.ch)
}

// Channel returns the channel used internally for communication
func (r *TestReceiver) Channel() chan domain.Notication {
	return r.ch
}

// NewSimpleEndpoint returns an endpoint with a default receiver and a new generated ID
func NewSimpleEndpoint(room string) *TestEndpoint {
	id := NewID()
	ep := NewTestEndpoint(id, room)
	return ep

}
