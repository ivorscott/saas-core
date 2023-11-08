// Package handler manages the presentation layer for handling incoming requests.
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/devpies/saas-core/internal/subscription/model"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v72"
	"go.uber.org/zap"
)

type subscriptionService interface {
	Refund(ctx context.Context) error
	Cancel(ctx context.Context) error
	Save(ctx context.Context, ns model.NewSubscription, now time.Time) (model.Subscription, error)
	GetOne(ctx context.Context, id string) (model.Subscription, error)
	BillingInfo(ctx context.Context, tenantID string) (model.BillingInfo, error)
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

// Create sets up a new subscription for the customer.
func (sh *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var (
		payload model.NewStripePayload
		err     error
	)

	if err = web.Decode(r, &payload); err != nil {
		return err
	}

	stripeSubscriptionID, stripeCustomerID, err := sh.subscriptionService.SubscribeStripeCustomer(payload)
	if err != nil {
		return err
	}

	customer, err := sh.customerService.Save(
		r.Context(),
		model.NewCustomer{
			StripeCustomerID: stripeCustomerID,
			FirstName:        payload.FirstName,
			LastName:         payload.LastName,
			Email:            payload.Email,
		},
		time.Now(),
	)
	if err != nil {
		return err
	}

	transaction, err := sh.transactionService.Save(
		r.Context(),
		model.NewTransaction{
			Amount:               payload.Amount,
			Currency:             payload.Currency,
			LastFour:             payload.LastFour,
			StatusID:             model.TransactionStatusCleared,
			ExpirationMonth:      payload.ExpirationMonth,
			ExpirationYear:       payload.ExpirationYear,
			PaymentMethod:        payload.PaymentMethod,
			StripeSubscriptionID: stripeSubscriptionID,
		},
		time.Now(),
	)
	if err != nil {
		return err
	}

	_, err = sh.subscriptionService.Save(
		r.Context(),
		model.NewSubscription{
			Plan:          payload.Plan,
			TransactionID: transaction.ID,
			StatusID:      model.SubscriptionStatusCleared,
			Amount:        payload.Amount,
			CustomerID:    customer.ID,
		},
		time.Now(),
	)
	if err != nil {
		return err
	}

	resp := struct {
		StripeSubscriptionID string `json:"stripeSubscriptionId"`
	}{
		StripeSubscriptionID: stripeSubscriptionID,
	}

	return web.Respond(r.Context(), w, resp, http.StatusOK)
}

// BillingInfo aggregates various stripe resource objects into a billing summary.
func (sh *SubscriptionHandler) BillingInfo(w http.ResponseWriter, r *http.Request) error {
	var (
		tenantID = chi.URLParam(r, "tenantID")
		info     model.BillingInfo
		err      error
	)

	info, err = sh.subscriptionService.BillingInfo(r.Context(), tenantID)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, info, http.StatusOK)
}

// GetAll retrieves all subscriptions.
func (sh *SubscriptionHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// GetOne retrieves a specific subscription by id.
func (sh *SubscriptionHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	var (
		s        model.Subscription
		tenantID = chi.URLParam(r, "tenantID")
		err      error
	)

	if _, err = uuid.Parse(tenantID); err != nil {
		return web.NewRequestError(err, http.StatusBadRequest)
	}

	s, err = sh.subscriptionService.GetOne(r.Context(), tenantID)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, s, http.StatusOK)
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
func (sh *SubscriptionHandler) Cancel(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Refund provides a refund for the customer.
func (sh *SubscriptionHandler) Refund(w http.ResponseWriter, r *http.Request) error {
	return nil
}
