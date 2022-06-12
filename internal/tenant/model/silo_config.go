package model

import "github.com/go-playground/validator/v10"

var siloConfigValidator *validator.Validate

func init() {
	v := NewValidator()
	siloConfigValidator = v
}

// NewSiloConfig represents a request to store configuration for a premium tenant.
type NewSiloConfig struct {
	TenantName       string `json:"tenantName"`
	UserPoolID       string `json:"userPoolId"`
	AppClientID      string `json:"appClientId"`
	DeploymentStatus string `json:"deploymentStatus"`
}

// Validate validates the  NewSiloConfig.
func (a *NewSiloConfig) Validate() error {
	return siloConfigValidator.Struct(a)
}
