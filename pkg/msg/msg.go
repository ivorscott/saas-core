package msg

import (
	"encoding/json"
	"fmt"
)

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

func Bytes(message interface{}) ([]byte, error) {
	v, ok := message.(*Msg)
	if !ok {
		return []byte{}, fmt.Errorf("not a message")
	}
	data, err := v.Marshal()
	if err != nil {
		return []byte{}, fmt.Errorf("mashalling failed")
	}
	return data, nil
}

// Metadata represents additional data about the request.
type Metadata struct {
	TraceID string `json:"traceId"`
	UserID  string `json:"userId"`
}

// Msg represents a message in being sent or received.
type Msg struct {
	Data     interface{} `json:"data"`
	Metadata Metadata    `json:"metadata"`
	Type     string      `json:"type"`
}
