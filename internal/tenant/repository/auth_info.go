package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/devpies/saas-core/internal/tenant/model"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// AuthInfoRepository manages data access to tenant authentication information.
type AuthInfoRepository struct {
	client *dynamodb.Client
	table  string
}

// NewAuthInfoRepository returns a new AuthInfoRepository.
func NewAuthInfoRepository(client *dynamodb.Client, table string) *AuthInfoRepository {
	return &AuthInfoRepository{
		client: client,
		table:  table,
	}
}

// SelectAuthInfo retrieves authentication information for a specific tenant.
func (ar *AuthInfoRepository) SelectAuthInfo(ctx context.Context, path string) (model.AuthInfo, error) {
	var authInfo model.AuthInfo
	input := dynamodb.GetItemInput{
		TableName: aws.String(ar.table),
		Key: map[string]types.AttributeValue{
			"tenantPath": &types.AttributeValueMemberS{Value: path},
		},
	}
	output, err := ar.client.GetItem(ctx, &input)
	if err != nil {
		return authInfo, err
	}
	err = attributevalue.UnmarshalMap(output.Item, &authInfo)
	if err != nil {
		return authInfo, err
	}
	return authInfo, nil
}
