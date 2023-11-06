package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var transactionValidator *validator.Validate

func init() {
	v := NewValidator()
	transactionValidator = v
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

// String returns the corresponding string value for the TransactionStatusType.
func (t TransactionStatusType) String() string {
	return [...]string{"Pending", "Cleared", "Declined", "Refunded", "Partially Refunded"}[t]
}

// NewTransaction represents a new transaction payload.
type NewTransaction struct {
	Amount               int                   `json:"amount" validate:"required,gt=0"`
	Currency             string                `json:"currency" validate:"required"`
	LastFour             string                `json:"lastFour" validate:"required,length=4"`
	BankReturnCode       string                `json:"bankReturnCode"`
	StatusID             TransactionStatusType `json:"statusId" validate:"required,oneof=0 1 2 3 4"`
	ExpirationMonth      int                   `json:"expirationMonth" validate:"required,gte=1,lte=12"`
	ExpirationYear       int                   `json:"expirationYear" validate:"required,length=4"`
	StripeSubscriptionID string                `json:"stripeSubscriptionId" validate:"required"`
	PaymentIntent        string                `json:"paymentIntent"`
	PaymentMethod        string                `json:"paymentMethod" validate:"required"`
}

// Validate validates NewTransaction.
func (nt *NewTransaction) Validate() error {
	return transactionValidator.Struct(nt)
}

// Transaction represents a stripe transaction.
type Transaction struct {
	ID                   string                `json:"id" db:"transaction_id"`
	Amount               int                   `json:"amount" db:"amount"`
	Currency             string                `json:"currency" db:"currency"`
	LastFour             string                `json:"lastFour" db:"last_four"`
	BankReturnCode       string                `json:"bankReturnCode" db:"bank_return_code"`
	StatusID             TransactionStatusType `json:"statusId" db:"transaction_status_id"`
	ExpirationMonth      int                   `json:"expirationMonth" db:"expiration_month"`
	ExpirationYear       int                   `json:"expirationYear" db:"expiration_year"`
	StripeSubscriptionID string                `json:"stripeSubscriptionID" db:"stripe_subscription_id"`
	PaymentIntent        string                `json:"paymentIntent" db:"payment_intent"`
	PaymentMethod        string                `json:"paymentMethod" db:"payment_method"`
	TenantID             string                `json:"tenantId" db:"tenant_id"`
	UpdatedAt            time.Time             `json:"updatedAt" db:"updated_at"`
	CreatedAt            time.Time             `json:"createdAt" db:"created_at"`
}

// UpdateTransaction represents a transaction update.
type UpdateTransaction struct {
	StatusID  TransactionStatusType `json:"statusId" validate:"required,oneof=0 1 2 3 4"`
	UpdatedAt time.Time             `json:"updatedAt" validate:"required"`
}

// Validate validates UpdateTransaction.
func (ut *UpdateTransaction) Validate() error {
	return transactionValidator.Struct(ut)
}
