#!/bin/bash

# Create Local tables
aws dynamodb create-table --table-name "auth-info" \
    --attribute-definitions \
        AttributeName=tenantPath,AttributeType=S \
    --key-schema \
        AttributeName=tenantPath,KeyType=HASH \
    --endpoint-url http://localhost:30008 \
    --provisioned-throughput \
        ReadCapacityUnits=5,WriteCapacityUnits=5
aws dynamodb create-table --table-name "tenants" \
    --attribute-definitions \
        AttributeName=tenantId,AttributeType=S \
    --key-schema \
        AttributeName=tenantId,KeyType=HASH \
    --endpoint-url http://localhost:30008 \
    --provisioned-throughput \
        ReadCapacityUnits=5,WriteCapacityUnits=5
aws dynamodb create-table --table-name "silo-config" \
    --attribute-definitions \
        AttributeName=tenantName,AttributeType=S \
    --key-schema \
        AttributeName=tenantName,KeyType=HASH \
    --endpoint-url http://localhost:30008 \
    --provisioned-throughput \
        ReadCapacityUnits=5,WriteCapacityUnits=5

