package model

import "github.com/go-playground/validator/v10"

var stripeValidator *validator.Validate

func init() {
	v := NewValidator()
	stripeValidator = v
}

// NewStripePayload represents the stripe payload for creating a new customer subscription.
type NewStripePayload struct {
	Currency        string `json:"currency" validate:"required"`
	Amount          int    `json:"amount" validate:"required,gt=0"`
	PaymentMethod   string `json:"paymentMethod" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	CardBrand       string `json:"cardBrand" validate:"required"`
	ExpirationMonth int    `json:"expirationMonth" validate:"required,gte=1,lte=12"`
	ExpirationYear  int    `json:"expirationYear" validate:"required,gte=1958"`
	ProductID       string `json:"productId" validate:"required"`
	FirstName       string `json:"firstName" validate:"required"`
	LastName        string `json:"lastName" validate:"required"`
	LastFour        string `json:"lastFour" validate:"required,len=4"`
	Plan            string `json:"plan" validate:"required"`
}

// Validate validates the NewStripePayload.
func (ns *NewStripePayload) Validate() error {
	return stripeValidator.Struct(ns)
}
