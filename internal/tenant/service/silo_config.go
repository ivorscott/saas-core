package service

import (
	"context"

	"github.com/devpies/saas-core/internal/tenant/model"
	"github.com/devpies/saas-core/pkg/msg"

	"go.uber.org/zap"
)

// SiloConfigService manages silo configuration for tenants.
type SiloConfigService struct {
	logger         *zap.Logger
	tenantRepo     tenantRepository
	siloConfigRepo siloConfigRepository
}

type siloConfigRepository interface {
	Insert(ctx context.Context, siloConfig model.NewSiloConfig) error
}

// NewSiloConfigService returns a new SiloConfigService.
func NewSiloConfigService(logger *zap.Logger, siloConfigRepo siloConfigRepository) *SiloConfigService {
	return &SiloConfigService{
		logger:         logger,
		siloConfigRepo: siloConfigRepo,
	}
}

// StoreConfigFromMessage stores tenant silo configuration from a message.
func (ts *SiloConfigService) StoreConfigFromMessage(ctx context.Context, message interface{}) error {
	m, err := msg.Bytes(message)
	if err != nil {
		return err
	}
	event, err := msg.UnmarshalTenantSiloedEvent(m)
	if err != nil {
		return err
	}
	config := model.NewSiloConfig(event.Data)
	err = ts.siloConfigRepo.Insert(ctx, config)
	if err != nil {
		return err
	}
	return nil
}
