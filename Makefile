include .env

.DEFAULT_GOAL := help

# =============================================================
# ADMIN SERVICE
# =============================================================
admin-test:	;@ ## Run admin tests. Add -- -v for verbosity.
	go test ./internal/admin/... -v $(val)
.PHONY: admin-test

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
registration-test:	;@ ## Run registration tests. Add -- -v for verbosity.
	go test ./internal/registration/... -v $(val)
.PHONY: registration-test

# =============================================================
# TENANT SERVICE
# =============================================================
tenant-test:	;@ ## Run tenant tests. Add -- -v for verbosity.
	go test ./internal/tenant/... -v $(val)
.PHONY: tenant-test

# =============================================================
# PROJECT SERVICE
# =============================================================
project-test: 	;@ ## Run project tests. Add -- -v for verbosity.
	@go test `go list ./internal/project/... | grep -v repository` $(val)
	@go test ./internal/project/repository/... $(val) -port $(PROJECT_DB_TEST_PORT)
.PHONY: project-test

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

# =============================================================
# SUBSCRIPTION SERVICE
# =============================================================
subscription-test:	;@ ## Run subscription tests. Add -- -v for verbosity.
	go test ./internal/subscription/... -v $(val)
.PHONY: subscription-test

subscription-db: ;@ ## Enter subscription database.
	@pgcli postgres://postgres:postgres@localhost:$(BILLING_DB_PORT)/subscription
.PHONY: subscription-db

subscription-db-gen: ;@ ## Generate migration files. Required <name> argument.
	@migrate create -ext sql -dir ./internal/subscription/res/migrations -seq $(val)
.PHONY: subscription-db-gen

subscription-db-migrate: ;@ ## Migrate subscription database. Optional <num> argument.
	@migrate -path ./internal/subscription/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(BILLING_DB_PORT)/subscription?sslmode=disable up $(val)
.PHONY: subscription-db-migrate

subscription-db-version: ;@ ## Print migration version for subscription database.
	@migrate -path ./internal/subscription/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(BILLING_DB_PORT)/subscription?sslmode=disable version
.PHONY: subscription-db-version

subscription-db-rollback: ;@ ## Rollback subscription database. Optional <num> argument.
	@migrate -path ./internal/subscription/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(BILLING_DB_PORT)/subscription?sslmode=disable down $(val)
.PHONY: subscription-db-rollback

subscription-db-force: ;@ ## Force version on subscription database. Optional <num> argument.
	@migrate -path ./internal/subscription/res/migrations -verbose -database postgres://postgres:postgres@localhost:$(BILLING_DB_PORT)/subscription?sslmode=disable force $(val)
.PHONY: subscription-db-force

init: ;@ ## Initialize project.
	@./scripts/setup-data-folder.sh
.PHONY: init

init-db: admin-db-migrate user-db-migrate project-db-migrate subscription-db-migrate ;@ ## Initialize databases with base schema.
	@echo "Databases initialized!"
.PHONY: init-db

ports: ;@ ## Port forward Traefik ports.
	kubectl port-forward --address 0.0.0.0 service/traefik 8000:8000 8080:8080 443:4443 -n default
.PHONY: ports

routes: ;@ ## Apply ingress routes.
	kubectl apply -f ./manifests/traefik-routes.yaml
.PHONY: routes

nats: ;## Port forward NATS port.
	kubectl port-forward statefulset.apps/nats 4222
.PHONY: nats

lint: ;@ ## Run linter. Optional <package path> argument.
	@golangci-lint run $(val)
.PHONY: lint

test: generate;@ ## Run all tests. Add -- -v for verbosity.
	go test  ./... $(val)
.PHONY: test

cover: ;@ ## Run coverage report.
	go test -cover ./...
.PHONY: cover

generate: ;@ ## Run Go generate.
	go generate ./...
.PHONY: generate

help:
	@grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'
	@echo
	@echo "- Setup Instructions -"
	@echo
	@echo "1. tilt up"
	@echo "2. make ports"
	@echo "3. make routes"
	@echo
.PHONY: help

# http://bit.ly/37TR1r2
# TODO: Find a better way
ifeq ($(firstword $(MAKECMDGOALS)),$(filter $(firstword $(MAKECMDGOALS)),test lint admin-test admin-db-gen admin-db-migrate admin-db-rollback admin-db-force project-test project-db-gen project-db-migrate project-db-rollback project-db-force user-test user-db-gen user-db-migrate user-db-rollback user-db-force subscription-test subscription-db-gen subscription-db-migrate subscription-db-rollback subscription-db-force registration-test tenant-test))
  val := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(val):;@:)
endif
