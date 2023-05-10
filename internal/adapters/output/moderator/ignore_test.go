package moderator

import (
	"testing"

	"github.com/posilva/simplechat/internal/core/domain"

	testutils "github.com/posilva/simplechat/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestNewIgnoreModerator(t *testing.T) {
	mod := NewIgnoreModerator()
	assert.NotNil(t, mod)
}

func TestIgnoreModerator_Check(t *testing.T) {
	mod := NewIgnoreModerator()
	assert.NotNil(t, mod)

	id1 := testutils.NewID()
	topic := testutils.NewUnique(testutils.Name(t))

	payload := "TestIgnoreModerator_Check Message"

	m := domain.Message{
		From:    id1,
		To:      topic,
		Payload: payload,
	}
	mm, err := mod.Check(m)
	assert.NoError(t, err)

	assert.Equal(t, mm.Message, m)

}
