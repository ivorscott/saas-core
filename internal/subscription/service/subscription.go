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
	GetPlan(planID string) (*stripe.Plan, error)
	GetProduct(productID string) (*stripe.Product, error)
	GetCards(customerID string) ([]*stripe.Card, error)
	GetCustomer(customerID string) (*stripe.Customer, error)
	GetDefaultPaymentMethod(customerID string) (*stripe.PaymentMethod, error)
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
func (ss *SubscriptionService) SubscribeStripeCustomer(payload model.NewStripePayload) (string, string, error) {
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
		return "", "", ErrCreatingStripeCustomer
	}

	// subscribe to plan
	stripeSubscription, err = ss.stripeClient.SubscribeToPlan(
		stripeCustomer,
		payload.Plan,
		payload.Email,
		payload.LastFour,
		"",
	)
	if err != nil {
		ss.logger.Error(
			err.Error(),
			zap.String("email", stripeCustomer.Email),
			zap.String("plan", payload.Plan),
		)
		return "", "", ErrSubscribingStripeCustomer
	}
	ss.logger.Info(
		"successfully subscribed",
		zap.String("email", stripeCustomer.Email),
		zap.String("plan", payload.Plan),
		zap.String("subscription_id", stripeSubscription.ID),
	)

	return stripeSubscription.ID, stripeCustomer.ID, nil
}

// BillingInfo aggregates various stripe resources to build a convenient billing info summary.
func (ss *SubscriptionService) BillingInfo(ctx context.Context, tenantID string) (model.BillingInfo, error) {
	var (
		info           model.BillingInfo
		customer       model.Customer
		transactions   []model.Transaction
		plan           *stripe.Plan
		cardList       []*stripe.Card
		product        *stripe.Product
		stripeCustomer *stripe.Customer
		paymentMethod  *stripe.PaymentMethod
		err            error
	)

	customer, err = ss.customerRepository.GetCustomer(ctx, tenantID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "not found"):
			return info, web.NewRequestError(err, http.StatusNotFound)
		}
		return info, err
	}

	// get stripe customer
	stripeCustomer, err = ss.stripeClient.GetCustomer(customer.StripeCustomerID)
	if err != nil {
		return info, err
	}

	// get plan
	plan, err = ss.stripeClient.GetPlan(stripeCustomer.Subscriptions.Data[0].Plan.ID)
	if err != nil {
		return info, err
	}
	ss.logger.Info(fmt.Sprintf("=====PLAN======%+v", plan))

	product, err = ss.stripeClient.GetProduct(plan.Product.ID)
	ss.logger.Info(fmt.Sprintf("=====PRODUCT======%+v", product))
	// get card list
	cardList, err = ss.stripeClient.GetCards(stripeCustomer.ID)
	if err != nil {
		return info, err
	}
	ss.logger.Info(fmt.Sprintf("=====cardLIST======%+v", cardList))

	// get customer (default payment method)
	paymentMethod, err = ss.stripeClient.GetDefaultPaymentMethod(stripeCustomer.ID)
	if err != nil {
		return info, err
	}
	ss.logger.Info(fmt.Sprintf("=====paymentMthod======%+v", paymentMethod))

	transactions, err = ss.transactionRepository.GetAllTransactions(ctx, tenantID)
	if err != nil {
		return info, err
	}

	info.DefaultPaymentMethod = paymentMethod
	info.Cards = cardList
	info.PlanSummary = model.PlanSummary{
		Name:        product.Name,
		Description: product.Description,
		Active:      plan.Active,
		Amount:      plan.Amount,
		Currency:    string(plan.Currency),
		ProductID:   plan.Product.ID,
		Interval:    string(plan.Interval),
	}
	info.Transactions = transactions

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
