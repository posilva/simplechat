// Package registry implements the Registry interface
package registry

import (
	"fmt"
	"strings"
	"sync"

	"github.com/posilva/simplechat/internal/core/ports"
)

type idSet map[string]struct{}

// InMemoryRegistry implements an in memory registry
type InMemoryRegistry struct {
	roomsMap map[string]idSet
	idsMap   map[string]string
	mu       sync.Mutex
}

// NewInMemoryRegistry creates a new in memory registry
func NewInMemoryRegistry() *InMemoryRegistry {
	return &InMemoryRegistry{
		idsMap:   make(map[string]string),
		roomsMap: make(map[string]idSet),
	}
}

// Register from the registry
func (r *InMemoryRegistry) Register(ep ports.Endpoint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := ep.ID()
	room := ep.Room()

	if v, ok := r.idsMap[id]; ok {
		if strings.Compare(v, room) != 0 {
			return fmt.Errorf("id: '%s' already registered to different topic: '%s'", id, room)
		}
	}
	r.idsMap[id] = room
	if v, ok := r.roomsMap[room]; ok {
		v[id] = struct{}{}
		r.roomsMap[room] = v
	} else {
		s := make(map[string]struct{})
		s[id] = struct{}{}
		r.roomsMap[room] = s
	}
	return nil
}

// DeRegister from the registry
func (r *InMemoryRegistry) DeRegister(ep ports.Endpoint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := ep.ID()
	if room, ok := r.idsMap[id]; ok {
		if v, ok := r.roomsMap[room]; ok {
			delete(v, id)
			r.roomsMap[room] = v
		}
	}
	delete(r.idsMap, id)
	return nil
}
