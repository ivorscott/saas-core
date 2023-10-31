package model

import (
	"github.com/go-playground/validator/v10"
	"time"
)

var customerValidator *validator.Validate

func init() {
	v := NewValidator()
	customerValidator = v
}

// NewCustomer represents a new customer.
type NewCustomer struct {
	ID        string `json:"id" validate:"required,uuid4"`
	FirstName string `json:"firstName" validate:"required,min=1,max=255"`
	LastName  string `json:"lastName" validate:"required,min=1,max=255"`
	Email     string `json:"email" validate:"required,email"`
}

// Validate validates NewCustomer.
func (nc *NewCustomer) Validate() error {
	return customerValidator.Struct(nc)
}

// Customer represents a paying customer.
type Customer struct {
	ID        string    `json:"id" db:"customer_id"`
	TenantID  string    `json:"tenantId" db:"tenant_id"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Email     string    `json:"email" db:"email"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// UpdateCustomer represents a customer update.
type UpdateCustomer struct {
	FirstName *string   `json:"firstName" validate:"omitempty"`
	LastName  *string   `json:"lastName" validate:"omitempty"`
	Email     *string   `json:"email" validate:"omitempty,email"`
	UpdatedAt time.Time `json:"updatedAt" validate:"required"`
}

// Validate validates UpdateCustomer.
func (uc *UpdateCustomer) Validate() error {
	return customerValidator.Struct(uc)
}
