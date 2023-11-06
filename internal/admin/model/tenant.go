package model

import "github.com/go-playground/validator/v10"

var tenantValidator *validator.Validate

func init() {
	v := NewValidator()
	tenantValidator = v
}

// NewTenant represents a new tenant.
type NewTenant struct {
	ID          string `json:"id" validate:"required"`
	Email       string `json:"email" validate:"required"`
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	CompanyName string `json:"companyName" validate:"required"`
	Plan        string `json:"plan"`
}

// Tenant represents a tenant.
type Tenant struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	CompanyName string `json:"companyName"`
	Plan        string `json:"plan"`
	Enabled     bool   `json:"enabled"`
	Status      string `json:"status"`
	Created     string `json:"createdAt"`
}

// Validate validates the NewTenant.
func (nt *NewTenant) Validate() error {
	return tenantValidator.Struct(nt)
}
