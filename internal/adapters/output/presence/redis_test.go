// Package presence defines the presence component of the application
package presence

import (
	"testing"

	"github.com/posilva/simplechat/internal/adapters/output/notifier"
	"github.com/posilva/simplechat/internal/adapters/output/notifier/codecs"

	"github.com/posilva/simplechat/internal/adapters/output/registry"
	"github.com/posilva/simplechat/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func newPresence(t *testing.T) *RedisPresence {
	registry := registry.NewInMemoryRegistry()
	notifier, err := notifier.NewRabbitMQNotifierWithLocal[*codecs.JSONNotifierCodec](testutil.RabbitMQLocalURL, registry)
	assert.NoError(t, err)
	ps, err := NewRedisPresence(DefaultLocalOpts(), notifier)
	assert.NoError(t, err)
	return ps
}
func TestNewRedisPresence(t *testing.T) {
	ps := newPresence(t)
	assert.NotNil(t, ps)
}

func TestRedisPresence_Join(t *testing.T) {

	ps := newPresence(t)
	room := testutil.NewUnique(testutil.Name(t))
	ep := testutil.NewSimpleEndpoint(room)
	err := ps.Join(ep)

	assert.NoError(t, err)
}

func TestRedisPresence_Leave(t *testing.T) {
	ps := newPresence(t)
	room := testutil.NewUnique(testutil.Name(t))
	ep := testutil.NewSimpleEndpoint(room)
	err := ps.Join(ep)
	assert.NoError(t, err)

	err = ps.Leave(ep)
	assert.NoError(t, err)
}

func TestRedisPresence_Presents(t *testing.T) {
	ps := newPresence(t)
	room := testutil.NewUnique(testutil.Name(t))
	ep1 := testutil.NewSimpleEndpoint(room)
	ep2 := testutil.NewSimpleEndpoint(room)
	err := ps.Join(ep1)
	assert.NoError(t, err)
	err = ps.Join(ep2)
	assert.NoError(t, err)

	r, err := ps.Presents(ep1.Room())
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Len(t, r, 2)
}

func TestRedisPresence_IsPresent(t *testing.T) {
	ps := newPresence(t)
	room := testutil.NewUnique(testutil.Name(t))
	ep1 := testutil.NewSimpleEndpoint(room)
	err := ps.Join(ep1)
	assert.NoError(t, err)
	r, err := ps.IsPresent(ep1)
	assert.NoError(t, err)
	assert.True(t, r)
}
