package handler

import (
	"github.com/posilva/simplechat/internal/core/domain"
	"github.com/posilva/simplechat/internal/core/ports"
)

type clientEndpoint struct {
	receiver ports.Receiver
	id       string
	room     string
}

func newclientEndpoint(id string, room string, receiver ports.Receiver) *clientEndpoint {
	return &clientEndpoint{
		receiver,
		id,
		room,
	}
}

// ID returns the id of the endpoint
func (ep *clientEndpoint) ID() string {
	return ep.id
}

// Room returns the room id of the endpoint
func (ep *clientEndpoint) Room() string {
	return ep.room
}

// Receive implements the Receiver interface
func (ep *clientEndpoint) Receive(m domain.ModeratedMessage) {
	ep.receiver.Receive(m)
}

// Recover implements the Receiver interface
func (ep *clientEndpoint) Recover() {
	ep.receiver.Recover()
}

// clientReceiver represents a Receiver interface used for Tests
type clientReceiver struct {
	ch chan domain.ModeratedMessage
	f  func()
}

// newclientReceiver creates a test receiver with a channel for tests
func newclientReceiver() *clientReceiver {
	return &clientReceiver{
		ch: make(chan domain.ModeratedMessage, 1),
		f:  func() {},
	}
}

// Receive is called every time a message should be delivered
func (r *clientReceiver) Receive(m domain.ModeratedMessage) {
	r.f()
	r.ch <- m
}

// Recover is called when there is an panic in go routine receiving the message
func (r *clientReceiver) Recover() {
	close(r.ch)
}

// Channel returns the channel used internally for communication
func (r *clientReceiver) Channel() chan domain.ModeratedMessage {
	return r.ch
}
