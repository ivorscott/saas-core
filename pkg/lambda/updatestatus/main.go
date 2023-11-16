package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/devpies/saas-core/pkg/lambda/updatestatus/repository"
)

func handler(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	var err error

	client := repository.NewDynamoRepository(ctx, event.Region)

	id, ok := event.Request.UserAttributes["custom:tenant-id"]
	if !ok {
		return event, nil
	}

	_, err = client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String("local-tenants"),
		Key: map[string]types.AttributeValue{
			"tenantId": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression: aws.String("set status = :status"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberS{Value: "confirmed"},
		},
	})
	if err != nil {
		return event, nil
	}

	return event, nil
}

func main() {
	lambda.Start(handler)
}
