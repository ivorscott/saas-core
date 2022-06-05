package msg

import "encoding/json"

// UnmarshalCreateTenantConfigCommand parses the JSON-encoded data and returns CreateTenantConfigCommand.
func UnmarshalCreateTenantConfigCommand(data []byte) (CreateTenantConfigCommand, error) {
	var m CreateTenantConfigCommand
	err := json.Unmarshal(data, &m)
	return m, err
}

// Marshal JSON encodes CreateTenantConfigCommand.
func (m *CreateTenantConfigCommand) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

const (
	// CreateTenantConfig represents a valid MessageType.
	CreateTenantConfig MessageType = "CreateTenantConfig"
)

const (
	// TypeCreateTenantConfig represents a concrete value for the CreateTenantConfigType.
	TypeCreateTenantConfig CreateTenantConfigType = "CreateTenantConfig"
)

// CreateTenantConfigType represents a CreateTenantConfig Message.
type CreateTenantConfigType string

type CreateTenantConfigCommand struct {
	Metadata Metadata               `json:"metadata"`
	Type     CreateTenantConfigType `json:"type"`
	Data     struct {
		TenantName       string `json:"tenantName"`
		UserPoolID       string `json:"userPoolId"`
		AppClientID      string `json:"appClientId"`
		DeploymentStatus string `json:"deploymentStatus"`
	} `json:"data"`
}
