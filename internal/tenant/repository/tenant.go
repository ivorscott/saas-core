package repository

import (
	"context"

	"github.com/devpies/saas-core/internal/tenant/model"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// TenantRepository manages data access to system tenants.
type TenantRepository struct {
	client *dynamodb.Client
}

// NewTenantRepository returns a new TenantRepository.
func NewTenantRepository(client *dynamodb.Client) *TenantRepository {
	return &TenantRepository{
		client: client,
	}
}

// Insert stores a new tenant.
func (tr *TenantRepository) Insert(ctx context.Context, tenant model.NewTenant) error {
	return nil
}

// SelectOne retrieves a single tenant.
func (tr *TenantRepository) SelectOne(ctx context.Context, tenantID string) (model.Tenant, error) {
	var tenant model.Tenant
	return tenant, nil
}

// SelectAll retrieves all tenants.
func (tr *TenantRepository) SelectAll(ctx context.Context) ([]model.Tenant, error) {
	var tenants []model.Tenant
	return tenants, nil
}

// Update updates a tenant.
func (tr *TenantRepository) Update(ctx context.Context, tenant model.UpdateTenant) error {
	return nil
}

// Delete removes a tenant.
func (tr *TenantRepository) Delete(ctx context.Context, tenantID string) error {
	return nil
}
