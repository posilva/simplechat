package testutil

import (
	"testing"

	"github.com/posilva/simplechat/internal/core/domain"
	"github.com/posilva/simplechat/internal/core/ports"
	uuid "github.com/segmentio/ksuid"
)

const (
	DynamoDBLocalTableName string = "local-dev-dev-simplechat"
	RabbitMQLocalURL       string = "amqp://guest:guest@localhost:5672/"
)

func Name(t *testing.T) string {
	return t.Name()
}

func NewID() string {
	return uuid.New().String()
}

func NewUnique(prefix string) string {
	return prefix + NewID()
}

type TestEndpoint struct {
	id       string
	room     string
	receiver ports.Receiver
}

func NewTestEndpoint(id string, room string, receiver ports.Receiver) *TestEndpoint {
	return &TestEndpoint{
		id,
		room,
		receiver,
	}
}
func (ep *TestEndpoint) ID() string {
	return ep.id
}
func (ep *TestEndpoint) Room() string {
	return ep.room
}
func (ep *TestEndpoint) Receive(m domain.ModeratedMessage) {
	ep.receiver.Receive(m)
}
func (ep *TestEndpoint) Recover() {
	ep.receiver.Recover()
}

type TestReceiver struct {
	ch chan domain.ModeratedMessage
	f  func()
}

func NewTestReceiverWithFunc(f func()) *TestReceiver {
	return &TestReceiver{
		ch: make(chan domain.ModeratedMessage, 1),
		f:  f,
	}
}

func NewTestReceiver() *TestReceiver {
	return &TestReceiver{
		ch: make(chan domain.ModeratedMessage, 1),
		f:  func() {},
	}
}
func (r *TestReceiver) Receive(m domain.ModeratedMessage) {
	r.f()
	r.ch <- m
}
func (r *TestReceiver) Recover() {
	close(r.ch)
}
func (r *TestReceiver) Channel() chan domain.ModeratedMessage {
	return r.ch
}
