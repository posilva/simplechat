// package ports should define the interfaces to interact
// with services
package ports

import (
	"time"

	"github.com/posilva/simplechat/internal/core/domain"
)

// ChatService defines the actions that a chat service provides
type ChatService interface {
	Send(domain.Message) error
	History(since time.Duration) ([]*domain.ModeratedMessage, error)
}

// Repository defines the interface to handle with
// with the storage layer of the chat messages
type Repository interface {
	Store(m domain.ModeratedMessage) error
	Fetch(since time.Duration) ([]*domain.ModeratedMessage, error)
}

// Notifier defines the interface to handle with
// the notifications of the chat messages
type Notifier interface {
	Registry
	Broadcast(m domain.ModeratedMessage) error
}

// Moderator defines the interface to handle with
// the moderation of the chat messages
type Moderator interface {
	Check(m domain.Message) (*domain.ModeratedMessage, error)
}

// Registry defines the interface of a EndPoint registry
type Registry interface {
	Register(id string, topic string) error
	DeRegister(id string) error
}
