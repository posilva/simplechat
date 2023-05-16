// Package codecs manages the existing codecs
package codecs

import (
	"encoding/json"
	"fmt"

	"github.com/posilva/simplechat/internal/core/domain"
)

// JSONNotifierCodec manages a JSON codec
type JSONNotifierCodec struct {
}

// Encode encodes a concrete type to a byte array or error
func (c *JSONNotifierCodec) Encode(n domain.Notication) ([]byte, error) {
	return json.Marshal(n)
}

// Decode decodes a byte into a concrete type
func (c *JSONNotifierCodec) Decode(d []byte, n *domain.Notication) error {
	err := json.Unmarshal(d, n)
	if err != nil {
		return err
	}

	s, err := json.Marshal(n.Payload)
	if err != nil {
		return fmt.Errorf("failed to convert to json: %v", err)
	}

	switch n.Kind {
	case domain.MessageKind:
		fmt.Println(n.Kind)
		var m domain.Message
		err = json.Unmarshal(s, &m)
		if err != nil {
			return fmt.Errorf("failed to decode from json: %v", err)
		}
		n.Payload = m
	case domain.ChatHistoryKind:
	case domain.ModeratedMessageKind:
		var m domain.ModeratedMessage
		err = json.Unmarshal(s, &m)
		if err != nil {
			return fmt.Errorf("failed to decode from json: %v", err)
		}
		n.Payload = m
	case domain.PresenceListKind:
	case domain.PresenceUpdateKind:
	default:
		return fmt.Errorf("unknown kind notification: %v", n.Kind)
	}
	return nil
}
