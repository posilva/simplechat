// Package presence defines the presence component of the application
package presence

import (
	"context"
	"fmt"
	"time"

	"github.com/posilva/simplechat/internal/core/ports"
	"github.com/redis/rueidis"
)

// RedisPresence component
type RedisPresence struct {
	client rueidis.Client
	// TODO add a reference to a notifier to allow to notify groups of join and leave events

}

// DefaultLocalOpts returns default local options
func DefaultLocalOpts() rueidis.ClientOption {
	return rueidis.ClientOption{InitAddress: []string{"127.0.0.1:6379"}}
}

// NewRedisPresence creates a new presence component using redis
func NewRedisPresence(opts rueidis.ClientOption) (*RedisPresence, error) {
	client, err := rueidis.NewClient(opts)
	if err != nil {
		return nil, err
	}
	return &RedisPresence{
		client: client,
	}, nil
}

// Join the Presence component
func (p *RedisPresence) Join(ep ports.Endpoint) error {
	err := p.doJoin(ep)
	if err != nil {
		return err
	}

	return nil
}

// Leave the Presence component
func (p *RedisPresence) Leave(ep ports.Endpoint) error {
	err := p.doLeave(ep)
	if err != nil {
		return err
	}
	return nil
}

func (p *RedisPresence) doJoin(ep ports.Endpoint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	key := fmt.Sprintf("room::%s", ep.Room())
	f := fmt.Sprintf("id::" + ep.ID())
	v := fmt.Sprintf("%d", time.Now().UTC().Unix())

	cmd := p.client.B().Hsetnx().Key(key).Field(f).Value(v).Build()
	err := p.client.Do(ctx, cmd).Error()
	if err != nil {
		return fmt.Errorf("failed to set field value in redis presence: %v", err)
	}
	return nil
}

func (p *RedisPresence) doLeave(ep ports.Endpoint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	key := fmt.Sprintf("room::%s", ep.Room())
	f := fmt.Sprintf("id::" + ep.ID())

	cmd := p.client.B().Hdel().Key(key).Field(f).Build()
	err := p.client.Do(ctx, cmd).Error()
	if err != nil {
		return fmt.Errorf("failed to set field value in redis presence: %v", err)
	}
	return nil
}
