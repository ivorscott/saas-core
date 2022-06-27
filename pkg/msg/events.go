package msg

import (
	"encoding/json"
)

// UnmarshalTenantIdentityCreatedEvent parses the JSON-encoded data and returns TenantIdentityCreatedEvent.
func UnmarshalTenantIdentityCreatedEvent(data []byte) (TenantIdentityCreatedEvent, error) {
	var m TenantIdentityCreatedEvent
	err := json.Unmarshal(data, &m)
	return m, err
}

// Marshal JSON encodes TenantIdentityCreatedEvent.
func (m *TenantIdentityCreatedEvent) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalTenantSiloedEvent parses the JSON-encoded data and returns TenantSiloedEvent.
func UnmarshalTenantSiloedEvent(data []byte) (TenantSiloedEvent, error) {
	var m TenantSiloedEvent
	err := json.Unmarshal(data, &m)
	return m, err
}

// Marshal JSON encodes TenantSiloedEvent.
func (m *TenantSiloedEvent) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalTenantRegisteredEvent parses the JSON-encoded data and returns TenantRegisteredEvent.
func UnmarshalTenantRegisteredEvent(data []byte) (TenantRegisteredEvent, error) {
	var m TenantRegisteredEvent
	err := json.Unmarshal(data, &m)
	return m, err
}

// Marshal JSON encodes CreateTenantConfigCommand.
func (m *TenantRegisteredEvent) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// TenantRegistered is a valid MessageType.
const TenantRegistered MessageType = "TenantRegistered"

// TenantRegisteredType represents a TenantRegistered event type.
type TenantRegisteredType string

const (
	// TypeTenantRegistered represents a concrete value for the TenantRegisteredType.
	TypeTenantRegistered TenantRegisteredType = "TenantRegistered"
)

// TenantRegisteredEvent represents a TenantRegistered Message.
type TenantRegisteredEvent struct {
	Metadata Metadata                  `json:"metadata"`
	Type     TenantRegisteredType      `json:"type"`
	Data     TenantRegisteredEventData `json:"data"`
}

type TenantRegisteredEventData struct {
	TenantID   string `json:"tenantId"`
	Email      string `json:"email"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Company    string `json:"company"`
	Plan       string `json:"plan"`
	UserPoolID string `json:"userPoolId"`
}

// TenantSiloed is a valid MessageType.
const TenantSiloed MessageType = "TenantSiloed"

const (
	// TypeTenantSiloed represents a concrete value for the TenantSiloedType.
	TypeTenantSiloed TenantSiloedType = "TenantSiloed"
)

// TenantSiloedType represents a TenantSiloed Message.
type TenantSiloedType string

type TenantSiloedEvent struct {
	Metadata Metadata              `json:"metadata"`
	Type     TenantSiloedType      `json:"type"`
	Data     TenantSiloedEventData `json:"data"`
}

type TenantSiloedEventData struct {
	TenantName       string `json:"tenantName"`
	UserPoolID       string `json:"userPoolId"`
	AppClientID      string `json:"appClientId"`
	DeploymentStatus string `json:"deploymentStatus"`
}

// TenantIdentityCreated is a valid MessageType.
const TenantIdentityCreated MessageType = "TenantIdentityCreated"

const (
	// TypeTenantIdentityCreated represents a concrete value for the TenantIdentityCreatedType.
	TypeTenantIdentityCreated TenantIdentityCreatedType = "TenantIdentityCreated"
)

// TenantIdentityCreatedType represents a TenantIdentityCreated Message.
type TenantIdentityCreatedType string

type TenantIdentityCreatedEvent struct {
	Metadata Metadata                       `json:"metadata"`
	Type     TenantIdentityCreatedType      `json:"type"`
	Data     TenantIdentityCreatedEventData `json:"data"`
}

type TenantIdentityCreatedEventData struct {
	TenantID  string `json:"tenantId"`
	UserID    string `json:"userId"`
	Company   string `json:"company"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Plan      string `json:"plan"`
	CreatedAt string `json:"createdAt"`
}
