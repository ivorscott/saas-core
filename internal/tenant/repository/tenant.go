package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/devpies/saas-core/internal/tenant/model"
)

// TenantRepository manages data access to system tenants.
type TenantRepository struct {
	client *dynamodb.Client
	table  string
}

// NewTenantRepository returns a new TenantRepository.
func NewTenantRepository(client *dynamodb.Client, table string) *TenantRepository {
	return &TenantRepository{
		client: client,
		table:  table,
	}
}

// Insert stores a new tenant.
func (tr *TenantRepository) Insert(ctx context.Context, tenant model.NewTenant) error {
	input := dynamodb.PutItemInput{
		Item: map[string]types.AttributeValue{
			"tenant_id": &types.AttributeValueMemberS{Value: tenant.ID},
			"email":     &types.AttributeValueMemberS{Value: tenant.Email},
			"fullname":  &types.AttributeValueMemberS{Value: tenant.FullName},
			"company":   &types.AttributeValueMemberS{Value: tenant.Company},
			"plan":      &types.AttributeValueMemberS{Value: tenant.Plan},
		},
		TableName: aws.String(tr.table),
	}
	_, err := tr.client.PutItem(ctx, &input)
	if err != nil {
		return err
	}
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
