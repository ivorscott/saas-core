# SaaS-Core

This project has 3 repositories:

- [sass-core](https://github.com/devpies/saas-core) (this repo)
- [sass-client](https://github.com/devpies/saas-client)
- [sass-infra](https://github.com/devpies/saas-infra)

## Overview

SaaS-Core is a multi-tenant SaaS backend and monorepo of Go services. It supports the SaaS-Client application: a 
simple project management tool like Jira or Trello. 

Necessary services include: registration, tenant and user management, 
project management, billing, and SaaS administration.

## Prerequisites

Read initial setup [instructions](docs/SETUP.md).

### Cheatsheets

- [Using NATS](docs/nats.md)
- [Using Postgres](docs/postgres.md)
- [Project Structure](docs/structure.md)

## Getting Started

```bash
# Run make to print instructions
> make

- Setup Instructions - 

1. tilt up
2. make ports
3. make routes

admin-test           Run admin tests. Add " -- -v" for verbosity.
admin-mock           Generate admin mocks.
admin-db             Enter admin database.
admin-db-gen         Generate migration files. Required <name> argument.
admin-db-migrate     Migrate admin database. Optional <num> argument.
admin-db-version     Print migration version for admin database.
admin-db-rollback    Rollback admin database. Optional <num> argument.
admin-db-force       Force version on admin database. Optional <num> argument.
...
```
> __TIP__
>
> Enable `bash-completion` for makefile targets. Open your `~/.zshrc` or `~/.bashrc` file and add:
> ```bash
> complete -W "\`grep -oE '^[a-zA-Z0-9_.-]+:([^=]|$)' ?akefile | sed 's/[^a-zA-Z0-9_.-]*$//'\`" make
> ```

