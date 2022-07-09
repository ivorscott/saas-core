package repository

import (
	"context"
	"strings"

	"github.com/devpies/saas-core/internal/tenant/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
		TableName: aws.String(tr.table),
		Item: map[string]types.AttributeValue{
			"tenantId":    &types.AttributeValueMemberS{Value: tenant.ID},
			"email":       &types.AttributeValueMemberS{Value: tenant.Email},
			"firstName":   &types.AttributeValueMemberS{Value: tenant.FirstName},
			"lastName":    &types.AttributeValueMemberS{Value: tenant.LastName},
			"companyName": &types.AttributeValueMemberS{Value: tenant.CompanyName},
			"plan":        &types.AttributeValueMemberS{Value: tenant.Plan},
			"enabled":     &types.AttributeValueMemberBOOL{Value: true},
			"status":      &types.AttributeValueMemberS{Value: tenant.Status},
			"created":     &types.AttributeValueMemberS{Value: tenant.Created.String()},
		},
	}
	_, err := tr.client.PutItem(ctx, &input)
	if err != nil {
		return err
	}
	return nil
}

// SelectOne retrieves a single tenant.
func (tr *TenantRepository) SelectOne(ctx context.Context, tenantID string) (model.Tenant, error) {
	var (
		tenant model.Tenant
		err    error
	)

	input := dynamodb.GetItemInput{
		TableName: aws.String(tr.table),
		Key: map[string]types.AttributeValue{
			"tenantId": &types.AttributeValueMemberS{Value: tenantID},
		},
	}
	output, err := tr.client.GetItem(ctx, &input)
	if err != nil {
		return tenant, err
	}

	err = attributevalue.UnmarshalMap(output.Item, &tenant)
	if err != nil {
		return tenant, err
	}

	return tenant, nil
}

// SelectAll retrieves all tenants.
func (tr *TenantRepository) SelectAll(ctx context.Context) ([]model.Tenant, error) {
	out, err := tr.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(tr.table),
	})
	if err != nil {
		return nil, err
	}

	var tenants = make([]model.Tenant, 0, len(out.Items))
	for _, v := range out.Items {
		var item model.Tenant
		err = attributevalue.UnmarshalMap(v, &item)
		if err != nil {
			return nil, err
		}
		tenants = append(tenants, item)
	}
	return tenants, nil
}

// Update updates a tenant.
func (tr *TenantRepository) Update(ctx context.Context, id string, update model.UpdateTenant) error {
	var (
		updateExp = "set"
		av        = make(map[string]types.AttributeValue)
	)

	if update.Email != nil {
		updateExp = updateExp + " email = :email,"
		av[":email"] = &types.AttributeValueMemberS{Value: *update.Email}
	}

	if update.FirstName != nil {
		updateExp = updateExp + " firstName = :firstName,"
		av[":firstName"] = &types.AttributeValueMemberS{Value: *update.FirstName}
	}

	if update.LastName != nil {
		updateExp = updateExp + " lastName = :lastName,"
		av[":lastName"] = &types.AttributeValueMemberS{Value: *update.LastName}
	}

	if update.CompanyName != nil {
		updateExp = updateExp + " companyName = :companyName,"
		av[":companyName"] = &types.AttributeValueMemberS{Value: *update.CompanyName}
	}

	if update.Plan != nil {
		updateExp = updateExp + " plan = :plan,"
		av[":plan"] = &types.AttributeValueMemberS{Value: *update.Plan}
	}

	if update.Status != nil {
		updateExp = updateExp + " status = :status,"
		av[":status"] = &types.AttributeValueMemberS{Value: *update.Status}
	}

	updateExp = strings.TrimSuffix(updateExp, ",")

	_, err := tr.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(tr.table),
		Key: map[string]types.AttributeValue{
			"tenantId": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression:          aws.String(updateExp),
		ExpressionAttributeValues: av,
	})
	if err != nil {
		return nil
	}
	return nil
}

// Delete removes a tenant.
func (tr *TenantRepository) Delete(ctx context.Context, tenantID string) error {
	_, err := tr.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(tr.table),
		Key: map[string]types.AttributeValue{
			"tenantId": &types.AttributeValueMemberS{Value: tenantID},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
