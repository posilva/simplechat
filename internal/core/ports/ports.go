// Package ports should define the interfaces to interact
// with services
package ports

import (
	"time"

	"github.com/posilva/simplechat/internal/core/domain"
)

// ChatService defines the actions that a chat service provides
type ChatService interface {
	Register(ep Endpoint) error
	DeRegister(ep Endpoint) error
	Send(domain.Message) error
	History(dst string, since time.Duration) ([]*domain.ModeratedMessage, error)
}

// Repository defines the interface to handle with
// with the storage layer of the chat messages
type Repository interface {
	Store(m domain.ModeratedMessage) error
	Fetch(key string, since time.Duration) ([]*domain.ModeratedMessage, error)
}

// Notifier defines the interface to handle with
// the notifications of the chat messages
type Notifier interface {
	Subscribe(ep Endpoint) error
	Unsubscribe(ep Endpoint) error
	Broadcast(m domain.ModeratedMessage) error
}

// Moderator defines the interface to handle with
// the moderation of the chat messages
type Moderator interface {
	Check(m domain.Message) (*domain.ModeratedMessage, error)
}

// Registry defines the interface of a endpoint registries
type Registry interface {
	// TODO: this may a generic Notication message as we may notify more than
	// Chat messages (Presence status updates, joins, leaves)
	Notify(m domain.ModeratedMessage)
	Register(ep Endpoint) error
	DeRegister(ep Endpoint) error
}

// Receiver define the interface to receive messages
type Receiver interface {
	// TODO: this may a generic Notication message as we may notify more than
	// Chat messages (Presence status updates, joins, leaves)
	Receive(domain.ModeratedMessage)
	Recover()
}

// Endpoint defines an interface of an endpoint that can receive a message
type Endpoint interface {
	Receiver
	ID() string
	Room() string
}

// Presence defines an interface to a presence component
type Presence interface {
	Join(ep Endpoint) error
	Leave(ep Endpoint) error
	// TODO: may introduce later the Participant
	Presents(room string) (v map[string]string, err error)
	IsPresent(ep Endpoint) (bool, error)
}
