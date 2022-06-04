package service

import (
	"context"

	"github.com/devpies/saas-core/internal/admin/model"

	"go.uber.org/zap"
)

// RegistrationService is responsible for managing tenant registration.
type RegistrationService struct {
	logger *zap.Logger
}

// NewRegistrationService returns a new registration service.
func NewRegistrationService(logger *zap.Logger) *RegistrationService {
	return &RegistrationService{
		logger: logger,
	}
}

// RegisterTenant sends new tenant to registration microservice.
func (rs *RegistrationService) RegisterTenant(ctx context.Context, tenant model.NewTenant) error {
	// TODO: Call registration microservice.
	return nil
}
