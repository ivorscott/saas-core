package service

import (
	"context"

	"github.com/devpies/saas-core/internal/tenant/config"
	"github.com/devpies/saas-core/internal/tenant/model"

	"go.uber.org/zap"
)

// TenantService manages tenant business operations.
type TenantService struct {
	logger           *zap.Logger
	config           config.Config
	tenantRepo       tenantRepository
	tenantConfigRepo tenantConfigRepository
	authInfoRepo     authInfoRepository
}

type tenantRepository interface {
	Insert(ctx context.Context, tenant model.NewTenant) error
	SelectOne(ctx context.Context, tenantID string) (model.Tenant, error)
	SelectAll(ctx context.Context) ([]model.Tenant, error)
	Update(ctx context.Context, tenant model.UpdateTenant) error
	Delete(ctx context.Context, tenantID string) error
}

type tenantConfigRepository interface {
	InsertConfiguration(ctx context.Context, tenantConfig model.NewTenantConfig) error
}

type authInfoRepository interface {
	SelectAuthInfo(ctx context.Context, path string) (model.AuthInfo, error)
}

// NewTenantService returns a new TenantService.
func NewTenantService(logger *zap.Logger, config config.Config, tenantRepo tenantRepository, tenantConfigRepo tenantConfigRepository, authInfoRepo authInfoRepository) *TenantService {
	return &TenantService{
		logger:           logger,
		config:           config,
		tenantRepo:       tenantRepo,
		tenantConfigRepo: tenantConfigRepo,
		authInfoRepo:     authInfoRepo,
	}
}

// Create creates a tenant.
func (ts *TenantService) Create(ctx context.Context, tenant model.NewTenant) error {
	return nil
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
