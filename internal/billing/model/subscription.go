package model

import (
	"github.com/go-playground/validator/v10"
	"time"
)

var subscriptionValidator *validator.Validate

func init() {
	v := NewValidator()
	subscriptionValidator = v
}

type SubscriptionPlan int

const (
	Basic SubscriptionPlan = iota
	Premium
)

func (s SubscriptionPlan) String() string {
	return [...]string{"Basic", "Premium"}[s]
}

type SubscriptionStatusType int

const (
	// SubscriptionStatusCleared represents a successfully cleared subscription.
	SubscriptionStatusCleared SubscriptionStatusType = iota
	// SubscriptionStatusRefunded represents a refunded subscription.
	SubscriptionStatusRefunded
	// SubscriptionStatusCancelled represents a cancelled subscription.
	SubscriptionStatusCancelled
)

func (s SubscriptionStatusType) String() string {
	return [...]string{"Cleared", "Refunded", "Cancelled"}[s]
}

type NewSubscription struct {
	ID            string                 `json:"id" validate:"required,uuid4"`
	Plan          SubscriptionPlan       `json:"plan" validate:"oneof=0 1"`
	TransactionID string                 `json:"transactionId" validate:"required,uuid4"`
	StatusID      SubscriptionStatusType `json:"statusId" validate:"required,oneof=0 1 2"`
	Amount        int                    `json:"amount" validate:"gt=0"`
	CustomerID    string                 `json:"customerId" validate:"required,uuid4"`
}

// Validate validates NewSubscription.
func (ns *NewSubscription) Validate() error {
	return subscriptionValidator.Struct(ns)
}

type Subscription struct {
	ID            string                 `json:"id" db:"subscription_id"`
	Plan          SubscriptionPlan       `json:"plan" db:"plan"`
	TransactionID string                 `json:"transactionId" db:"transaction_id"`
	StatusID      SubscriptionStatusType `json:"statusId" db:"status_id"`
	Amount        int                    `json:"amount" db:"amount"`
	TenantID      string                 `json:"tenantId" db:"tenant_id"`
	CustomerID    string                 `json:"customerId" db:"customer_id"`
	UpdatedAt     time.Time              `json:"updatedAt" db:"updated_at"`
	CreatedAt     time.Time              `json:"createdAt" db:"created_at"`
}

type UpdateSubscription struct {
	Plan          *SubscriptionPlan       `json:"plan" validate:"omitempty,oneof=0 1"`
	TransactionID *string                 `json:"transactionId" validate:"omitempty,uuid4"`
	StatusID      *SubscriptionStatusType `json:"statusId" validate:"omitempty,oneof=0 1 2"`
	Amount        *int                    `json:"amount" validate:"omitempty,gt=0"`
	UpdatedAt     time.Time               `json:"updatedAt" validate:"required"`
}

// Validate validates UpdateSubscription.
func (ns *UpdateSubscription) Validate() error {
	return subscriptionValidator.Struct(ns)
}
