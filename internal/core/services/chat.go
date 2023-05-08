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
func NewChatService(r ports.Repository, n ports.Notifier, m ports.Moderator) *ChatService {
	return &ChatService{
		repository: r,
		notifier:   n,
		moderator:  m,
	}
}

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

func (c *ChatService) History(dest string, since time.Duration) ([]*domain.ModeratedMessage, error) {
	return c.repository.Fetch(dest, since)
}
