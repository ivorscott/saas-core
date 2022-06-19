package msg

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	StreamTenants     = "TENANTS"
	SubjectRegistered = "TENANTS.registered"
	SubjectSiloed     = "TENANTS.siloed"

	StreamMemberships        = "MEMBERSHIPS"
	SubjectMembershipCreated = "MEMBERSHIPS.created"
	SubjectMembershipUpdated = "MEMBERSHIPS.updated"
	SubjectMembershipDeleted = "MEMBERSHIPS.deleted"
	SubjectMembershipInvited = "MEMBERSHIPS.invited"

	StreamProjects             = "PROJECTS"
	SubjectProjectCreated      = "PROJECTS.created"
	SubjectProjectUpdated      = "PROJECTS.updated"
	SubjectProjectDeleted      = "PROJECTS.deleted"
	SubjectProjectTeamAssigned = "PROJECTS.assigned"
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

// Bytes returns message bytes.
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
	TraceID  string `json:"traceId"`
	UserID   string `json:"userId"`
	TenantID string `json:"tenantID"`
}

// MessageType is a type of message.
type MessageType string

// Msg represents a message in being sent or received.
type Msg struct {
	Data     interface{} `json:"data"`
	Metadata Metadata    `json:"metadata"`
	Type     MessageType `json:"type"`
}

// ParseTime converts a time string to time.Time.
func ParseTime(ts string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", ts)
	if err != nil {
		panic("failed to parse time")
	}
	return t
}
