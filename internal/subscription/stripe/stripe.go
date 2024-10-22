// Package stripe provides the stripe client and subscription logic.
package stripe

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"github.com/stripe/stripe-go/v72/refund"
	"github.com/stripe/stripe-go/v72/sub"

	"go.uber.org/zap"
)

// Client manages stripe credit cards.
type Client struct {
	logger    *zap.Logger
	key       string
	secretKey string
}

// Transaction represents a monetary transaction.
type Transaction struct {
	TransactionStatusID int
	Amount              int
	Currency            string
	// Only store the last four digits of the credit card.
	// We never store the entire credit card. That information goes to stripe.
	LastFour string
	// Code returned from Stripe.
	BankReturnCode string
}

// NewStripeClient represents a client with access to stripe methods.
func NewStripeClient(logger *zap.Logger, stripeKey string, stripeSecretKey string) *Client {
	return &Client{
		logger:    logger,
		key:       stripeKey,
		secretKey: stripeSecretKey,
	}
}

// CreatePaymentIntent creates a payment intent. PaymentIntent encapsulates details about the transaction,
// such as the supported payment methods, the amount to collect, and the desired currency.
func (c *Client) CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	var (
		pi  *stripe.PaymentIntent
		err error
	)

	stripe.Key = c.secretKey

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount)),
		Currency: stripe.String(currency),
	}
	// You can always add additional metadata to transactions.
	//params.AddMetadata("key", "value")

	pi, err = paymentintent.New(params)
	if err != nil {
		msg := ""
		if stripeErr, ok := err.(*stripe.Error); ok {
			msg = cardErrorMessage(stripeErr.Code)
		}
		return nil, msg, err
	}

	return pi, "", err
}

// GetPaymentMethod get the payment method information via payment intend id.
func (c *Client) GetPaymentMethod(paymentMethodID string) (*stripe.PaymentMethod, error) {
	stripe.Key = c.secretKey
	return paymentmethod.Get(paymentMethodID, nil)
}

// GetExistingPaymentIntent retrieves an existing payment intent.
// PaymentIntent information changes during its lifecycle.
func (c *Client) GetExistingPaymentIntent(intent string) (*stripe.PaymentIntent, error) {
	stripe.Key = c.secretKey
	return paymentintent.Get(intent, nil)
}

// SubscribeToPlan subscribes a stripe customer to a plan.
func (c *Client) SubscribeToPlan(customer *stripe.Customer, plan, last4, cardType string) (*stripe.Subscription, error) {
	var (
		subscription *stripe.Subscription
		err          error
	)

	stripe.Key = c.secretKey

	items := []*stripe.SubscriptionItemsParams{{Plan: stripe.String(plan)}}

	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customer.ID),
		Items:    items,
	}

	params.AddMetadata("last_four", last4)
	params.AddMetadata("card_type", cardType)
	params.AddExpand("latest_invoice.payment_intent")
	params.AddExpand("latest_invoice.charge.balance_transaction")

	subscription, err = sub.New(params)
	if err != nil {
		return subscription, err
	}
	return subscription, err
}

// GetCustomer retrieves the customer's profile.
func (c *Client) GetCustomer(customerID string) (*stripe.Customer, error) {
	var (
		stripeCustomer *stripe.Customer
		err            error
	)

	stripe.Key = c.secretKey

	params := &stripe.CustomerParams{}
	params.AddExpand("subscriptions")
	stripeCustomer, err = customer.Get(customerID, params)
	if err != nil {
		return stripeCustomer, err
	}

	return stripeCustomer, nil
}

// CreateCustomer creates a stripe customer.
func (c *Client) CreateCustomer(pm, fullName, email string) (*stripe.Customer, string, error) {
	var (
		newCustomer *stripe.Customer
		msg         string
		err         error
	)

	stripe.Key = c.secretKey

	params := &stripe.CustomerParams{
		Name:          stripe.String(fullName),
		PaymentMethod: stripe.String(pm),
		Email:         stripe.String(email),
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(pm),
		},
	}
	newCustomer, err = customer.New(params)
	if err != nil {
		if stripeError, ok := err.(*stripe.Error); ok {
			msg = cardErrorMessage(stripeError.Code)
		}
	}

	return newCustomer, msg, err
}

// Refund refunds an amount for a payment intent.
func (c *Client) Refund(chargeID string) error {
	stripe.Key = c.secretKey

	refundParams := &stripe.RefundParams{
		Charge: stripe.String(chargeID),
	}

	_, err := refund.New(refundParams) // Returns a refund object.
	if err != nil {
		return err
	}
	return nil
}

// CancelSubscription cancel stripe subscription.
func (c *Client) CancelSubscription(subID string) error {
	stripe.Key = c.secretKey

	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}

	_, err := sub.Update(subID, params)
	if err != nil {
		return err
	}
	return nil
}

// cardErrorMessage translates the error code into a human readable message.
func cardErrorMessage(code stripe.ErrorCode) string {
	var msg string
	switch code {
	case stripe.ErrorCodeCardDeclined:
		msg = "Your card was declined"
	case stripe.ErrorCodeExpiredCard:
		msg = "Your card is expired"
	case stripe.ErrorCodeIncorrectCVC:
		msg = "Your card verification code (CVC) is incorrect"
	case stripe.ErrorCodeIncorrectZip:
		msg = "Incorrect zip/postal code"
	case stripe.ErrorCodePostalCodeInvalid:
		msg = "Your postal code is invalid"
	case stripe.ErrorCodeAmountTooLarge:
		msg = "The amount is too large to charge to your card"
	case stripe.ErrorCodeAmountTooSmall:
		msg = "The amount is too small to charge to your card"
	case stripe.ErrorCodeBalanceInsufficient:
		msg = "Insufficient Balance"
	default:
		msg = "Your card was declined"
	}
	return msg
}
