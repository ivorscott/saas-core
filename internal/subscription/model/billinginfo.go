package model

import "github.com/stripe/stripe-go/v72"

// BillingInfo represents a collection so essential billing information.
type BillingInfo struct {
	PlanSummary          PlanSummary           `json:"plan"`
	Transactions         []Transaction         `json:"transactions"`
	Cards                []*stripe.Card        `json:"cards"`
	DefaultPaymentMethod *stripe.PaymentMethod `json:"defaultPaymentMethod"`
}

// PlanSummary represents a plan overview.
type PlanSummary struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	ProductID   string `json:"productId"`
	Interval    string `json:"interval"`
}
