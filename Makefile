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
	--registration-service-address=${ADMIN_REGISTRATION_SERVICE_ADDRESS} \
	--registration-service-port=${ADMIN_REGISTRATION_SERVICE_PORT} \
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
	@pgcli postgres://$(ADMIN_POSTGRES_USER):$(ADMIN_POSTGRES_PASSWORD)@$(ADMIN_POSTGRES_HOST):$(ADMIN_POSTGRES_PORT)/$(ADMIN_POSTGRES_DB)
.PHONY: admin-db

admin-db-gen: ;@ ## Generate migration files. Required <name> argument.
	@migrate create -ext sql -dir ./internal/admin/res/migrations -seq $(val)
.PHONY: admin-db-gen

admin-db-migrate: ;@ ## Migrate admin database. Optional <num> argument.
	@migrate -path ./internal/admin/res/migrations -verbose -database postgres://$(ADMIN_POSTGRES_USER):$(ADMIN_POSTGRES_PASSWORD)@$(ADMIN_POSTGRES_HOST):$(ADMIN_POSTGRES_PORT)/$(ADMIN_POSTGRES_DB)?sslmode=disable up $(val)
.PHONY: admin-db-migrate

admin-db-version: ;@ ## Print migration version for admin database.
	@migrate -path ./internal/admin/res/migrations -verbose -database postgres://$(ADMIN_POSTGRES_USER):$(ADMIN_POSTGRES_PASSWORD)@$(ADMIN_POSTGRES_HOST):$(ADMIN_POSTGRES_PORT)/$(ADMIN_POSTGRES_DB)?sslmode=disable up $(val)
.PHONY: admin-db-version

admin-db-rollback: ;@ ## Rollback admin database. Optional <num> argument.
	@migrate -path ./internal/admin/res/migrations -verbose -database postgres://$(ADMIN_POSTGRES_USER):$(ADMIN_POSTGRES_PASSWORD)@$(ADMIN_POSTGRES_HOST):$(ADMIN_POSTGRES_PORT)/$(ADMIN_POSTGRES_DB)?sslmode=disable down $(val)
.PHONY: admin-db-rollback

registration: ;@ ## Run registration api with live reload.
	@CompileDaemon \
	-build="go build -o ./bin/registration ./cmd/registration" \
	-command="./bin/registration \
	--web-address=${REGISTRATION_WEB_ADDRESS} \
	--web-port=${REGISTRATION_WEB_PORT} \
	--cognito-admin-user-pool-client-id=${REGISTRATION_COGNITO_ADMIN_USER_POOL_CLIENT_ID} \
	--dynamodb-tenant-table=${REGISTRATION_DYNAMODB_TENANT_TABLE} \
	--dynamodb-auth-table=${REGISTRATION_DYNAMODB_AUTH_TABLE} \
	--dynamodb-config-table=${REGISTRATION_DYNAMODB_CONFIG_TABLE}" \
	-log-prefix=false
.PHONY: registration

registration-test: registration-mock	;@ ## Run registration tests. Add " -- -v" for verbosity.
	go test $(val) -cover ./internal/registration/...
.PHONY: registration-test

registration-mock: ;@ ## Generate registration mocks.
	go generate ./internal/registration/...
.PHONY: registration-mock

tenant: ;@ ## Run tenant api with live reload.
	@CompileDaemon \
	-build="go build -o ./bin/tenant ./cmd/tenant" \
	-command="./bin/tenant \
	--web-address=${TENANT_WEB_ADDRESS} \
	--web-port=${TENANT_WEB_PORT} \
	--cognito-user-pool-client-id=${TENANT_COGNITO_USER_POOL_CLIENT_ID} \
	--dynamodb-tenant-table=${TENANT_DYNAMODB_TENANT_TABLE} \
	--dynamodb-auth-table=${TENANT_DYNAMODB_AUTH_TABLE} \
	--dynamodb-config-table=${TENANT_DYNAMODB_CONFIG_TABLE}" \
	-log-prefix=false
.PHONY: tenant

tenant-test: tenant-mock	;@ ## Run tenant tests. Add " -- -v" for verbosity.
	go test $(val) -cover ./internal/tenant/...
.PHONY: tenant-test

tenant-mock: ;@ ## Generate tenant mocks.
	go generate ./internal/tenant/...
.PHONY: tenant-mock

identity: ;@ ## Run identity api with live reload.
	@CompileDaemon \
	-build="go build -o ./bin/identity ./cmd/identity" \
	-command="./bin/identity \
	--web-address=${IDENTITY_WEB_ADDRESS} \
	--web-port=${IDENTITY_WEB_PORT} \
	--cognito-user-pool-client-id=${IDENTITY_COGNITO_USER_POOL_CLIENT_ID}" \
	-log-prefix=false
.PHONY: identity

identity-test: identity-mock	;@ ## Run identity tests. Add " -- -v" for verbosity.
	go test $(val) -cover ./internal/identity/...
.PHONY: identity-test

identity-mock: ;@ ## Generate identity mocks.
	go generate ./internal/identity/...
.PHONY: identity-mock

project: ;@ ## Run project api with live reload.
	@CompileDaemon \
	-build="go build -o ./bin/project ./cmd/project" \
	-command="./bin/project \
	--web-address=${PROJECT_WEB_ADDRESS} \
	--web-port=${PROJECT_WEB_PORT} \
	--cognito-shared-user-pool-client-id=${PROJECT_COGNITO_SHARED_USER_POOL_CLIENT_ID} \
	--db-user=${PROJECT_DB_USER} \
	--db-password=${PROJECT_DB_PASSWORD} \
	--db-host=${PROJECT_DB_HOST} \
	--db-port=${PROJECT_DB_PORT} \
	--db-name=${PROJECT_DB_NAME} \
	--db-disable-tls=true" \
	-log-prefix=false
.PHONY: project

project-db: ;@ ## Enter project database.
	@pgcli postgres://$(PROJECT_DB_USER):$(PROJECT_DB_PASSWORD)@$(PROJECT_DB_HOST):$(PROJECT_DB_PORT)/$(PROJECT_DB_NAME)
.PHONY: project-db

project-db-gen: ;@ ## Generate migration files. Required <name> argument.
	@migrate create -ext sql -dir ./internal/project/res/migrations -seq $(val)
.PHONY: project-db-gen

project-db-migrate: ;@ ## Migrate project database. Optional <num> argument.
	@migrate -path ./internal/project/res/migrations -verbose -database postgres://$(PROJECT_DB_USER):$(PROJECT_DB_PASSWORD)@$(PROJECT_DB_HOST):$(PROJECT_DB_PORT)/$(PROJECT_DB_NAME)?sslmode=disable up $(val)
.PHONY: project-db-migrate

project-db-version: ;@ ## Print migration version for project database.
	@migrate -path ./internal/project/res/migrations -verbose -database postgres://$(PROJECT_DB_USER):$(PROJECT_DB_PASSWORD)@$(PROJECT_DB_HOST):$(PROJECT_DB_PORT)/$(PROJECT_DB_NAME)?sslmode=disable up $(val)
.PHONY: project-db-version

project-db-rollback: ;@ ## Rollback project database. Optional <num> argument.
	@migrate -path ./internal/project/res/migrations -verbose -database postgres://$(PROJECT_DB_USER):$(PROJECT_DB_PASSWORD)@$(PROJECT_DB_HOST):$(PROJECT_DB_PORT)/$(PROJECT_DB_NAME)?sslmode=disable down $(val)
.PHONY: project-db-rollback

user: ;@ ## Run user api with live reload.
	@CompileDaemon \
	-build="go build -o ./bin/user ./cmd/user" \
	-command="./bin/user \
	--web-address=${USER_WEB_ADDRESS} \
	--web-port=${USER_WEB_PORT} \
	--cognito-shared-user-pool-client-id=${USER_COGNITO_SHARED_USER_POOL_CLIENT_ID} \
	--db-user=${USER_DB_USER} \
	--db-password=${USER_DB_PASSWORD} \
	--db-host=${USER_DB_HOST} \
	--db-port=${USER_DB_PORT} \
	--db-name=${USER_DB_NAME} \
	--db-disable-tls=true" \
	-log-prefix=false
.PHONY: user

user-db: ;@ ## Enter user database.
	@pgcli postgres://$(USER_DB_USER):$(USER_DB_PASSWORD)@$(USER_DB_HOST):$(USER_DB_PORT)/$(USER_DB_NAME)
.PHONY: user-db

user-db-gen: ;@ ## Generate migration files. Required <name> argument.
	@migrate create -ext sql -dir ./internal/user/res/migrations -seq $(val)
.PHONY: user-db-gen

user-db-migrate: ;@ ## Migrate user database. Optional <num> argument.
	@migrate -path ./internal/user/res/migrations -verbose -database postgres://$(USER_DB_USER):$(USER_DB_PASSWORD)@$(USER_DB_HOST):$(USER_DB_PORT)/$(USER_DB_NAME)?sslmode=disable up $(val)
.PHONY: user-db-migrate

user-db-version: ;@ ## Print migration version for user database.
	@migrate -path ./internal/user/res/migrations -verbose -database postgres://$(USER_DB_USER):$(USER_DB_PASSWORD)@$(USER_DB_HOST):$(USER_DB_PORT)/$(USER_DB_NAME)?sslmode=disable up $(val)
.PHONY: user-db-version

user-db-rollback: ;@ ## Rollback user database. Optional <num> argument.
	@migrate -path ./internal/user/res/migrations -verbose -database postgres://$(USER_DB_USER):$(USER_DB_PASSWORD)@$(USER_DB_HOST):$(USER_DB_PORT)/$(USER_DB_NAME)?sslmode=disable down $(val)
.PHONY: user-db-rollback

tables:	;@ ## List Dynamodb tables.
	@aws dynamodb list-tables --endpoint-url http://localhost:30008
.PHONY: tables

lint: ;@ ## Run linter.
	@golangci-lint run
.PHONY: lint

help:
	@cat ./setup.txt
	@grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
.PHONY: help

# http://bit.ly/37TR1r2
ifeq ($(firstword $(MAKECMDGOALS)),$(filter $(firstword $(MAKECMDGOALS)),admin-test admin-db-gen admin-db-migrate admin-db-rollback registration-test project-test project-db-gen project-db-migrate project-db-rollback user-test user-db-gen user-db-migrate user-db-rollback))
  val := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(val):;@:)
endif
