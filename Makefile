include .env

.DEFAULT_GOAL := help

admin-client: ;@ ## Run admin frontend with live reload.
	@CompileDaemon \
	-build="go build -o ./bin/admin-client ./cmd/admin-client" \
	-command="./bin/admin-client \
	--web-frontend-port=${ADMIN_WEB_FRONTEND_PORT} \
	--web-backend-port=${ADMIN_WEB_BACKEND_PORT}" \
	-include="*.gohtml" \
	-log-prefix=false
.PHONY: admin-client

admin-api: ;@ ## Run admin backend with live reload.
	@CompileDaemon \
	-build="go build -o ./bin/admin-api ./cmd/admin-api" \
	-command="./bin/admin-api \
	--web-frontend-port=${ADMIN_WEB_FRONTEND_PORT} \
	--web-backend-port=${ADMIN_WEB_BACKEND_PORT}" \
	-log-prefix=false

.PHONY: admin-api

lint: ;@ ## Run linter.
	@golangci-lint run

help:
	@grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'
.PHONY: help