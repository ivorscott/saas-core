// Package handler manages the presentation layer for handling incoming requests.
package handler

import (
	"context"
	"net/http"

	"github.com/devpies/saas-core/internal/subscription/model"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/go-chi/chi/v5"
	"github.com/stripe/stripe-go/v72"
	"go.uber.org/zap"
)

type subscriptionService interface {
	Refund(ctx context.Context, subID string) error
	Cancel(subID string) error
	SubscriptionInfo(ctx context.Context, tenantID string) (model.SubscriptionInfo, error)
	CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error)
	SubscribeStripeCustomer(ctx context.Context, nsp model.NewStripePayload) error
}

// SubscriptionHandler handles subscription related requests.
type SubscriptionHandler struct {
	logger              *zap.Logger
	subscriptionService subscriptionService
}

// NewSubscriptionHandler returns a SubscriptionHandler.
func NewSubscriptionHandler(
	logger *zap.Logger,
	subscriptionService subscriptionService,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		logger:              logger,
		subscriptionService: subscriptionService,
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

	err = sh.subscriptionService.SubscribeStripeCustomer(r.Context(), payload)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

// SubscriptionInfo aggregates various stripe resource objects into a subscription summary.
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
func (sh *SubscriptionHandler) Cancel(w http.ResponseWriter, r *http.Request) error {
	var (
		subID = chi.URLParam(r, "subID")
		err   error
	)

	err = sh.subscriptionService.Cancel(subID)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

// Refund provides a refund for the customer.
func (sh *SubscriptionHandler) Refund(w http.ResponseWriter, r *http.Request) error {
	var (
		subID = chi.URLParam(r, "subID")
		err   error
	)

	err = sh.subscriptionService.Refund(r.Context(), subID)
	if err != nil {
		return err
	}
	return web.Respond(r.Context(), w, nil, http.StatusOK)
}
