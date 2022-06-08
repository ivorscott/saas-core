package repository

import (
	"context"

	"github.com/devpies/saas-core/internal/tenant/model"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// TenantConfigRepository manages data access to tenant configuration.
type TenantConfigRepository struct {
	client *dynamodb.Client
}

// NewTenantConfigRepository returns a new TenantConfigRepository.
func NewTenantConfigRepository(client *dynamodb.Client) *TenantConfigRepository {
	return &TenantConfigRepository{
		client: client,
	}
}

// InsertConfiguration stores premium tenant configuration.
func (cr *TenantConfigRepository) InsertConfiguration(ctx context.Context, tenantConfig model.NewTenantConfig) error {
	return nil
}
