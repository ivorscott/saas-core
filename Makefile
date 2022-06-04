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

admin-end:	;@ ## Run end-to-end admin tests with Cypress.
	@cypress run --project e2e/admin/
.PHONY: admin-end

admin-test: admin-mock	;@ ## Run admin tests. Add " -- -v" for verbosity.
	go test $(val) -cover ./internal/admin/...
.PHONY: admin-test

admin-mock: ;@ ## Generate admin mocks.
	go generate ./internal/admin/...
.PHONY: admin-mock

admin-db: ;@ ## Enter admin database.
	pgcli postgres://$(ADMIN_POSTGRES_USER):$(ADMIN_POSTGRES_PASSWORD)@$(ADMIN_POSTGRES_HOST):$(ADMIN_POSTGRES_PORT)/$(ADMIN_POSTGRES_DB)
.PHONY: admin-db

admin-db-gen: ;@ ## Generate migration files. Required <name> argument.
	@migrate create -ext sql -dir ./internal/adminapi/res/migrations $(val)
.PHONY: admin-db-gen

admin-db-migrate: ;@ ## Migrate admin database. Optional <num> argument.
	@migrate -path ./internal/adminapi/res/migrations -verbose -database postgres://$(ADMIN_POSTGRES_USER):$(ADMIN_POSTGRES_PASSWORD)@$(ADMIN_POSTGRES_HOST):$(ADMIN_POSTGRES_PORT)/$(ADMIN_POSTGRES_DB)?sslmode=disable up $(val)
.PHONY: admin-db-migrate

admin-db-version: ;@ ## Print migration version for admin database.
	@migrate -path ./internal/adminapi/res/migrations -verbose -database postgres://$(ADMIN_POSTGRES_USER):$(ADMIN_POSTGRES_PASSWORD)@$(ADMIN_POSTGRES_HOST):$(ADMIN_POSTGRES_PORT)/$(ADMIN_POSTGRES_DB)?sslmode=disable up $(val)
.PHONY: admin-db-version

admin-db-rollback: ;@ ## Rollback admin database. Optional <num> argument.
	@migrate -path ./internal/adminapi/res/migrations -verbose -database postgres://$(ADMIN_POSTGRES_USER):$(ADMIN_POSTGRES_PASSWORD)@$(ADMIN_POSTGRES_HOST):$(ADMIN_POSTGRES_PORT)/$(ADMIN_POSTGRES_DB)?sslmode=disable down $(val)
.PHONY: admin-db-rollback

lint: ;@ ## Run linter.
	@golangci-lint run
.PHONY: lint

help:
	@grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'
.PHONY: help

# http://bit.ly/37TR1r2
ifeq ($(firstword $(MAKECMDGOALS)),$(filter $(firstword $(MAKECMDGOALS)),admin-test db-admin-gen db-admin-migrate db-admin-rollback))
  val := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(val):;@:)
endif
