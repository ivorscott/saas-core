package service

import (
	"context"
	"time"

	"github.com/devpies/saas-core/internal/billing/model"

	"go.uber.org/zap"
)

type customerRepository interface {
	SaveCustomer(ctx context.Context, nc model.NewCustomer, now time.Time) (model.Customer, error)
}

// CustomerService is responsible for managing customer related business logic.
type CustomerService struct {
	logger *zap.Logger
	repo   customerRepository
}

// NewCustomerService returns a new CustomerService.
func NewCustomerService(logger *zap.Logger, repo customerRepository) *CustomerService {
	return &CustomerService{
		logger: logger,
		repo:   repo,
	}
}

// Save persists the new customer details.
func (cs *CustomerService) Save(ctx context.Context, nc model.NewCustomer, now time.Time) (model.Customer, error) {
	var (
		c   model.Customer
		err error
	)

	c, err = cs.repo.SaveCustomer(ctx, nc, now)
	if err != nil {
		switch err {
		default:
			return c, err
		}
	}
	return c, nil
}
