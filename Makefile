include .env

.DEFAULT_GOAL := help

# =============================================================
# ADMIN SERVICE
# =============================================================
admin-test: admin-mock	;@ ## Run admin tests. Add -- -v for verbosity.
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

# =============================================================
# REGISTRATION SERVICE
# =============================================================
registration-test: registration-mock	;@ ## Run registration tests. Add -- -v for verbosity.
	go test $(val) -cover ./internal/registration/...
.PHONY: registration-test

registration-mock: ;@ ## Generate registration mocks.
	go generate ./internal/registration/...
.PHONY: registration-mock

# =============================================================
# TENANT SERVICE
# =============================================================
tenant-test: tenant-mock	;@ ## Run tenant tests. Add -- -v for verbosity.
	go test $(val) -cover ./internal/tenant/...
.PHONY: tenant-test

tenant-mock: ;@ ## Generate tenant mocks.
	go generate ./internal/tenant/...
.PHONY: tenant-mock

# =============================================================
# PROJECT SERVICE
# =============================================================
project-test: project-mock	;@ ## Run project tests. Add -- -v for verbosity.
	go test $(val) -cover ./internal/project/...
.PHONY: project-test

project-mock: ;@ ## Generate project mocks.
	go generate ./internal/project/...
.PHONY: project-mock

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

# =============================================================
# USER SERVICE
# =============================================================
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

# Billing service =============================================================
billing-test: billing-test	;@ ## Run billing tests. Add -- -v for verbosity.
	go test $(val) -cover ./internal/billing/...
.PHONY: billing-test

billing-mock: ;@ ## Generate billing mocks.
	go generate ./internal/billing/...
.PHONY: billing-mock


billing-db: ;@ ## Enter billing database.
	@pgcli postgres://postgres:postgres@localhost:$(BILLING_DB_PORT)/billing
.PHONY: billing-db

billing-db-gen: ;@ ## Generate migration files. Required <name> argument.
	@migrate create -ext sql -dir ./internal/billing/res/migrations -seq $(val)
.PHONY: billing-db-gen

billing-db-migrate: ;@ ## Migrate billing database. Optional <num> argument.
	@migrate -path ./internal/billing/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(BILLING_DB_PORT)/billing?sslmode=disable up $(val)
.PHONY: billing-db-migrate

billing-db-version: ;@ ## Print migration version for billing database.
	@migrate -path ./internal/billing/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(BILLING_DB_PORT)/billing?sslmode=disable version
.PHONY: billing-db-version

billing-db-rollback: ;@ ## Rollback billing database. Optional <num> argument.
	@migrate -path ./internal/billing/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(BILLING_DB_PORT)/billing?sslmode=disable down $(val)
.PHONY: billing-db-rollback

billing-db-force: ;@ ## Force version on billing database. Optional <num> argument.
	@migrate -path ./internal/billing/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(BILLING_DB_PORT)/billing?sslmode=disable force $(val)
.PHONY: billing-db-force

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

lint: ;@ ## Run linter. Optional <package path> argument.
	@golangci-lint run $(val)
.PHONY: lint

test: ;@ ## Run all tests. Add -- -v for verbosity.
	go test $(val) -cover ./...
.PHONY: test

help:
	@echo
	@echo "- Setup Instructions -"
	@echo
	@echo 1. tilt up
	@echo 2. make ports
	@echo 3. make routes
	@echo
	@grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
.PHONY: help

# http://bit.ly/37TR1r2
# TODO: Find a better way
ifeq ($(firstword $(MAKECMDGOALS)),$(filter $(firstword $(MAKECMDGOALS)),test lint admin-test admin-db-gen admin-db-migrate admin-db-rollback admin-db-force project-test project-db-gen project-db-migrate project-db-rollback project-db-force user-test user-db-gen user-db-migrate user-db-rollback user-db-force billing-test billing-db-gen billing-db-migrate billing-db-rollback billing-db-force registration-test tenant-test))
  val := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(val):;@:)
endif
