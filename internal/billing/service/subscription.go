// Package service manages the application layer for handling business logic.
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/devpies/saas-core/internal/billing/model"

	"github.com/stripe/stripe-go/v72"
	"go.uber.org/zap"
)

type stripeClient interface {
	CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error)
	GetPaymentMethod(method string) (*stripe.PaymentMethod, error)
	GetExistingPaymentIntent(intent string) (*stripe.PaymentIntent, error)
	SubscribeToPlan(cust *stripe.Customer, plan, email, last4, cardType string) (*stripe.Subscription, error)
	CreateCustomer(pm, fullName, email string) (*stripe.Customer, string, error)
	Refund(pi string, amount int) error
	CancelSubscription(subID string) error
}

type subscriptionRepository interface {
	SaveSubscription(ctx context.Context, ns model.NewSubscription, now time.Time) (model.Subscription, error)
	GetTenantSubscription(ctx context.Context, id string) (model.Subscription, error)
}

// SubscriptionService is responsible for managing subscription related business logic.
type SubscriptionService struct {
	logger           *zap.Logger
	stripeClient     stripeClient
	subscriptionRepo subscriptionRepository
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
) *SubscriptionService {
	return &SubscriptionService{
		logger:           logger,
		stripeClient:     stripeClient,
		subscriptionRepo: subscriptionRepo,
	}
}

// CreatePaymentIntent creates a payment intent and returns the stripe client_secret (among other things), required
// for the frontend to collect payments.
func (ss *SubscriptionService) CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	return ss.stripeClient.CreatePaymentIntent(currency, amount)
}

// SubscribeStripeCustomer creates a new stripe customer and attaches them to a stripe subscription.
func (ss *SubscriptionService) SubscribeStripeCustomer(payload model.NewStripePayload) (string, error) {
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
		return "", ErrCreatingStripeCustomer
	}

	// subscribe to plan
	stripeSubscription, err = ss.stripeClient.SubscribeToPlan(
		stripeCustomer,
		payload.Plan.String(),
		payload.Email,
		payload.LastFour,
		"",
	)
	if err != nil {
		ss.logger.Error(
			err.Error(),
			zap.String("email", stripeCustomer.Email),
			zap.String("plan", payload.Plan.String()),
		)
		return "", ErrSubscribingStripeCustomer
	}
	ss.logger.Info(
		"successfully subscribed",
		zap.String("email", stripeCustomer.Email),
		zap.String("plan", payload.Plan.String()),
		zap.String("subscription_id", stripeSubscription.ID),
	)

	return stripeSubscription.ID, nil
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

// GetOne returns the subscription for the provided tenant.
func (ss *SubscriptionService) GetOne(ctx context.Context, tenantID string) (model.Subscription, error) {
	var (
		s   model.Subscription
		err error
	)

	s, err = ss.subscriptionRepo.GetTenantSubscription(ctx, tenantID)
	if err != nil {
		switch err {
		default:
			return s, err
		}
	}

	return s, nil
}

// Cancel cancels a stripe subscription, transitioning the customer to the free tier.
func (ss *SubscriptionService) Cancel(ctx context.Context) error {
	return nil
}

// Refund refunds a subscription payment.
func (ss *SubscriptionService) Refund(ctx context.Context) error {
	return nil
}
