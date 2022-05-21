include .env

.DEFAULT_GOAL := help

admin-client: ;@ ## Run admin frontend with live reload.
	@CompileDaemon \
	-build="go build -o ./bin/adminclient ./cmd/adminclient" \
	-command="./bin/adminclient \
	--web-backend=${ADMIN_WEB_BACKEND} \
	--web-backend-port=${ADMIN_WEB_BACKEND_PORT} \
	--web-frontend-port=${ADMIN_WEB_FRONTEND_PORT}" \
	-include="*.gohtml" \
	-log-prefix=false
.PHONY: admin-client

admin-api: ;@ ## Run admin backend with live reload.
	@CompileDaemon \
	-build="go build -o ./bin/adminapi ./cmd/adminapi" \
	-command="./bin/adminapi \
	--web-backend=${ADMIN_WEB_BACKEND} \
	--web-backend-port=${ADMIN_WEB_BACKEND_PORT} \
	--web-frontend-port=${ADMIN_WEB_FRONTEND_PORT} \
	--cognito-app-client-id=${ADMIN_COGNITO_APP_CLIENT_ID} \
	--cognito-user-pool-client-id=${ADMIN_COGNITO_USER_POOL_CLIENT_ID}" \
	-log-prefix=false

.PHONY: admin-api

db-sessions:  ;@ ## Enter user sessions database.
	pgcli postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)

lint: ;@ ## Run linter.
	@golangci-lint run

help:
	@grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'
.PHONY: help