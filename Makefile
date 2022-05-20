include .env

.DEFAULT_GOAL := help

admin-client: ;@ ## Run admin frontend with live reload.
	@CompileDaemon \
	-build="go build -o ./bin/admin-client ./cmd/admin-client" \
	-command="./bin/admin-client \
	--web-backend=${ADMIN_WEB_BACKEND} \
	--web-backend-port=${ADMIN_WEB_BACKEND_PORT} \
	--web-frontend-port=${ADMIN_WEB_FRONTEND_PORT}" \
	-include="*.gohtml" \
	-log-prefix=false
.PHONY: admin-client

admin-api: ;@ ## Run admin backend with live reload.
	@CompileDaemon \
	-build="go build -o ./bin/admin-api ./cmd/admin-api" \
	-command="./bin/admin-api \
	--web-backend=${ADMIN_WEB_BACKEND} \
	--web-backend-port=${ADMIN_WEB_BACKEND_PORT} \
	--web-frontend-port=${ADMIN_WEB_FRONTEND_PORT} \
	--cognito-app-client-id=${ADMIN_COGNITO_APP_CLIENT_ID} \
	--cognito-user-pool-client-id=${ADMIN_COGNITO_USER_POOL_CLIENT_ID}" \
	-log-prefix=false

.PHONY: admin-api

lint: ;@ ## Run linter.
	@golangci-lint run

help:
	@grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'
.PHONY: help