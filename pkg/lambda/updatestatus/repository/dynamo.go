package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
