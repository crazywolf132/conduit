package conduit

import (
	"encoding/json"
	"fmt"
)

// Message represents a structured message that can be sent over the Unix socket
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// NewMessage creates a new Message with the given type and payload
func NewMessage(msgType string, payload interface{}) (*Message, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return &Message{
		Type:    msgType,
		Payload: payloadBytes,
	}, nil
}

// UnmarshalPayload unmarshals the message payload into the provided interface
func (m *Message) UnmarshalPayload(v interface{}) error {
	return json.Unmarshal(m.Payload, v)
}
