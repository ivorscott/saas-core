package db

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// NewProductionDynamoDBClient returns a production client for DynamoDB.
func NewProductionDynamoDBClient(ctx context.Context) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	return dynamodb.NewFromConfig(cfg)
}

// NewDevelopmentDynamoDBClient returns a development client for DynamoDB.
func NewDevelopmentDynamoDBClient(ctx context.Context, port string) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("eu-central-1"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: fmt.Sprintf("http://localhost:%s", port)}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)
	if err != nil {
		panic(err)
	}

	return dynamodb.NewFromConfig(cfg)
}
