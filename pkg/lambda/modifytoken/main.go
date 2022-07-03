// This package adds tenant connections as custom attributes before id token generation.
// https://github.com/aws/aws-lambda-go/blob/main/events/README_Cognito_UserPools_PreTokenGen.md
// https://docs.aws.amazon.com/cognito/latest/developerguide/user-pool-lambda-pre-token-generation.html
//
// Build instructions:
// 1. GOARCH=amd64 GOOS=linux go build -o main
// 2. zip function.zip main
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Tenant struct {
	TenantID    string `json:"id"`
	CompanyName string `json:"companyName"`
	Plan        string `json:"plan"`
}

// TenantConnection represents a valid tenant connection.
type TenantConnection struct {
	TenantID    string `json:"id"`
	CompanyName string `json:"companyName"`
	Plan        string `json:"plan"`
	Path        string `json:"path"`
}

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

func (r *DynamoRepository) FindTenants(ctx context.Context, userID string) ([]Tenant, error) {
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
	return newTenants(out.Responses["local-tenants"])
}

func newTenants(items []map[string]types.AttributeValue) ([]Tenant, error) {
	var tenants = make([]Tenant, 0, len(items))
	for _, v := range items {
		var item Tenant
		if err := attributevalue.UnmarshalMap(v, &item); err != nil {
			return nil, err
		}
		tenants = append(tenants, item)
	}
	return tenants, nil
}

func newTenantConnections(tenants []Tenant) []TenantConnection {
	var connections = make([]TenantConnection, 0, len(tenants))
	for _, tenant := range tenants {
		item := TenantConnection{
			TenantID:    tenant.TenantID,
			CompanyName: tenant.CompanyName,
			Path:        fmt.Sprintf("/%s", formatPath(tenant.CompanyName)),
			Plan:        tenant.Plan,
		}
		connections = append(connections, item)
	}
	return connections
}

func handler(ctx context.Context, event events.CognitoEventUserPoolsPreTokenGen) (events.CognitoEventUserPoolsPreTokenGen, error) {
	var err error

	client := NewDynamoRepository(ctx, event.Region)

	tenants, err := client.FindTenants(ctx, event.Request.UserAttributes["sub"])
	if err != nil {
		return event, err
	}
	fmt.Printf("Tenants: %+v \n", tenants)

	connections := newTenantConnections(tenants)
	fmt.Printf("Tenants Connections: %+v \n", connections)

	bytes, err := json.Marshal(&connections)
	if err != nil {
		return event, err
	}

	event.Response.ClaimsOverrideDetails.ClaimsToAddOrOverride = map[string]string{
		"tenant-connections": string(bytes),
	}
	return event, nil
}

func formatPath(company string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic("regex failed to compile" + err.Error())
	}
	return strings.ToLower(reg.ReplaceAllString(company, ""))
}

func main() {
	lambda.Start(handler)
}
