package msg

import "encoding/json"

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
	ID         string `json:"id"`
	Email      string `json:"email"`
	FullName   string `json:"fullName"`
	Company    string `json:"company"`
	Plan       string `json:"plan"`
	UserPoolID string `json:"userPoolId"`
}

// TenantSiloed is a valid MessageType.
const TenantSiloed MessageType = "TenantSiloed"

const (
	// TypeCreateTenantConfig represents a concrete value for the TenantSiloedType.
	TypeCreateTenantConfig TenantSiloedType = "TenantSiloed"
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
