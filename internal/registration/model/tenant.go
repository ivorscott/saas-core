// Package model provides data transfer objects and validation.
package model

import "github.com/go-playground/validator/v10"

var tenantValidator *validator.Validate

func init() {
	v := NewValidator()
	tenantValidator = v
}

// NewTenant represents the new tenant.
type NewTenant struct {
	ID        string `json:"id" validate:"required"`
	Email     string `json:"email" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Company   string `json:"companyName" validate:"required"`
	Plan      string `json:"plan" validate:"required,oneof=basic premium"`
}

// Validate validates the NewTenant.
func (nt *NewTenant) Validate() error {
	return tenantValidator.Struct(nt)
}
