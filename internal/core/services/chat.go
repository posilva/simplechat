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
	registry   ports.Registry
	repository ports.Repository
	notifier   ports.Notifier
	moderator  ports.Moderator
}

// NewChatService creates a new instance of a chat service using
// existing repository, notifier and moderator providers
func NewChatService(repo ports.Repository, notif ports.Notifier, mod ports.Moderator, reg ports.Registry) *ChatService {
	return &ChatService{
		repository: repo,
		notifier:   notif,
		moderator:  mod,
		registry:   reg,
	}
}

// Login registers an Endpoint in the chat service
func (c *ChatService) Login(ep ports.Endpoint) error {
	err := c.registry.Register(ep)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to register id: %s in room %s: %v", ep.ID(), ep.Room(), err))
	}
	return c.notifier.Subscribe(ep)
}

// Logout unregisters an Endpoint in the chat service
func (c *ChatService) Logout(ep ports.Endpoint) error {
	err := c.registry.DeRegister(ep)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to deregister id: %s: %v", ep.ID(), err))
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
