package repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/devpies/saas-core/pkg/lambda/modifytoken/model"
)

// DynamoRepository manages data access to DynamoDB.
type DynamoRepository struct {
	*dynamodb.Client
}

// NewDynamoRepository returns a new Repository for DynamoDB.
func NewDynamoRepository(ctx context.Context, region string) *DynamoRepository {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	return &DynamoRepository{dynamodb.NewFromConfig(cfg)}
}

func (r *DynamoRepository) LookupTenantKeys(ctx context.Context, userID string) ([]map[string]types.AttributeValue, error) {
	var (
		filter     expression.ConditionBuilder
		projection expression.ProjectionBuilder
		expr       expression.Expression
		err        error
	)

	filter = expression.Name("userId").Equal(expression.Value(userID))
	projection = expression.NamesList(expression.Name("tenantId"))
	expr, err = expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		return nil, fmt.Errorf("error building expression: %w", err)
	}

	out, err := r.Scan(ctx, &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String("local-connections"),
	})
	if err != nil {
		return nil, fmt.Errorf("error performing scan: %w", err)
	}

	return out.Items, nil
}

func (r *DynamoRepository) FindTenantConnections(ctx context.Context, userID string) (model.TenantConnectionMap, error) {
	tenantKeys, err := r.LookupTenantKeys(ctx, userID)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Tenant Keys: %+v \n", tenantKeys)

	out, err := r.BatchGetItem(ctx, &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			"local-tenants": {Keys: tenantKeys},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error performing batch get: %w", err)
	}

	tenants, err := model.NewTenants(out.Responses["local-tenants"])
	if err != nil {
		return nil, err
	}
	fmt.Printf("Tenants: %+v \n", tenants)
	connectionMap := model.NewTenantConnectionMap(tenants)
	fmt.Printf("Tenant Connection Map: %+v \n", connectionMap)
	return connectionMap, nil
}
