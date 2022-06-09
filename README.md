# SaaS-Core

This project is a part of "AWS SaaS app in 30 days" - _Proof of Concept_

### Requirements
- mac or linux machine
- aws account
- install [go v1.18 or higher](https://go.dev/doc/install)
- install [tilt](https://tilt.dev/)
- install docker-desktop - enable kubernetes
- install [mockery](https://github.com/vektra/mockery)
- install cypress `npm install cypress -g`
- install [pgcli](https://www.pgcli.com/)
- install [golangci-lint](https://github.com/golangci/golangci-lint)
- install [go-migrate](https://github.com/golang-migrate/migrate)
- [saas-infra resources](https://github.com/devpies/saas-infra/tree/main/local/saas) 

## Getting Started
Print instructions.

```bash
> make
admin             Run admin app with live reload.
admin-end         Run end-to-end admin tests with Cypress.
admin-test        Run admin tests. Add " -- -v" for verbosity.
admin-mock        Generate admin mocks.
admin-db          Enter admin database.
admin-db-gen      Generate migration files. Required <name> argument.
admin-db-migrate  Migrate admin database. Optional <num> argument.
admin-db-version  Print migration version for admin database.
admin-db-rollback Rollback admin database. Optional <num> argument.
lint              Run linter.
...
```

## Environment Variables

The `.env` file contains variables for all programs. Using `make` automatically references these values.
Program requirements are also documented in help text. 
```bash
> go run ./cmd/admin -h
Usage: admin [options] [arguments]

OPTIONS
  --web-production/$ADMIN_WEB_PRODUCTION                              <bool>      (default: false)
  --web-read-timeout/$ADMIN_WEB_READ_TIMEOUT                          <duration>  (default: 5s)
  --web-write-timeout/$ADMIN_WEB_WRITE_TIMEOUT                        <duration>  (default: 5s)
  --web-shutdown-timeout/$ADMIN_WEB_SHUTDOWN_TIMEOUT                  <duration>  (default: 5s)
  --web-address/$ADMIN_WEB_ADDRESS                                    <string>    (default: localhost)
  --web-port/$ADMIN_WEB_PORT                                          <string>    (default: 4000)
  --cognito-app-client-id/$ADMIN_COGNITO_APP_CLIENT_ID                <string>    (required)
  --cognito-user-pool-client-id/$ADMIN_COGNITO_USER_POOL_CLIENT_ID    <string>    (required)
  --cognito-region/$ADMIN_COGNITO_REGION                              <string>    (default: eu-central-1)
  --postgres-user/$ADMIN_POSTGRES_USER                                <string>    (required)
  --postgres-password/$ADMIN_POSTGRES_PASSWORD                        <string>    (required)
  --postgres-host/$ADMIN_POSTGRES_HOST                                <string>    (required)
  --postgres-port/$ADMIN_POSTGRES_PORT                                <int>       (required)
  --postgres-db/$ADMIN_POSTGRES_DB                                    <string>    (required)
  --postgres-disable-tls/$ADMIN_POSTGRES_DISABLE_TLS                  <bool>      (default: false)
  --registration-service-address/$ADMIN_REGISTRATION_SERVICE_ADDRESS  <string>    (default: localhost)
  --registration-service-port/$ADMIN_REGISTRATION_SERVICE_PORT        <string>    (default: 4001)
  --help/-h                                                           
  display this help message
```

> __TIP__  
> 
> 1. Using `make` is the easiest way to get started. However, if you choose to run go binaries directly, you can export the `.env` file variables to avoid using CLI flags:  
> ```bash
> export $(grep -v '^#' .env | xargs)
>```
> 2. Enable bash-completion of the makefile targets. Add this in your `~/.bash_profile` file or `~/.bashrc` file.
> ```bash
> complete -W "\`grep -oE '^[a-zA-Z0-9_.-]+:([^=]|$)' ?akefile | sed 's/[^a-zA-Z0-9_.-]*$//'\`" make
> ```

## Debugging NATS 

```bash
kubectl port-forward statefulset.apps/nats 4222
nats account info
nats stream info my_stream
nats consumer info
nats stream view my_stream # view messages
nats stream rmm # remove message
nats stream purge # remove all messages
```
https://docs.nats.io/nats-concepts/jetstream/js_walkthrough

