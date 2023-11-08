# SaaS-Core

This project has 3 repositories:

- [saas-core](https://github.com/devpies/saas-core) (this repo)
- [saas-client](https://github.com/devpies/saas-client)
- [saas-infra](https://github.com/devpies/saas-infra)

## Overview

SaaS-Core is a multi-tenant SaaS backend and monorepo of Go services. It supports the SaaS-Client application: a 
simple project management tool like Jira or Trello. 

Necessary services include: registration, tenant and user management, 
project management, subscriptions, and SaaS administration.

## Prerequisites

Read initial setup [instructions](doc/guide/SETUP.md).

### Guides

- [Using NATS](doc/guide/nats.md)
- [Using Postgres](doc/guide/postgres.md)
- [Project Structure](doc/guide/structure.md)

## Getting Started

```bash
# Run make to print instructions
> make
...
subscription-test         Run subscription tests. Add -- -v for verbosity.
subscription-mock         Generate subscription mocks.
subscription-db           Enter subscription database.
subscription-db-gen       Generate migration files. Required <name> argument.
subscription-db-migrate   Migrate subscription database. Optional <num> argument.
subscription-db-version   Print migration version for subscription database.
subscription-db-rollback  Rollback subscription database. Optional <num> argument.
subscription-db-force     Force version on subscription database. Optional <num> argument.
init                      Initialize project. Do once.
ports                     Port forward Traefik ports.
routes                    Apply ingress routes.
nats                      Port forward NATS port.
lint                      Run linter. Optional <package path> argument.
test                      Run all tests. Add -- -v for verbosity.

- Setup Instructions -

1. tilt up
2. make ports
3. make routes

```
> __TIP__
>
> Enable `bash-completion` for makefile targets. Open your `~/.zshrc` or `~/.bashrc` file and add:
> ```bash
> complete -W "\`grep -oE '^[a-zA-Z0-9_.-]+:([^=]|$)' ?akefile | sed 's/[^a-zA-Z0-9_.-]*$//'\`" make
> ```

