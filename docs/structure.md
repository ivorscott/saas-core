# Project Structure

The `internal` folder contains the application source code. Services commonly adopt a [3-layered architecture](https://en.wikipedia.org/wiki/Multitier_architecture) 
represented as 3 packages:

1. __handler__ 
   - presentation layer for handling incoming requests.
2. __service__
   - application layer for handling business logic.
3. __repository__
   - data access layer for handling queries.

In addition to these conventions, services may also contain:

1. __config__
   - environment variable configuration.
2. __db__
   - required database clients.
3. __model__
   - data transfer objects and validation.
4. __res__
   - additional resources for the services. [Learn more](#res)


Shared libraries, are kept in `pkg` to enforce consistency across services.

1. __web__ 
   - a custom web framework. 
2. __log__
   - logging configuration.
3. __events__
   - events used for service communication.

## Res
Additional resources for services.

- [Fixtures](#fixtures)
- [Migrations](#migrations)
- [Golden Files](#golden-files)
- [Test Utils](#test-utils)

The resource folder includes, test fixtures, golden files, and migrations.

```bash
fixtures # test fixtures
golden # golden files
migrations # migrations
testutils # integration test helpers
```

### Fixtures

`res/fixtures`

Test fixtures are only loaded into test databases. Fixtures data is feed into our repository level tests. Fixtures allow
the Go service to be tested against a real database instead of running it against mocks, which may lead to production bugs
not being caught in the tests.

__Before every test, the test database is cleaned and the fixture data is loaded into
the database.__ https://github.com/go-testfixtures/testfixtures

### Migrations

`res/migrations`

Migrations use [go-migrate](https://github.com/golang-migrate/migrate) and are managed via `res.MigrateUp()` and `make` commands.
Always create backups before migrating (even locally). [Learn more.](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

```bash
# Example make targets
project-db-gen       # Generate migration files. Required <name> argument.
project-db-migrate   # Migrate project database. Optional <num> argument.
project-db-version   # Print migration version for project database.
project-db-rollback  # Rollback project database. Optional <num> argument.
project-db-force     # Force version on project database. Optional <num> argument.
```

### Golden Files

`res/golden`

Goldenfiles are used in tests to compare database responses with previous queries preserved as snapshots in json format.
If a database response changes, the golden file test fails and a new snapshot must be saved for the test to pass.

Goldenfiles help detect unwanted changes in the data access layer. Always track down why something has changed before updating
a goldenfile to pass a failing test. To update all golden files change the following function argument to `true`:
```
golden := testutils.NewGoldenConfig(false)
```
Doing this will always update the goldenfiles so make sure to change it back to `false` after running the tests again, otherwise,
the goldenfiles will always update. 

Alternatively, if you want one golden file to update you can use comments:

```go
// repository/repository_test.go

goldenFile := "projects.json"

//if golden.ShouldUpdate() {
    testutils.SaveGoldenFile(&actual, goldenFile)
//}
```
Then re-run the  tests.

### Test Utils

Test utils are integration test helpers. The package contains a database client that sets up
the test database with migrations and fixtures. The package provides the ability to switch between database users. This is useful because postgres 
[row-level security](https://www.postgresql.org/docs/current/ddl-rowsecurity.html) is bypassed for the `root` user. To 
ensure data isolation is working as expected use the `non-root` user in tests unless you have a good reason to do otherwise.

Initialize `testutils.NewDatabaseClient` in TestMain, and load test fixtures.

```go
package repository_test

import (
    ...
)

var (
	dbConnect    *testutils.DatabaseClient
	testProjects []model.Project
	testColumns  []model.Column
)

func TestMain(m *testing.M) {
   db, dbClose := testutils.NewDatabaseClient()
   dbConnect = db
   defer dbClose()

   testutils.LoadGoldenFile(&testProjects, "projects.json")
   testutils.LoadGoldenFile(&testColumns, "columns.json")

   os.Exit(m.Run())
}
```
Then in your tests use the `dbConnect` variable:

__Connecting to the test database as `root`.__

```go
db, Close := dbConnect.AsRoot()
```

__Connecting to the test database as `non-root`.__

```go
db, Close := dbConnect.AsNonRoot()
```