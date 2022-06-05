#!/bin/bash

# Create Local tables
aws dynamodb create-table --table-name "auth-info" \
    --attribute-definitions \
        AttributeName=tenant_path,AttributeType=S \
    --key-schema \
        AttributeName=tenant_path,KeyType=HASH \
    --endpoint-url http://localhost:30008 \
    --provisioned-throughput \
        ReadCapacityUnits=5,WriteCapacityUnits=5
aws dynamodb create-table --table-name "tenants" \
    --attribute-definitions \
        AttributeName=tenant_id,AttributeType=S \
    --key-schema \
        AttributeName=tenant_id,KeyType=HASH \
    --endpoint-url http://localhost:30008 \
    --provisioned-throughput \
        ReadCapacityUnits=5,WriteCapacityUnits=5
aws dynamodb create-table --table-name "tenant-config" \
    --attribute-definitions \
        AttributeName=tenant_name,AttributeType=S \
    --key-schema \
        AttributeName=tenant_name,KeyType=HASH \
    --endpoint-url http://localhost:30008 \
    --provisioned-throughput \
        ReadCapacityUnits=5,WriteCapacityUnits=5

