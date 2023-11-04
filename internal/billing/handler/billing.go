// Package handler manages the presentation layer for handling incoming requests.
package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/devpies/saas-core/internal/billing/model"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/google/uuid"

	"go.uber.org/zap"
)

type subscriptionService interface {
	Refund(ctx context.Context) error
	Cancel(ctx context.Context) error
	Save(ctx context.Context, ns model.NewSubscription, now time.Time) (model.Subscription, error)
	GetAll(ctx context.Context) ([]model.Subscription, error)
	GetOne(ctx context.Context, id string) (model.Subscription, error)
	SubscribeStripeCustomer(nsp model.NewStripePayload) (string, error)
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

	stripeSubscriptionID, err := sh.subscriptionService.SubscribeStripeCustomer(payload)
	if err != nil {
		return err
	}

	customer, err := sh.customerService.Save(
		r.Context(),
		model.NewCustomer{
			ID:        uuid.New().String(),
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
			Email:     payload.Email,
		},
		time.Now(),
	)
	if err != nil {
		return err
	}

	transaction, err := sh.transactionService.Save(
		r.Context(),
		model.NewTransaction{
			ID:                   uuid.New().String(),
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
			ID:            uuid.New().String(),
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

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

// GetAll retrieves all subscriptions.
func (sh *SubscriptionHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// GetOne retrieves a specific subscription by id.
func (sh *SubscriptionHandler) GetOne(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// GetPaymentIntent retrieves the paymentIntent from stripe.
func (sh *SubscriptionHandler) GetPaymentIntent(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Cancel cancels a paid subscription.
func (sh *SubscriptionHandler) Cancel(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Refund provides a refund for the customer.
func (sh *SubscriptionHandler) Refund(w http.ResponseWriter, r *http.Request) error {
	return nil
}
