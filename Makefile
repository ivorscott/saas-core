include .env

.DEFAULT_GOAL := help

admin: ;@ ## Run admin app with live reload.
	@CompileDaemon \
	-build="go build -o ./bin/admin ./cmd/admin" \
	-command="./bin/admin \
	--web-port=${ADMIN_WEB_PORT} \
	--cognito-user-pool-id=${ADMIN_USER_POOL_ID} \
	--cognito-user-pool-client-id=${ADMIN_USER_POOL_CLIENT_ID} \
	--cognito-shared-user-pool-id=${SHARED_USER_POOL_ID} \
	--cognito-region=${REGION} \
	--registration-service-port=${REGISTRATION_WEB_PORT} \
	--db-port=${ADMIN_DB_PORT} \
	--db-disable-tls=true" \
	-include="*.gohtml" \
	-log-prefix=false
.PHONY: admin

admin-test: admin-mock	;@ ## Run admin tests. Add " -- -v" for verbosity.
	go test $(val) -cover ./internal/admin/...
.PHONY: admin-test

admin-mock: ;@ ## Generate admin mocks.
	go generate ./internal/admin/...
.PHONY: admin-mock

admin-db: ;@ ## Enter admin database.
	@pgcli postgres://postgres:postgres@localhost:$(ADMIN_DB_PORT)/admin
.PHONY: admin-db

admin-db-gen: ;@ ## Generate migration files. Required <name> argument.
	@migrate create -ext sql -dir ./internal/admin/res/migrations -seq $(val)
.PHONY: admin-db-gen

admin-db-migrate: ;@ ## Migrate admin database. Optional <num> argument.
	@migrate -path ./internal/admin/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(ADMIN_DB_PORT)/admin?sslmode=disable up $(val)
.PHONY: admin-db-migrate

admin-db-version: ;@ ## Print migration version for admin database.
	@migrate -path ./internal/admin/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(ADMIN_DB_PORT)/admin?sslmode=disable version
.PHONY: admin-db-version

admin-db-rollback: ;@ ## Rollback admin database. Optional <num> argument.
	@migrate -path ./internal/admin/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(ADMIN_DB_PORT)/admin?sslmode=disable down $(val)
.PHONY: admin-db-rollback

admin-db-force: ;@ ## Force version on admin database. Optional <num> argument.
	@migrate -path ./internal/admin/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(ADMIN_DB_PORT)/admin?sslmode=disable force $(val)
.PHONY: admin-db-force

registration: ;@ ## Run registration api with live reload.
	@CompileDaemon \
	-build="go build -o ./bin/registration ./cmd/registration" \
	-command="./bin/registration \
	--web-port=${REGISTRATION_WEB_PORT} \
	--cognito-user-pool-id=${ADMIN_USER_POOL_ID} \
	--cognito-region=${REGION} \
	--dynamodb-tenant-table=${DYNAMODB_TENANT_TABLE} \
	--dynamodb-auth-table=${DYNAMODB_AUTH_TABLE} \
	--dynamodb-config-table=${DYNAMODB_CONFIG_TABLE}" \
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
	--web-port=${TENANT_WEB_PORT} \
	--cognito-user-pool-id=${ADMIN_USER_POOL_ID} \
	--cognito-shared-user-pool-id=${SHARED_USER_POOL_ID} \
	--cognito-region=${REGION} \
	--dynamodb-connection-table=${DYNAMODB_CONNECTION_TABLE} \
	--dynamodb-tenant-table=${DYNAMODB_TENANT_TABLE} \
	--dynamodb-auth-table=${DYNAMODB_AUTH_TABLE} \
	--dynamodb-config-table=${DYNAMODB_CONFIG_TABLE}" \
	-log-prefix=false
.PHONY: tenant

tenant-test: tenant-mock	;@ ## Run tenant tests. Add " -- -v" for verbosity.
	go test $(val) -cover ./internal/tenant/...
.PHONY: tenant-test

tenant-mock: ;@ ## Generate tenant mocks.
	go generate ./internal/tenant/...
.PHONY: tenant-mock

project: ;@ ## Run project api with live reload.
	@CompileDaemon \
	-build="go build -o ./bin/project ./cmd/project" \
	-command="./bin/project \
	--web-port=${PROJECT_WEB_PORT} \
	--cognito-user-pool-id=${SHARED_USER_POOL_ID} \
	--cognito-region=${REGION} \
	--db-port=${PROJECT_DB_PORT} \
	--db-disable-tls=true" \
	-log-prefix=false
.PHONY: project

project-db: ;@ ## Enter project database.
	@pgcli postgres://postgres:postgres@localhost:$(PROJECT_DB_PORT)/project
.PHONY: project-db

project-db-gen: ;@ ## Generate migration files. Required <name> argument.
	@migrate create -ext sql -dir ./internal/project/res/migrations -seq $(val)
.PHONY: project-db-gen

project-db-migrate: ;@ ## Migrate project database. Optional <num> argument.
	@migrate -path ./internal/project/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(PROJECT_DB_PORT)/project?sslmode=disable up $(val)
.PHONY: project-db-migrate

project-db-version: ;@ ## Print migration version for project database.
	@migrate -path ./internal/project/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(PROJECT_DB_PORT)/project?sslmode=disable version
.PHONY: project-db-version

project-db-rollback: ;@ ## Rollback project database. Optional <num> argument.
	@migrate -path ./internal/project/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(PROJECT_DB_PORT)/project?sslmode=disable down $(val)
.PHONY: project-db-rollback

project-db-force: ;@ ## Force version on project database. Optional <num> argument.
	@migrate -path ./internal/project/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(PROJECT_DB_PORT)/project?sslmode=disable force $(val)
.PHONY: project-db-force

user: ;@ ## Run user api with live reload.
	@CompileDaemon \
	-build="go build -o ./bin/user ./cmd/user" \
	-command="./bin/user \
	--web-port=${USER_WEB_PORT} \
	--cognito-shared-user-pool-id=${SHARED_USER_POOL_ID} \
	--cognito-region=${REGION} \
	--dynamodb-connection-table=${DYNAMODB_CONNECTION_TABLE} \
	--db-port=${USER_DB_PORT} \
	--db-disable-tls=true" \
	-log-prefix=false
.PHONY: user

user-db: ;@ ## Enter user database.
	@pgcli postgres://postgres:postgres@localhost:$(USER_DB_PORT)/user
.PHONY: user-db

user-db-gen: ;@ ## Generate migration files. Required <name> argument.
	@migrate create -ext sql -dir ./internal/user/res/migrations -seq $(val)
.PHONY: user-db-gen

user-db-migrate: ;@ ## Migrate user database. Optional <num> argument.
	@migrate -path ./internal/user/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(USER_DB_PORT)/user?sslmode=disable up $(val)
.PHONY: user-db-migrate

user-db-version: ;@ ## Print migration version for user database.
	@migrate -path ./internal/user/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(USER_DB_PORT)/user?sslmode=disable version
.PHONY: user-db-version

user-db-rollback: ;@ ## Rollback user database. Optional <num> argument.
	@migrate -path ./internal/user/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(USER_DB_PORT)/user?sslmode=disable down $(val)
.PHONY: user-db-rollback

user-db-force: ;@ ## Force version on user database. Optional <num> argument.
	@migrate -path ./internal/user/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(USER_DB_PORT)/user?sslmode=disable force $(val)
.PHONY: user-db-force

dynamodb-tables:	;@ ## List Dynamodb tables.
	@aws dynamodb list-tables
.PHONY: dynamodb-tables

routes: ;@ ## Apply ingress routes.
	kubectl apply -f ./manifests/traefik-routes.yaml
.PHONY: routes

ports: ;@ ## Port forward Traefik ports.
	kubectl port-forward --address 0.0.0.0 service/traefik 8000:8000 8080:8080 443:4443 -n default
.PHONY: ports

nats: ;## Port forward NATS port.
	kubectl port-forward statefulset.apps/nats 4222
.PHONY: nats

lint: ;@ ## Run linter.
	@golangci-lint run
.PHONY: lint

help:
	@cat ./setup.txt
	@grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
.PHONY: help

# http://bit.ly/37TR1r2
ifeq ($(firstword $(MAKECMDGOALS)),$(filter $(firstword $(MAKECMDGOALS)),admin-test admin-db-gen admin-db-migrate admin-db-rollback admin-db-force registration-test project-test project-db-gen project-db-migrate project-db-rollback project-db-force user-test user-db-gen user-db-migrate user-db-rollback user-db-force))
  val := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(val):;@:)
endif
