#!/bin/bash

# Create local dynamodb tables and store "shared user pool" info

SHARED_POOL=$1
SHARED_POOL_CLIENT=$2

aws dynamodb create-table --table-name "auth-info" \
    --attribute-definitions \
        AttributeName=tenantPath,AttributeType=S \
    --key-schema \
        AttributeName=tenantPath,KeyType=HASH \
    --provisioned-throughput \
        ReadCapacityUnits=5,WriteCapacityUnits=5

aws dynamodb put-item --table-name "auth-info" \
    --item \
        '{"tenantPath": {"S": "/app"}, "userPoolId": {"S": "'${SHARED_POOL}'"}, "userPoolClientId": {"S": "'${SHARED_POOL_CLIENT}'"}}'

aws dynamodb create-table --table-name "tenants" \
    --attribute-definitions \
        AttributeName=tenantId,AttributeType=S \
    --key-schema \
        AttributeName=tenantId,KeyType=HASH \
    --provisioned-throughput \
        ReadCapacityUnits=5,WriteCapacityUnits=5
aws dynamodb create-table --table-name "silo-config" \
    --attribute-definitions \
        AttributeName=tenantName,AttributeType=S \
    --key-schema \
        AttributeName=tenantName,KeyType=HASH \
    --provisioned-throughput \
        ReadCapacityUnits=5,WriteCapacityUnits=5

