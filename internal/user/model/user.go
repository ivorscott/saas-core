package model

import "github.com/go-playground/validator/v10"

var userValidator *validator.Validate

func init() {
	v := NewValidator()
	userValidator = v
}

// NewUser represents a new system user.
type NewUser struct {
	UserPoolID string `json:"userPoolId"`
	Email      string `json:"email"`
	TenantID   string `json:"tenantId"`
	FullName   string `json:"fullName"`
	Company    string `json:"company"`
}

// Validate validates the NewUser.
func (a *NewUser) Validate() error {
	return userValidator.Struct(a)
}
