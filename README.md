# SaaS-Core

Multi-tenant SaaS app built on AWS

## Setup

Read [instructions](docs/SETUP.md).

### Cheatsheets

- [Using NATS](/docs/nats.md)
- [Using Postgres](/docs/postgres.md)

## Using Make
By default, using Tilt allows you to develop against running containers. Alternatively, you can simultaneously run
go binaries natively for an idiomatic go development experience.

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
  --cognito-user-pool-id/$ADMIN_COGNITO_USER_POOL_ID                  <string>    (required)
  --cognito-user-pool-client-id/$ADMIN_COGNITO_USER_POOL_CLIENT_ID    <string>    (required)
  --cognito-region/$ADMIN_COGNITO_REGION                              <string>    (required)
  --db-user/$ADMIN_DB_USER                                            <string>    (noprint,default: postgres)
  --db-password/$ADMIN_DB_PASSWORD                                    <string>    (noprint,default: postgres)
  --db-host/$ADMIN_DB_HOST                                            <string>    (noprint,default: localhost)
  --db-port/$ADMIN_DB_PORT                                            <int>       (noprint,default: 5432)
  --db-name/$ADMIN_DB_NAME                                            <string>    (noprint,default: admin)
  --db-disable-tls/$ADMIN_DB_DISABLE_TLS                              <bool>      (default: false)
  --registration-service-address/$ADMIN_REGISTRATION_SERVICE_ADDRESS  <string>    (default: localhost)
  --registration-service-port/$ADMIN_REGISTRATION_SERVICE_PORT        <string>    (default: 4001)
  --help/-h                                                           display this help message
```

> __TIP__  
>
> Enable `bash-completion` for makefile targets. Add this in your `~/.bash_profile` or `~/.bashrc` file.
> ```bash
> complete -W "\`grep -oE '^[a-zA-Z0-9_.-]+:([^=]|$)' ?akefile | sed 's/[^a-zA-Z0-9_.-]*$//'\`" make
> ```
