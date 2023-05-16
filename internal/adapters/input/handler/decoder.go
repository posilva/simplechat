package handler

import (
	"encoding/json"

	"github.com/posilva/simplechat/internal/core/domain"
)

func decode(data []byte) (domain.Message, error) {
	var m domain.Message
	err := json.Unmarshal(data, &m)
	return m, err
}

func encode(m domain.Notication) ([]byte, error) {
	return json.Marshal(&m.Payload)
}
