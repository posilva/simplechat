// Package presence defines the presence component of the application
package presence

import (
	"context"
	"fmt"
	"time"

	"github.com/posilva/simplechat/internal/core/domain"
	"github.com/posilva/simplechat/internal/core/ports"
	"github.com/redis/rueidis"
)

// RedisPresence component
type RedisPresence struct {
	client   rueidis.Client
	notifier ports.Notifier
	log      ports.Logger
}

// DefaultLocalOpts returns default local options
func DefaultLocalOpts() rueidis.ClientOption {
	return rueidis.ClientOption{InitAddress: []string{"127.0.0.1:6379"}}
}

// NewRedisPresence creates a new presence component using redis
func NewRedisPresence(opts rueidis.ClientOption, n ports.Notifier, log ports.Logger) (*RedisPresence, error) {
	client, err := rueidis.NewClient(opts)
	if err != nil {
		return nil, err
	}
	return &RedisPresence{
		client:   client,
		notifier: n,
		log:      log,
	}, nil
}

// Join the Presence component
func (p *RedisPresence) Join(ep ports.Endpoint) error {
	err := p.doJoin(ep)
	if err != nil {
		return fmt.Errorf("failed to join: %v", err)
	}
	err = p.notifier.Broadcast(domain.Notication{
		Kind: domain.PresenceJoinKind,
		To:   ep.Room(),
		Payload: domain.PresenceUpdate{
			ID:        ep.ID(),
			Action:    domain.PresenceUpdateJoinAction,
			Timestamp: uint64(time.Now().Unix()),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to broadcast join update: %v", err)
	}
	return nil
}

// Leave the Presence component
func (p *RedisPresence) Leave(ep ports.Endpoint) error {
	err := p.doLeave(ep)
	if err != nil {
		return fmt.Errorf("failed to leave: %v", err)
	}
	err = p.notifier.Broadcast(domain.Notication{
		UUID: fmt.Sprintf("ep:%v:%v", ep.ID(), ep.Room()),
		Kind: domain.PresenceLeaveKind,
		Payload: domain.PresenceUpdate{
			ID:        ep.ID(),
			Action:    domain.PresenceUpdateLeaveAction,
			Timestamp: uint64(time.Now().Unix()),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to broadcast join update: %v", err)
	}
	return nil
}

func (p *RedisPresence) doJoin(ep ports.Endpoint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	key := key(ep.Room())
	f := idField(ep.ID())
	v := fmt.Sprintf("%d", time.Now().UTC().Unix())

	cmd := p.client.B().Hsetnx().Key(key).Field(f).Value(v).Build()
	err := p.client.Do(ctx, cmd).Error()
	if err != nil {
		return fmt.Errorf("failed to set field value in redis presence: %v", err)
	}
	return nil
}

// Presents returns the participants of a room
func (p *RedisPresence) Presents(room string) (v map[string]string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	key := key(room)

	cmd := p.client.B().Hgetall().Key(key).Build()
	res, err := p.client.Do(ctx, cmd).AsStrMap()

	if err != nil {
		return nil, fmt.Errorf("failed to set field value in redis presence: %v", err)
	}

	return res, nil
}

// IsPresent returns true if a id is present in the room
func (p *RedisPresence) IsPresent(ep ports.Endpoint) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	key := key(ep.Room())
	f := idField(ep.ID())

	cmd := p.client.B().Hexists().Key(key).Field(f).Build()
	b, err := p.client.Do(ctx, cmd).AsBool()
	if err != nil {
		return false, fmt.Errorf("failed check if endpoint is present: %v", err)
	}
	return b, err
}
func (p *RedisPresence) doLeave(ep ports.Endpoint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	key := key(ep.Room())
	f := idField(ep.ID())

	cmd := p.client.B().Hdel().Key(key).Field(f).Build()
	err := p.client.Do(ctx, cmd).Error()
	if err != nil {
		return fmt.Errorf("failed to set field value in redis presence: %v", err)
	}
	return nil
}

func key(room string) string {
	return fmt.Sprintf("room::%s", room)
}

func idField(id string) string {
	return fmt.Sprintf("id::" + id)
}
