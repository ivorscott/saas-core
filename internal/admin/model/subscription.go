package model

import (
	"time"

	"github.com/stripe/stripe-go/v72"
)

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
	ID            string                 `json:"id"`
	Plan          string                 `json:"plan"`
	TransactionID string                 `json:"transactionId"`
	StatusID      SubscriptionStatusType `json:"statusId"`
	Amount        int                    `json:"amount"`
	TenantID      string                 `json:"tenantId"`
	CustomerID    string                 `json:"customerId"`
	UpdatedAt     time.Time              `json:"updatedAt"`
	CreatedAt     time.Time              `json:"createdAt"`
}

// SubscriptionInfo represents a collection of subscription information.
type SubscriptionInfo struct {
	Subscription  Subscription          `json:"subscription"`
	Transactions  []Transaction         `json:"transactions"`
	PaymentMethod *stripe.PaymentMethod `json:"paymentMethod"`
}

// Transaction represents a stripe transaction.
type Transaction struct {
	ID              string                `json:"id" db:"transaction_id"`
	Amount          int                   `json:"amount" db:"amount"`
	Currency        string                `json:"currency" db:"currency"`
	LastFour        string                `json:"lastFour" db:"last_four"`
	BankReturnCode  string                `json:"bankReturnCode" db:"bank_return_code"`
	StatusID        TransactionStatusType `json:"statusId" db:"transaction_status_id"`
	ExpirationMonth int                   `json:"expirationMonth" db:"expiration_month"`
	ExpirationYear  int                   `json:"expirationYear" db:"expiration_year"`
	SubscriptionID  string                `json:"subscriptionID" db:"subscription_id"`
	PaymentIntent   string                `json:"paymentIntent" db:"payment_intent"`
	PaymentMethod   string                `json:"paymentMethod" db:"payment_method"`
	TenantID        string                `json:"tenantId" db:"tenant_id"`
	ChargeID        string                `json:"chargeId" db:"charge_id"`
	UpdatedAt       time.Time             `json:"updatedAt" db:"updated_at"`
	CreatedAt       time.Time             `json:"createdAt" db:"created_at"`
}

// TransactionStatusType describes a transaction status type.
type TransactionStatusType int

const (
	// TransactionStatusPending represents a pending transaction.
	TransactionStatusPending TransactionStatusType = iota
	// TransactionStatusCleared represents a successfully cleared transaction.
	TransactionStatusCleared
	// TransactionStatusDeclined represents a declined transaction.
	TransactionStatusDeclined
	// TransactionStatusRefunded represents a refunded transaction.
	TransactionStatusRefunded
	// TransactionStatusPartiallyRefunded represents a partially refunded transaction.
	TransactionStatusPartiallyRefunded
)
