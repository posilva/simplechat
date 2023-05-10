// Package moderator implements Moderator implementations
package moderator

import "github.com/posilva/simplechat/internal/core/domain"

// IgnoreModerator implements a Moderator that does not moderate
type IgnoreModerator struct {
}

// NewIgnoreModerator creates an ignore moderator
func NewIgnoreModerator() *IgnoreModerator {
	return &IgnoreModerator{}
}

// Check the messages (in this case it just ignores the moderations)
func (m *IgnoreModerator) Check(msg domain.Message) (*domain.ModeratedMessage, error) {
	return &domain.ModeratedMessage{
		Message: msg,
	}, nil
}
