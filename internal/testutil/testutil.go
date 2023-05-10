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
	// RabbitMQLocalURL defines the local url to connect to Rabbit MQ running in docker
	RabbitMQLocalURL string = "amqp://guest:guest@localhost:5672/"
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

// NewTestEndpoint creates a new test endpoint
func NewTestEndpoint(id string, room string, receiver ports.Receiver) *TestEndpoint {
	return &TestEndpoint{
		receiver,
		id,
		room,
	}
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
func (ep *TestEndpoint) Receive(m domain.ModeratedMessage) {
	ep.receiver.Receive(m)
}

// Recover implements the Receiver interface
func (ep *TestEndpoint) Recover() {
	ep.receiver.Recover()
}

// TestReceiver represents a Receiver interface used for Tests
type TestReceiver struct {
	ch chan domain.ModeratedMessage
	f  func()
}

// NewTestReceiverWithFunc creates a test receiver that can intercept the receive with a function
// used to generate panic
func NewTestReceiverWithFunc(f func()) *TestReceiver {
	return &TestReceiver{
		ch: make(chan domain.ModeratedMessage, 1),
		f:  f,
	}
}

// NewTestReceiver creates a test receiver with a channel for tests
func NewTestReceiver() *TestReceiver {
	return &TestReceiver{
		ch: make(chan domain.ModeratedMessage, 1),
		f:  func() {},
	}
}

// Receive is called every time a message should be delivered
func (r *TestReceiver) Receive(m domain.ModeratedMessage) {
	r.f()
	r.ch <- m
}

// Recover is called when there is an panic in go routine receiving the message
func (r *TestReceiver) Recover() {
	close(r.ch)
}

// Channel returns the channel used internally for communication
func (r *TestReceiver) Channel() chan domain.ModeratedMessage {
	return r.ch
}
