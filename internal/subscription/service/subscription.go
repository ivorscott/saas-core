// Package service manages the application layer for handling business logic.
package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/devpies/saas-core/internal/subscription/model"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/jmoiron/sqlx"
	"github.com/stripe/stripe-go/v72"
	"go.uber.org/zap"
)

type stripeClient interface {
	GetCustomer(customerID string) (*stripe.Customer, error)
	CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error)
	GetPaymentMethod(paymentMethodID string) (*stripe.PaymentMethod, error)
	GetExistingPaymentIntent(intent string) (*stripe.PaymentIntent, error)
	SubscribeToPlan(customer *stripe.Customer, plan, last4, cardType string) (*stripe.Subscription, error)
	CreateCustomer(pm, fullName, email string) (*stripe.Customer, string, error)
	Refund(pi string, amount int) error
	CancelSubscription(subID string) error
}

type subscriptionRepository interface {
	SaveSubscriptionTx(ctx context.Context, tx *sqlx.Tx, ns model.NewSubscription, now time.Time) (model.Subscription, error)
	GetTenantSubscription(ctx context.Context, tenantID string) (model.Subscription, error)
}

// SubscriptionService is responsible for managing subscription related business logic.
type SubscriptionService struct {
	logger           *zap.Logger
	stripeClient     stripeClient
	subscriptionRepo subscriptionRepository
	customerRepo     customerRepository
	transactionRepo  transactionRepository
}

var (
	// ErrCreatingStripeCustomer represents an error creating a stripe customer.
	ErrCreatingStripeCustomer = errors.New("error creating stripe customer")
	// ErrSubscribingStripeCustomer represents an error creating a subscription for the stripe customer.
	ErrSubscribingStripeCustomer = errors.New("error subscribing stripe customer")
)

// NewSubscriptionService returns a new SubscriptionService.
func NewSubscriptionService(
	logger *zap.Logger,
	stripeClient stripeClient,
	subscriptionRepo subscriptionRepository,
	customerRepo customerRepository,
	transactionRepo transactionRepository,
) *SubscriptionService {
	return &SubscriptionService{
		logger:           logger,
		stripeClient:     stripeClient,
		subscriptionRepo: subscriptionRepo,
		customerRepo:     customerRepo,
		transactionRepo:  transactionRepo,
	}
}

// CreatePaymentIntent creates a payment intent and returns the stripe client_secret (among other things), required
// for the frontend to collect payments.
func (ss *SubscriptionService) CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	return ss.stripeClient.CreatePaymentIntent(currency, amount)
}

// SubscribeStripeCustomer creates a new stripe customer and attaches them to a stripe subscription.
func (ss *SubscriptionService) SubscribeStripeCustomer(ctx context.Context, payload model.NewStripePayload) error {
	var (
		stripeSubscription *stripe.Subscription
		stripeCustomer     *stripe.Customer
		customerMsg        string
		fullName           = fmt.Sprintf("%s %s", payload.FirstName, payload.LastName)
		err                error
	)

	// create customer
	stripeCustomer, customerMsg, err = ss.stripeClient.CreateCustomer(payload.PaymentMethod, fullName, payload.Email)
	if err != nil {
		ss.logger.Error(
			fmt.Sprintf("%s: %s", customerMsg, err.Error()),
			zap.String("email", stripeCustomer.Email),
		)
		return ErrCreatingStripeCustomer
	}

	// subscribe to plan
	stripeSubscription, err = ss.stripeClient.SubscribeToPlan(
		stripeCustomer,
		payload.Plan,
		payload.LastFour,
		"",
	)
	if err != nil {
		ss.logger.Error(
			err.Error(),
			zap.String("email", stripeCustomer.Email),
			zap.String("plan", payload.Plan),
		)
		return ErrSubscribingStripeCustomer
	}
	ss.logger.Info(
		"successfully subscribed",
		zap.String("email", stripeCustomer.Email),
		zap.String("plan", payload.Plan),
		zap.String("subscription_id", stripeSubscription.ID),
	)
	var transactionID string

	if stripeSubscription.LatestInvoice != nil {
		if stripeSubscription.LatestInvoice.Charge != nil {
			if stripeSubscription.LatestInvoice.Charge.BalanceTransaction != nil {
				transactionID = stripeSubscription.LatestInvoice.Charge.BalanceTransaction.ID
			}
		}
	}

	return ss.saveStripeResources(ctx, payload, stripeSubscription.ID, stripeCustomer.ID, transactionID)
}

func (ss *SubscriptionService) saveStripeResources(
	ctx context.Context,
	payload model.NewStripePayload,
	subscriptionID,
	customerID,
	transactionID string,
) error {
	var err error

	c, err := newCustomer(customerID, payload)
	if err != nil {
		return err
	}
	s, err := newSubscription(customerID, transactionID, subscriptionID, payload)
	if err != nil {
		return err
	}
	t, err := newTransaction(transactionID, subscriptionID, payload)
	if err != nil {
		return err
	}

	err = ss.customerRepo.RunTx(ctx, func(tx *sqlx.Tx) error {
		_, err = ss.customerRepo.SaveCustomerTx(ctx, tx, c, time.Now())
		if err != nil {
			return err
		}
		_, err = ss.subscriptionRepo.SaveSubscriptionTx(ctx, tx, s, time.Now())
		if err != nil {
			return err
		}
		_, err = ss.transactionRepo.SaveTransactionTx(ctx, tx, t, time.Now())
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}

func newCustomer(customerID string, payload model.NewStripePayload) (model.NewCustomer, error) {
	c := model.NewCustomer{
		ID:              customerID,
		FirstName:       payload.FirstName,
		LastName:        payload.LastName,
		Email:           payload.Email,
		PaymentMethodID: payload.PaymentMethod,
	}
	if err := c.Validate(); err != nil {
		return model.NewCustomer{}, web.NewRequestError(err, http.StatusBadRequest)
	}
	return c, nil
}

func newTransaction(transactionID, subscriptionID string, payload model.NewStripePayload) (model.NewTransaction, error) {
	t := model.NewTransaction{
		ID:              transactionID,
		Amount:          payload.Amount,
		Currency:        payload.Currency,
		LastFour:        payload.LastFour,
		StatusID:        model.TransactionStatusCleared,
		ExpirationMonth: payload.ExpirationMonth,
		ExpirationYear:  payload.ExpirationYear,
		PaymentMethod:   payload.PaymentMethod,
		SubscriptionID:  subscriptionID,
	}
	if err := t.Validate(); err != nil {
		return model.NewTransaction{}, web.NewRequestError(err, http.StatusBadRequest)
	}
	return t, nil
}

func newSubscription(customerID, transactionID, subscriptionID string, payload model.NewStripePayload) (model.NewSubscription, error) {
	s := model.NewSubscription{
		ID:            subscriptionID,
		Plan:          payload.Plan,
		TransactionID: transactionID,
		StatusID:      model.SubscriptionStatusCleared,
		Amount:        payload.Amount,
		CustomerID:    customerID,
	}
	if err := s.Validate(); err != nil {
		return model.NewSubscription{}, web.NewRequestError(err, http.StatusBadRequest)
	}
	return s, nil
}

// SubscriptionInfo aggregates various stripe resources to show convenient subscription information.
func (ss *SubscriptionService) SubscriptionInfo(ctx context.Context, tenantID string) (model.SubscriptionInfo, error) {
	var (
		info          model.SubscriptionInfo
		customer      model.Customer
		transactions  []model.Transaction
		subscription  model.Subscription
		paymentMethod *stripe.PaymentMethod
		err           error
	)

	customer, err = ss.customerRepo.GetCustomer(ctx, tenantID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return info, web.NewRequestError(err, http.StatusNotFound)
		}
		return info, err
	}

	transactions, err = ss.transactionRepo.GetAllTransactions(ctx, tenantID)
	if err != nil {
		return info, err
	}

	subscription, err = ss.subscriptionRepo.GetTenantSubscription(ctx, tenantID)
	if err != nil {
		return info, err
	}

	paymentMethod, err = ss.stripeClient.GetPaymentMethod(customer.PaymentMethodID)
	if err != nil {
		return info, err
	}

	info.PaymentMethod = paymentMethod
	info.Transactions = transactions
	info.Subscription = subscription

	return info, nil
}

// Cancel cancels a stripe subscription, transitioning the customer to the free tier.
func (ss *SubscriptionService) Cancel(subID string) error {
	return ss.stripeClient.CancelSubscription(subID)
}

// Refund refunds a subscription payment.
func (ss *SubscriptionService) Refund() error {
	var (
		premiumPlanAmount = 1000
		pi                *stripe.PaymentIntent
		msg               string
		err               error
	)
	pi, msg, err = ss.stripeClient.CreatePaymentIntent("eur", premiumPlanAmount)
	if err != nil {
		ss.logger.Error("failed creating payment intent", zap.String("message", msg))
		return err
	}
	err = ss.stripeClient.Refund(pi.ID, premiumPlanAmount)
	if err != nil {
		return err
	}
	return nil
}
