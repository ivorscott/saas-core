// Package handler manages the presentation layer for handling incoming requests.
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/devpies/saas-core/internal/subscription/model"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/go-chi/chi/v5"
	"github.com/stripe/stripe-go/v72"
	"go.uber.org/zap"
)

type subscriptionService interface {
	Refund(ctx context.Context) error
	Cancel(ctx context.Context) error
	Save(ctx context.Context, ns model.NewSubscription, now time.Time) (model.Subscription, error)
	SubscriptionInfo(ctx context.Context, tenantID string) (model.SubscriptionInfo, error)
	CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error)
	SubscribeStripeCustomer(nsp model.NewStripePayload) (string, string, error)
}

type transactionService interface {
	Save(ctx context.Context, nt model.NewTransaction, now time.Time) (model.Transaction, error)
}

type customerService interface {
	Save(ctx context.Context, nc model.NewCustomer, now time.Time) (model.Customer, error)
}

// SubscriptionHandler handles subscription related requests.
type SubscriptionHandler struct {
	logger              *zap.Logger
	subscriptionService subscriptionService
	transactionService  transactionService
	customerService     customerService
}

// NewSubscriptionHandler returns a SubscriptionHandler.
func NewSubscriptionHandler(
	logger *zap.Logger,
	subscriptionService subscriptionService,
	transactionService transactionService,
	customerService customerService,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		logger:              logger,
		subscriptionService: subscriptionService,
		transactionService:  transactionService,
		customerService:     customerService,
	}
}

func newCustomer(customerID string, payload model.NewStripePayload) (model.NewCustomer, error) {
	c := model.NewCustomer{
		ID:        customerID,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
	}
	if err := c.Validate(); err != nil {
		return model.NewCustomer{}, web.NewRequestError(err, http.StatusBadRequest)
	}
	return c, nil
}

func newTransaction(subscriptionID string, payload model.NewStripePayload) (model.NewTransaction, error) {
	t := model.NewTransaction{
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

// Create sets up a new subscription for the customer.
func (sh *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var (
		payload model.NewStripePayload
		err     error
	)

	if err = web.Decode(r, &payload); err != nil {
		return err
	}

	// TODO: opt for running in postgres transaction
	subscriptionID, customerID, err := sh.subscriptionService.SubscribeStripeCustomer(payload)
	if err != nil {
		return err
	}

	c, err := newCustomer(customerID, payload)
	if err != nil {
		return err
	}

	_, err = sh.customerService.Save(r.Context(), c, time.Now())
	if err != nil {
		return err
	}

	t, err := newTransaction(subscriptionID, payload)
	if err != nil {
		return err
	}

	transaction, err := sh.transactionService.Save(r.Context(), t, time.Now())
	if err != nil {
		return err
	}

	s, err := newSubscription(customerID, transaction.ID, subscriptionID, payload)
	if err != nil {
		return err
	}

	_, err = sh.subscriptionService.Save(r.Context(), s, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

// SubscriptionInfo aggregates various stripe resource objects into a billing summary.
func (sh *SubscriptionHandler) SubscriptionInfo(w http.ResponseWriter, r *http.Request) error {
	var (
		tenantID = chi.URLParam(r, "tenantID")
		info     model.SubscriptionInfo
		err      error
	)

	info, err = sh.subscriptionService.SubscriptionInfo(r.Context(), tenantID)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, info, http.StatusOK)
}

// GetPaymentIntent retrieves the paymentIntent from stripe.
func (sh *SubscriptionHandler) GetPaymentIntent(w http.ResponseWriter, r *http.Request) error {
	var (
		payload struct {
			Currency string
			Amount   int
		}
		err error
	)

	if err = web.Decode(r, &payload); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	pi, msg, err := sh.subscriptionService.CreatePaymentIntent(payload.Currency, payload.Amount)
	if err != nil {
		sh.logger.Error("creating payment intent failed", zap.String("message", msg))
		return err
	}

	return web.Respond(r.Context(), w, pi, http.StatusOK)
}

// Cancel cancels a paid subscription.
func (sh *SubscriptionHandler) Cancel(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

// Refund provides a refund for the customer.
func (sh *SubscriptionHandler) Refund(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
