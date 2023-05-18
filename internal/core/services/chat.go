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
	presence   ports.Presence
}

// NewChatService creates a new instance of a chat service using
// existing repository, notifier and moderator providers
func NewChatService(repo ports.Repository, notif ports.Notifier, mod ports.Moderator, ps ports.Presence) *ChatService {
	return &ChatService{
		repository: repo,
		notifier:   notif,
		moderator:  mod,
		presence:   ps,
	}
}

// Register registers an Endpoint in the chat service
func (c *ChatService) Register(ep ports.Endpoint) error {
	err := c.notifier.Subscribe(ep)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to subscribe notifications"))
	}
	err = c.presence.Join(ep)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to join presence"))
	}
	return nil
}

// DeRegister unregisters an Endpoint in the chat service
func (c *ChatService) DeRegister(ep ports.Endpoint) error {
	err := c.presence.Leave(ep)
	if err != nil {
		return errors.Join(err, fmt.Errorf("failed to join presence"))
	}
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

	err = c.notifier.Broadcast(domain.Notication{
		Payload: mm,
		Kind:    domain.ModeratedMessageKind,
		To:      mm.To,
	})
	if err != nil {
		return errors.Join(fmt.Errorf("failed to broadcast message: %v", err))
	}
	return nil
}

// History retrieves the chat history since a point in the past
func (c *ChatService) History(dest string, since time.Duration) ([]*domain.ModeratedMessage, error) {
	return c.repository.Fetch(dest, since)
}
