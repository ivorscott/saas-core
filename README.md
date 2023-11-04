# SaaS-Core

Multi-tenant SaaS app built on AWS

## Setup

Read [instructions](docs/SETUP.md).

### Cheatsheets

- [Using NATS](docs/nats.md)
- [Using Postgres](docs/postgres.md)
- [Project Structure](docs/structure.md)

## Using Make
By default, using Tilt allows you to develop against running containers. Alternatively, you can simultaneously run
go binaries natively for an idiomatic go development experience.

```bash
> make

- Setup Instructions - 

1. tilt up
2. make ports
3. make routes

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
> __TIP__
>
> Enable `bash-completion` for makefile targets. Open your `~/.zshrc` or `~/.bashrc` file and add:
> ```bash
> complete -W "\`grep -oE '^[a-zA-Z0-9_.-]+:([^=]|$)' ?akefile | sed 's/[^a-zA-Z0-9_.-]*$//'\`" make
> ```

