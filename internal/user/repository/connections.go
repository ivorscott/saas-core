package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/devpies/saas-core/internal/user/model"
)

// ConnectionRepository manages data access to user tenant connections.
type ConnectionRepository struct {
	client *dynamodb.Client
	table  string
}

// NewConnectionRepository returns a new ConnectionRepository.
func NewConnectionRepository(client *dynamodb.Client, table string) *ConnectionRepository {
	return &ConnectionRepository{
		client: client,
		table:  table,
	}
}

// Insert stores a new tenant connection.
func (tr *ConnectionRepository) Insert(ctx context.Context, connection model.NewConnection) error {
	input := dynamodb.PutItemInput{
		TableName: aws.String(tr.table),
		Item: map[string]types.AttributeValue{
			"userId":   &types.AttributeValueMemberS{Value: connection.UserID},
			"tenantId": &types.AttributeValueMemberS{Value: connection.TenantID},
		},
	}
	_, err := tr.client.PutItem(ctx, &input)
	if err != nil {
		return err
	}
	return nil
}