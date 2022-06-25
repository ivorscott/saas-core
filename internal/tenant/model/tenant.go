package model

import "github.com/go-playground/validator/v10"

var tenantValidator *validator.Validate

func init() {
	v := NewValidator()
	tenantValidator = v
}

// NewTenant represents a new Tenant.
type NewTenant struct {
	ID        string `json:"id" validate:"required"`
	Email     string `json:"email" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Company   string `json:"companyName" validate:"required"`
	Plan      string `json:"plan" validate:"required,oneof=basic premium"`
}

// Tenant represents a system Tenant.
type Tenant struct {
	TenantID    string `json:"id" validate:"required"`
	Email       string `json:"email" validate:"required"`
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	CompanyName string `json:"companyName" validate:"required"`
	Plan        string `json:"plan"  validate:"required,oneof=basic premium"`
	Address     string `json:"address,omitempty"`
	City        string `json:"city,omitempty"`
	Zipcode     string `json:"zipcode,omitempty"`
	Country     string `json:"country,omitempty"`
	TaxNumber   string `json:"taxNumber,omitempty"`
}

// UpdateTenant represents a request to Tenant data.
type UpdateTenant struct {
	Email       *string `json:"email"`
	FirstName   *string `json:"firstName"`
	LastName    *string `json:"lastName"`
	CompanyName *string `json:"companyName"`
	Plan        *string `json:"plan"  validate:"oneof=basic premium"`
	Address     *string `json:"address"`
	City        *string `json:"city"`
	Zipcode     *string `json:"zipcode"`
	Country     *string `json:"country"`
	TaxNumber   *string `json:"taxNumber"`
}

// Validate validates the NewTenant.
func (nt *NewTenant) Validate() error {
	return tenantValidator.Struct(nt)
}

// Validate validates the Tenant.
func (t *Tenant) Validate() error {
	return tenantValidator.Struct(t)
}
