package model

import "github.com/go-playground/validator/v10"

var tenantValidator *validator.Validate

func init() {
	v := NewValidator()
	tenantValidator = v
}

// NewTenant represents a new Tenant.
type NewTenant struct {
	ID       string `json:"id" validate:"required"`
	Email    string `json:"email" validate:"required"`
	FullName string `json:"fullName" validate:"required"`
	Company  string `json:"companyName" validate:"required"`
	Plan     string `json:"plan" validate:"required,oneof=basic premium"`
}

// Tenant represents a system Tenant.
type Tenant struct {
	ID        string `json:"id" validate:"required"`
	Email     string `json:"email" validate:"required"`
	FullName  string `json:"fullName" validate:"required"`
	Company   string `json:"companyName" validate:"required"`
	Plan      string `json:"plan"  validate:"required,oneof=basic premium"`
	Address   string `json:"address,omitempty"`
	City      string `json:"city,omitempty"`
	Zipcode   string `json:"zipcode,omitempty"`
	Country   string `json:"country,omitempty"`
	TaxNumber string `json:"taxNumber,omitempty"`
}

// UpdateTenant represents a request to Tenant data.
type UpdateTenant struct {
	Email     string  `json:"email" validate:"required"`
	FullName  string  `json:"fullName" validate:"required"`
	Company   string  `json:"companyName" validate:"required"`
	Plan      string  `json:"plan"  validate:"required,oneof=basic premium"`
	Address   *string `json:"address,omitempty"`
	City      *string `json:"city,omitempty"`
	Zipcode   *string `json:"zipcode,omitempty"`
	Country   *string `json:"country,omitempty"`
	TaxNumber *string `json:"taxNumber,omitempty"`
}

type NewSiloConfig struct {
	TenantName       string
	UserPoolID       string
	AppClientID      string
	DeploymentStatus string
}

// Validate validates the NewTenant.
func (nt *NewTenant) Validate() error {
	return tenantValidator.Struct(nt)
}

// Validate validates the Tenant.
func (t *Tenant) Validate() error {
	return tenantValidator.Struct(t)
}
