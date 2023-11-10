// Package model provides data transfer objects and validation.
package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stripe/stripe-go/v72"
)

var subscriptionValidator *validator.Validate

func init() {
	v := NewValidator()
	subscriptionValidator = v
}

// SubscriptionInfo represents a collection of subscription information.
type SubscriptionInfo struct {
	Subscription         Subscription          `json:"subscription"`
	Transactions         []Transaction         `json:"transactions"`
	DefaultPaymentMethod *stripe.PaymentMethod `json:"defaultPaymentMethod"`
}

// SubscriptionStatusType describes a subscription status type.
type SubscriptionStatusType int

const (
	// SubscriptionStatusCleared represents a successfully cleared subscription.
	SubscriptionStatusCleared SubscriptionStatusType = iota
	// SubscriptionStatusRefunded represents a refunded subscription.
	SubscriptionStatusRefunded
	// SubscriptionStatusCancelled represents a cancelled subscription.
	SubscriptionStatusCancelled
)

// String returns the corresponding string value for a subscription status.
func (s SubscriptionStatusType) String() string {
	return [...]string{"Cleared", "Refunded", "Cancelled"}[s]
}

// NewSubscription represents a new subscription payload.
type NewSubscription struct {
	ID            string                 `validate:"required"`
	Plan          string                 `validate:"required"`
	TransactionID string                 `validate:"required"`
	StatusID      SubscriptionStatusType `validate:"oneof=0 1 2"`
	Amount        int                    `validate:"gt=0"`
	CustomerID    string                 `validate:"required"`
}

// Validate validates NewSubscription.
func (ns *NewSubscription) Validate() error {
	return subscriptionValidator.Struct(ns)
}

// Subscription represents a stripe subscription for a tenant.
type Subscription struct {
	ID            string                 `json:"id" db:"subscription_id"`
	Plan          string                 `json:"plan" db:"plan"`
	TransactionID string                 `json:"transactionId" db:"transaction_id"`
	StatusID      SubscriptionStatusType `json:"statusId" db:"status_id"`
	Amount        int                    `json:"amount" db:"amount"`
	TenantID      string                 `json:"tenantId" db:"tenant_id"`
	CustomerID    string                 `json:"customerId" db:"customer_id"`
	UpdatedAt     time.Time              `json:"updatedAt" db:"updated_at"`
	CreatedAt     time.Time              `json:"createdAt" db:"created_at"`
}

// UpdateSubscription represents a subscription update.
type UpdateSubscription struct {
	Plan          string                  `json:"plan" validate:"omitempty,min=12"`
	TransactionID *string                 `json:"transactionId" validate:"omitempty,uuid4"`
	StatusID      *SubscriptionStatusType `json:"statusId" validate:"omitempty,oneof=0 1 2"`
	Amount        *int                    `json:"amount" validate:"omitempty,gt=0"`
	UpdatedAt     time.Time               `json:"updatedAt" validate:"required"`
}

// Validate validates UpdateSubscription.
func (ns *UpdateSubscription) Validate() error {
	return subscriptionValidator.Struct(ns)
}
