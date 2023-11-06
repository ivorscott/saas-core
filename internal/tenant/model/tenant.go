// Package model provides data transfer objects and validation.
package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var tenantValidator *validator.Validate

func init() {
	v := NewValidator()
	tenantValidator = v
}

// NewTenant represents a new Tenant.
// Tenants start on the basic plan.
type NewTenant struct {
	ID          string    `json:"id" validate:"required"`
	Email       string    `json:"email" validate:"required"`
	FirstName   string    `json:"firstName" validate:"required"`
	LastName    string    `json:"lastName" validate:"required"`
	CompanyName string    `json:"companyName" validate:"required"`
	Plan        string    `json:"plan" validate:"required,oneof=basic"`
	Status      string    `json:"status"`
	Created     time.Time `json:"createdAt"`
}

// Tenant represents a system Tenant.
type Tenant struct {
	TenantID    string `json:"id" validate:"required"`
	Email       string `json:"email" validate:"required"`
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	CompanyName string `json:"companyName" validate:"required"`
	Plan        string `json:"plan"  validate:"required,oneof=basic premium"`
	Enabled     bool   `json:"enabled"`
	Status      string `json:"status"`
	Created     string `json:"createdAt"`
}

// UpdateTenant represents a request to Tenant data.
type UpdateTenant struct {
	Email       *string `json:"email"`
	FirstName   *string `json:"firstName"`
	LastName    *string `json:"lastName"`
	CompanyName *string `json:"companyName"`
	Plan        *string `json:"plan"  validate:"oneof=basic premium"`
	Status      *string `json:"status"`
}

// Validate validates the NewTenant.
func (nt *NewTenant) Validate() error {
	return tenantValidator.Struct(nt)
}

// Validate validates the Tenant.
func (t *Tenant) Validate() error {
	return tenantValidator.Struct(t)
}
