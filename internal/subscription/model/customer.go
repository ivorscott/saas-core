package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var customerValidator *validator.Validate

func init() {
	v := NewValidator()
	customerValidator = v
}

// NewCustomer represents a new customer.
type NewCustomer struct {
	ID              string `validate:"required"`
	FirstName       string `validate:"required,max=255"`
	LastName        string `validate:"required,max=255"`
	Email           string `validate:"required,email"`
	PaymentMethodID string `validate:"required"`
}

// Validate validates NewCustomer.
func (nc *NewCustomer) Validate() error {
	return customerValidator.Struct(nc)
}

// Customer represents a paying customer.
type Customer struct {
	ID              string    `json:"id" db:"customer_id"`
	TenantID        string    `json:"tenantId" db:"tenant_id"`
	FirstName       string    `json:"firstName" db:"first_name"`
	LastName        string    `json:"lastName" db:"last_name"`
	Email           string    `json:"email" db:"email"`
	PaymentMethodID string    `json:"paymentMethodId" db:"payment_method"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
}

// UpdateCustomer represents a customer update.
type UpdateCustomer struct {
	FirstName       *string   `json:"firstName" validate:"omitempty,min=1,max=255"`
	LastName        *string   `json:"lastName" validate:"omitempty,min=1,max=255"`
	Email           *string   `json:"email" validate:"omitempty,email"`
	PaymentMethodID *string   `json:"paymentMethodId" validate:"omitempty,min=1"`
	UpdatedAt       time.Time `json:"updatedAt" validate:"required"`
}

// Validate validates UpdateCustomer.
func (uc *UpdateCustomer) Validate() error {
	return customerValidator.Struct(uc)
}
