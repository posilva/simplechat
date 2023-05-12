package registry

import (
	"testing"

	"github.com/posilva/simplechat/internal/core/domain"
	"github.com/posilva/simplechat/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryRegistry_New(t *testing.T) {

	r := NewInMemoryRegistry()
	exp := &InMemoryRegistry{}
	assert.NotNil(t, r)
	assert.IsType(t, exp, r)
}

func TestInMemoryRegistry_Register(t *testing.T) {
	r := NewInMemoryRegistry()
	assert.NotNil(t, r)

	id := testutil.NewID()
	id2 := testutil.NewID()
	room := testutil.NewUnique(testutil.Name(t))
	rc := testutil.NewTestReceiver()

	ep := testutil.NewTestEndpoint(id, room, rc)
	err := r.Register(ep)
	assert.NoError(t, err)

	ep2 := testutil.NewTestEndpoint(id2, room, rc)
	err = r.Register(ep2)
	assert.NoError(t, err)
}
func TestInMemoryRegistry_DeRegister(t *testing.T) {
	r := NewInMemoryRegistry()
	assert.NotNil(t, r)

	id := testutil.NewID()
	room := testutil.NewUnique(testutil.Name(t))
	rc := testutil.NewTestReceiver()

	ep := testutil.NewTestEndpoint(id, room, rc)
	err := r.Register(ep)
	assert.NoError(t, err)
	err = r.DeRegister(ep)
	assert.NoError(t, err)
}
func TestInMemoryRegistry_Notify(t *testing.T) {
	r := NewInMemoryRegistry()
	assert.NotNil(t, r)

	id := testutil.NewID()
	room := testutil.NewUnique(testutil.Name(t))
	rc := testutil.NewTestReceiver()

	ep := testutil.NewTestEndpoint(id, room, rc)
	err := r.Register(ep)
	assert.NoError(t, err)

	m := domain.ModeratedMessage{
		Message: domain.Message{
			Payload: testutil.Name(t),
			From:    id,
			To:      room,
		},
		Level:           0,
		FilteredPayload: testutil.Name(t),
	}
	r.Notify(m)

	mm := <-rc.Channel()
	assert.Equal(t, mm, m)
}
