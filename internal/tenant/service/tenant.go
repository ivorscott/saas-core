package service

import (
	"context"

	"github.com/devpies/saas-core/internal/tenant/model"
	"github.com/devpies/saas-core/pkg/msg"

	"go.uber.org/zap"
)

type publisher interface {
	Publish(subject string, message []byte)
}

// TenantService manages tenant business operations.
type TenantService struct {
	logger     *zap.Logger
	js         publisher
	tenantRepo tenantRepository
}

type tenantRepository interface {
	Insert(ctx context.Context, tenant model.NewTenant) error
	SelectOne(ctx context.Context, tenantID string) (model.Tenant, error)
	SelectAll(ctx context.Context) ([]model.Tenant, error)
	Update(ctx context.Context, id string, tenant model.UpdateTenant) error
	Delete(ctx context.Context, tenantID string) error
}

// NewTenantService returns a new TenantService.
func NewTenantService(logger *zap.Logger, js publisher, tenantRepo tenantRepository) *TenantService {
	return &TenantService{
		logger:     logger,
		js:         js,
		tenantRepo: tenantRepo,
	}
}

// CreateTenantFromMessage creates a tenant from a message.
func (ts *TenantService) CreateTenantFromMessage(ctx context.Context, message interface{}) error {
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

//TODO: add status and created time
func newTenant(data msg.TenantRegisteredEventData) model.NewTenant {
	return model.NewTenant{
		ID:          data.TenantID,
		Email:       data.Email,
		FirstName:   data.FirstName,
		LastName:    data.LastName,
		CompanyName: data.Company,
		Plan:        data.Plan,
	}
}

// FindOne finds a single tenant.
func (ts *TenantService) FindOne(ctx context.Context, tenantID string) (model.Tenant, error) {
	var (
		tenant model.Tenant
		err    error
	)
	tenant, err = ts.tenantRepo.SelectOne(ctx, tenantID)
	if err != nil {
		return tenant, err
	}
	return tenant, nil
}

// FindAll finds all tenants.
func (ts *TenantService) FindAll(ctx context.Context) ([]model.Tenant, error) {
	var err error
	tenants, err := ts.tenantRepo.SelectAll(ctx)
	if err != nil {
		return nil, err
	}
	return tenants, nil
}

// Update updates a single tenant.
func (ts *TenantService) Update(ctx context.Context, id string, tenant model.UpdateTenant) error {
	var err error
	err = ts.tenantRepo.Update(ctx, id, tenant)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a tenant.
func (ts *TenantService) Delete(ctx context.Context, tenantID string) error {
	var err error
	err = ts.tenantRepo.Delete(ctx, tenantID)
	if err != nil {
		return err
	}
	return nil
}
