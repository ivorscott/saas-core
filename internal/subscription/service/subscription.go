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

	"github.com/stripe/stripe-go/v72"
	"go.uber.org/zap"
)

type stripeClient interface {
	GetCustomer(customerID string) (*stripe.Customer, error)
	GetDefaultPaymentMethod(paymentMethodID string) (*stripe.PaymentMethod, error)
	CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error)
	GetPaymentMethod(method string) (*stripe.PaymentMethod, error)
	GetExistingPaymentIntent(intent string) (*stripe.PaymentIntent, error)
	SubscribeToPlan(customer *stripe.Customer, plan, last4, cardType string) (*stripe.Subscription, error)
	CreateCustomer(pm, fullName, email string) (*stripe.Customer, string, error)
	Refund(pi string, amount int) error
	CancelSubscription(subID string) error
}

type subscriptionRepository interface {
	SaveSubscription(ctx context.Context, ns model.NewSubscription, now time.Time) (model.Subscription, error)
	GetTenantSubscription(ctx context.Context, tenantID string) (model.Subscription, error)
}

// SubscriptionService is responsible for managing subscription related business logic.
type SubscriptionService struct {
	logger                *zap.Logger
	stripeClient          stripeClient
	subscriptionRepo      subscriptionRepository
	customerRepository    customerRepository
	transactionRepository transactionRepository
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
	customerRepository customerRepository,
	transactionRepository transactionRepository,
) *SubscriptionService {
	return &SubscriptionService{
		logger:                logger,
		stripeClient:          stripeClient,
		subscriptionRepo:      subscriptionRepo,
		customerRepository:    customerRepository,
		transactionRepository: transactionRepository,
	}
}

// CreatePaymentIntent creates a payment intent and returns the stripe client_secret (among other things), required
// for the frontend to collect payments.
func (ss *SubscriptionService) CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	return ss.stripeClient.CreatePaymentIntent(currency, amount)
}

// SubscribeStripeCustomer creates a new stripe customer and attaches them to a stripe subscription.
func (ss *SubscriptionService) SubscribeStripeCustomer(payload model.NewStripePayload) (string, string, string, error) {
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
		return "", "", "", ErrCreatingStripeCustomer
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
		return "", "", "", ErrSubscribingStripeCustomer
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

	return stripeSubscription.ID, stripeCustomer.ID, transactionID, nil
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

	customer, err = ss.customerRepository.GetCustomer(ctx, tenantID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			return info, web.NewRequestError(err, http.StatusNotFound)
		}
		return info, err
	}

	transactions, err = ss.transactionRepository.GetAllTransactions(ctx, tenantID)
	if err != nil {
		return info, err
	}

	subscription, err = ss.subscriptionRepo.GetTenantSubscription(ctx, tenantID)
	if err != nil {
		return info, err
	}

	paymentMethod, err = ss.stripeClient.GetDefaultPaymentMethod(customer.PaymentMethodID)
	if err != nil {
		return info, err
	}

	info.DefaultPaymentMethod = paymentMethod
	info.Transactions = transactions
	info.Subscription = subscription

	return info, nil
}

// Save persists the new subscription details.
func (ss *SubscriptionService) Save(ctx context.Context, ns model.NewSubscription, now time.Time) (model.Subscription, error) {
	var (
		s   model.Subscription
		err error
	)

	s, err = ss.subscriptionRepo.SaveSubscription(ctx, ns, now)
	if err != nil {
		switch err {
		default:
			return s, err
		}
	}
	return s, nil
}

// Cancel cancels a stripe subscription, transitioning the customer to the free tier.
func (ss *SubscriptionService) Cancel(_ context.Context) error {
	return nil
}

// Refund refunds a subscription payment.
func (ss *SubscriptionService) Refund(_ context.Context) error {
	return nil
}
