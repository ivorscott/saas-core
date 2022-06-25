#!/bin/bash

# Delete local tables
aws dynamodb delete-table --table-name "auth-info" \
    --endpoint-url http://localhost:30008
aws dynamodb delete-table --table-name "tenants" \
    --endpoint-url http://localhost:30008
aws dynamodb delete-table --table-name "silo-config" \
    --endpoint-url http://localhost:30008

