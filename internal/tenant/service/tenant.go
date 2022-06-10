package service

import (
	"context"

	"github.com/devpies/saas-core/internal/tenant/model"
	"github.com/devpies/saas-core/pkg/msg"

	"go.uber.org/zap"
)

// TenantService manages tenant business operations.
type TenantService struct {
	logger         *zap.Logger
	tenantRepo     tenantRepository
	siloConfigRepo siloConfigRepository
	authInfoRepo   authInfoRepository
}

type tenantRepository interface {
	Insert(ctx context.Context, tenant model.NewTenant) error
	SelectOne(ctx context.Context, tenantID string) (model.Tenant, error)
	SelectAll(ctx context.Context) ([]model.Tenant, error)
	Update(ctx context.Context, tenant model.UpdateTenant) error
	Delete(ctx context.Context, tenantID string) error
}

type siloConfigRepository interface {
	Insert(ctx context.Context, siloConfig model.NewSiloConfig) error
}

type authInfoRepository interface {
	SelectAuthInfo(ctx context.Context, path string) (model.AuthInfo, error)
}

// NewTenantService returns a new TenantService.
func NewTenantService(logger *zap.Logger, tenantRepo tenantRepository, siloConfigRepo siloConfigRepository, authInfoRepo authInfoRepository) *TenantService {
	return &TenantService{
		logger:         logger,
		tenantRepo:     tenantRepo,
		siloConfigRepo: siloConfigRepo,
		authInfoRepo:   authInfoRepo,
	}
}

// CreateFromMessage creates a tenant from a message.
func (ts *TenantService) CreateFromMessage(ctx context.Context, message interface{}) error {
	m, err := msg.Bytes(message)
	if err != nil {
		return err
	}

	event, err := msg.UnmarshalTenantRegisteredEvent(m)
	if err != nil {
		return err
	}

	tenant := newTenant(event.Data)

	err = ts.tenantRepo.Insert(ctx, tenant)
	if err != nil {
		return err
	}

	return nil
}

// StoreConfigFromMessage stores tenant silo configuration from a message.
func (ts *TenantService) StoreConfigFromMessage(ctx context.Context, message interface{}) error {
	m, err := msg.Bytes(message)
	if err != nil {
		return err
	}

	event, err := msg.UnmarshalTenantSiloedEvent(m)
	if err != nil {
		return err
	}

	config := newSiloConfig(event.Data)

	err = ts.siloConfigRepo.Insert(ctx, config)
	if err != nil {
		return err
	}

	return nil
}

func newTenant(data msg.TenantRegisteredEventData) model.NewTenant {
	return model.NewTenant{
		ID:       data.ID,
		Email:    data.Email,
		FullName: data.FullName,
		Company:  data.Company,
		Plan:     data.Plan,
	}
}

func newSiloConfig(data msg.TenantSiloedEventData) model.NewSiloConfig {
	return model.NewSiloConfig{
		TenantName:       data.TenantName,
		UserPoolID:       data.UserPoolID,
		AppClientID:      data.AppClientID,
		DeploymentStatus: data.DeploymentStatus,
	}
}

// AddConfiguration adds tenant configuration for premium tenants.
func (ts *TenantService) AddConfiguration(ctx context.Context, tenantConfig model.NewTenantConfig) error {
	return nil
}

// GetAuthInfo gets the tenant authentication information.
func (ts *TenantService) GetAuthInfo(ctx context.Context, referer string) (model.AuthInfo, error) {
	var authInfo model.AuthInfo
	return authInfo, nil
}

// FindOne finds a single tenant.
func (ts *TenantService) FindOne(ctx context.Context, tenantID string) (model.Tenant, error) {
	var tenant model.Tenant
	return tenant, nil
}

// FindAll finds all tenants.
func (ts *TenantService) FindAll(ctx context.Context) ([]model.Tenant, error) {
	return nil, nil
}

// Update updates a single tenant.
func (ts *TenantService) Update(ctx context.Context, tenant model.UpdateTenant) error {
	return nil
}

// Delete removes a tenant.
func (ts *TenantService) Delete(ctx context.Context, tenantID string) error {
	return nil
}

// getPath parses the request URI and retrieves the base path. The base path is either "app" or the shortened tenant name.
func (ts *TenantService) getPath(referer string) error {
	return nil
}
