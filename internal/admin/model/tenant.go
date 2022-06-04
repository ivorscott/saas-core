package model

import "github.com/go-playground/validator/v10"

var tenantValidator *validator.Validate

func init() {
	v := NewValidator()
	tenantValidator = v
}

// NewTenant represents the new tenant.
type NewTenant struct {
	Email    string `json:"email" validate:"required"`
	FullName string `json:"fullName" validate:"required"`
	Company  string `json:"company" validate:"required"`
	Plan     string `json:"plan" validate:"required,oneof=basic premium"`
}

// Validate validates the NewTenant.
func (nt *NewTenant) Validate() error {
	return tenantValidator.Struct(nt)
}
