package model

import "github.com/go-playground/validator/v10"

var tenantConfigValidator *validator.Validate

func init() {
	v := NewValidator()
	tenantConfigValidator = v
}

// NewTenantConfig represents a request to store configuration for a premium tenant.
type NewTenantConfig struct {
	TenantName       string `json:"tenantName"`
	UserPoolID       string `json:"userPoolID"`
	AppClientID      string `json:"appClientID"`
	DeploymentStatus string `json:"deploymentStatus"`
}

// Validate validates the  NewTenantConfig.
func (a *NewTenantConfig) Validate() error {
	return tenantConfigValidator.Struct(a)
}
