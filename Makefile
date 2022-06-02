include .env

.DEFAULT_GOAL := help

admin: ;@ ## Run admin app with live reload.
	@CompileDaemon \
	-build="go build -o ./bin/admin ./cmd/admin" \
	-command="./bin/admin \
	--web-address=${ADMIN_WEB_ADDRESS} \
	--web-port=${ADMIN_WEB_PORT} \
	--cognito-app-client-id=${ADMIN_COGNITO_APP_CLIENT_ID} \
	--cognito-user-pool-client-id=${ADMIN_COGNITO_USER_POOL_CLIENT_ID} \
	--postgres-user=${ADMIN_POSTGRES_USER} \
	--postgres-password=${ADMIN_POSTGRES_PASSWORD} \
	--postgres-host=${ADMIN_POSTGRES_HOST} \
	--postgres-port=${ADMIN_POSTGRES_PORT} \
	--postgres-db=${ADMIN_POSTGRES_DB} \
	--postgres-disable-tls=true" \
	-include="*.gohtml" \
	-log-prefix=false
.PHONY: admin

# http://bit.ly/37TR1r2
ifeq ($(firstword $(MAKECMDGOALS)),$(filter $(firstword $(MAKECMDGOALS)),db-admin-gen db-admin-migrate db-admin-rollback))
  val := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(val):;@:)
endif

db-admin-gen: ;@ ## Generate migration files. Required <name> argument.
	@migrate create -ext sql -dir ./internal/adminapi/res/migrations $(val)
.PHONY: db-admin-gen

db-admin: ;@ ## Enter admin database.
	pgcli postgres://$(ADMIN_POSTGRES_USER):$(ADMIN_POSTGRES_PASSWORD)@$(ADMIN_POSTGRES_HOST):$(ADMIN_POSTGRES_PORT)/$(ADMIN_POSTGRES_DB)
.PHONY: db-admin

db-admin-migrate: ;@ ## Migrate admin database. Optional <num> argument.
	@migrate -path ./internal/adminapi/res/migrations -verbose -database postgres://$(ADMIN_POSTGRES_USER):$(ADMIN_POSTGRES_PASSWORD)@$(ADMIN_POSTGRES_HOST):$(ADMIN_POSTGRES_PORT)/$(ADMIN_POSTGRES_DB)?sslmode=disable up $(val)
.PHONY: db-admin-migrate

db-admin-version: ;@ ## Print migration version for admin database.
	@migrate -path ./internal/adminapi/res/migrations -verbose -database postgres://$(ADMIN_POSTGRES_USER):$(ADMIN_POSTGRES_PASSWORD)@$(ADMIN_POSTGRES_HOST):$(ADMIN_POSTGRES_PORT)/$(ADMIN_POSTGRES_DB)?sslmode=disable up $(val)
.PHONY: db-admin-version

db-admin-rollback: ;@ ## Rollback admin database. Optional <num> argument.
	@migrate -path ./internal/adminapi/res/migrations -verbose -database postgres://$(ADMIN_POSTGRES_USER):$(ADMIN_POSTGRES_PASSWORD)@$(ADMIN_POSTGRES_HOST):$(ADMIN_POSTGRES_PORT)/$(ADMIN_POSTGRES_DB)?sslmode=disable down $(val)
.PHONY: db-admin-rollback

lint: ;@ ## Run linter.
	@golangci-lint run
.PHONY: lint

help:
	@grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'
.PHONY: help