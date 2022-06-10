package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/devpies/saas-core/internal/tenant/model"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// SiloConfigRepository manages data access to silo configuration.
type SiloConfigRepository struct {
	client *dynamodb.Client
	table  string
}

// NewSiloConfigRepository returns a new SiloConfigRepository.
func NewSiloConfigRepository(client *dynamodb.Client, table string) *SiloConfigRepository {
	return &SiloConfigRepository{
		client: client,
		table:  table,
	}
}

// Insert stores silo configuration for premium tenants.
func (sr *SiloConfigRepository) Insert(ctx context.Context, siloConfig model.NewSiloConfig) error {
	input := dynamodb.PutItemInput{
		TableName: aws.String(sr.table),
		Item: map[string]types.AttributeValue{
			"tenant_name":       &types.AttributeValueMemberS{Value: siloConfig.TenantName},
			"user_pool_id":      &types.AttributeValueMemberS{Value: siloConfig.UserPoolID},
			"app_client_id":     &types.AttributeValueMemberS{Value: siloConfig.AppClientID},
			"deployment_status": &types.AttributeValueMemberS{Value: siloConfig.DeploymentStatus},
		},
	}
	_, err := sr.client.PutItem(ctx, &input)
	if err != nil {
		return err
	}
	return nil
}
