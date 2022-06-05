// Package msg defines all the events/commands sent or received by services.
package msg

import "encoding/json"

// UnmarshalMsg parses the JSON-encoded data and returns Msg.
func UnmarshalMsg(data []byte) (Msg, error) {
	var m Msg
	err := json.Unmarshal(data, &m)
	return m, err
}

// Marshal JSON encodes Msg.
func (m *Msg) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalMetadata parses the JSON-encoded data and returns Metadata.
func UnmarshalMetadata(data []byte) (Metadata, error) {
	var m Metadata
	err := json.Unmarshal(data, &m)
	return m, err
}

// Marshal JSON encodes Metadata.
func (m *Metadata) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Metadata represents additional data about the request.
type Metadata struct {
	TraceID string `json:"traceId"`
	UserID  string `json:"userId"`
}

// MessageType is a type of message: a command or an event.
type MessageType string

// Msg represents a message in being sent or received.
type Msg struct {
	Data     []byte      `json:"data"`
	ID       string      `json:"id"`
	Metadata Metadata    `json:"metadata"`
	Type     MessageType `json:"type"`
}
