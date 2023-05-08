package moderator

import "github.com/posilva/simplechat/internal/core/domain"

type IgnoreModerator struct {
}

func NewIgnoreModerator() *IgnoreModerator {
	return &IgnoreModerator{}
}

func (m *IgnoreModerator) Check(msg domain.Message) (*domain.ModeratedMessage, error) {
	return &domain.ModeratedMessage{
		Message: msg,
	}, nil
}
