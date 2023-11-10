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
	ID              string `validate:"required"`
	Amount          int    `validate:"required"`
	Currency        string `validate:"required"`
	LastFour        string `validate:"required,len=4"`
	BankReturnCode  string
	StatusID        TransactionStatusType `validate:"required,oneof=0 1 2 3 4"`
	ExpirationMonth int                   `validate:"gte=1,lte=12"`
	ExpirationYear  int                   `validate:"min=1958"`
	SubscriptionID  string                `validate:"required"`
	PaymentIntent   string
	PaymentMethod   string `validate:"required"`
}

// Validate validates NewTransaction.
func (nt *NewTransaction) Validate() error {
	return transactionValidator.Struct(nt)
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
	UpdatedAt       time.Time             `json:"updatedAt" db:"updated_at"`
	CreatedAt       time.Time             `json:"createdAt" db:"created_at"`
}

// UpdateTransaction represents a transaction update.
type UpdateTransaction struct {
	StatusID  TransactionStatusType `json:"statusId" validate:"oneof=0 1 2 3 4"`
	UpdatedAt time.Time             `json:"updatedAt" validate:"required"`
}

// Validate validates UpdateTransaction.
func (ut *UpdateTransaction) Validate() error {
	return transactionValidator.Struct(ut)
}
