package service

import (
	"context"

	"github.com/devpies/saas-core/internal/registration/model"

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

// PublishTenantMessages publishes messages in response to a new tenant being onboarded.
func (rs *RegistrationService) PublishTenantMessages(ctx context.Context, tenant model.NewTenant) error {
	return nil
}
