package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/devpies/saas-core/pkg/lambda/updatestatus/repository"
)

func handler(ctx context.Context, event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	client := repository.NewDynamoRepository(ctx, event.Region)

	id, ok := event.Request.UserAttributes["custom:tenant-id"]
	if !ok {
		return event, fmt.Errorf("error: missing tenant-id")
	}

	fmt.Printf("%s", id)

	out, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String("local-tenants"),
		Key: map[string]types.AttributeValue{
			"tenantId": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression: aws.String("set #S = :status"),
		ExpressionAttributeNames: map[string]string{
			"#S": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status": &types.AttributeValueMemberS{Value: "confirmed"},
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	})
	if err != nil {
		return event, fmt.Errorf("error: %w", err)
	}
	fmt.Printf("%v", out.Attributes)

	return event, nil
}

func main() {
	lambda.Start(handler)
}
