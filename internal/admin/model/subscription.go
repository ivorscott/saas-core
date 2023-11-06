package model

import "time"

// SubscriptionPlanType describes the type of plan.
type SubscriptionPlanType int

const (
	// Basic is the free tier option.
	Basic SubscriptionPlanType = iota
	// Premium is the paid tier option.
	Premium
)

// String returns the corresponding string value for a subscription plan.
func (s SubscriptionPlanType) String() string {
	return [...]string{"Basic", "Premium"}[s]
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

// Subscription represents a stripe subscription for a tenant.
type Subscription struct {
	ID            string                 `json:"id" db:"subscription_id"`
	Plan          SubscriptionPlanType   `json:"plan" db:"plan"`
	TransactionID string                 `json:"transactionId" db:"transaction_id"`
	StatusID      SubscriptionStatusType `json:"statusId" db:"status_id"`
	Amount        int                    `json:"amount" db:"amount"`
	TenantID      string                 `json:"tenantId" db:"tenant_id"`
	CustomerID    string                 `json:"customerId" db:"customer_id"`
	UpdatedAt     time.Time              `json:"updatedAt" db:"updated_at"`
	CreatedAt     time.Time              `json:"createdAt" db:"created_at"`
}
