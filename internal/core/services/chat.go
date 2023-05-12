// Package services implements services
package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/posilva/simplechat/internal/core/domain"
	"github.com/posilva/simplechat/internal/core/ports"
)

// ChatService implements the chat logic
type ChatService struct {
	repository ports.Repository
	notifier   ports.Notifier
	moderator  ports.Moderator
}

// NewChatService creates a new instance of a chat service using
// existing repository, notifier and moderator providers
func NewChatService(repo ports.Repository, notif ports.Notifier, mod ports.Moderator) *ChatService {
	return &ChatService{
		repository: repo,
		notifier:   notif,
		moderator:  mod,
	}
}

// Register registers an Endpoint in the chat service
func (c *ChatService) Register(ep ports.Endpoint) error {
	return c.notifier.Subscribe(ep)
}

// UnRegister unregisters an Endpoint in the chat service
func (c *ChatService) UnRegister(ep ports.Endpoint) error {
	return c.notifier.Unsubscribe(ep)
}

// Send a message
func (c *ChatService) Send(m domain.Message) error {
	mm, err := c.moderator.Check(m)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to check message: %v", err))
	}
	err = c.repository.Store(*mm)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to store message: %v", err))
	}
	err = c.notifier.Broadcast(*mm)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to broadcast message: %v", err))
	}
	return nil
}

// History retrieves the chat history since a point in the past
func (c *ChatService) History(dest string, since time.Duration) ([]*domain.ModeratedMessage, error) {
	return c.repository.Fetch(dest, since)
}
