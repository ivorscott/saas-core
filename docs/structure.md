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
   - resources for the project. [Learn more](#res)


Shared libraries, are kept in `pkg` to enforce consistency across services.

1. __web__ 
   - a custom web framework. 
2. __log__
   - logging configuration.
3. __events__
   - events used for service communication.

## Res
Resources for projects.

- [Fixtures](#fixtures)
- [Seeds](#seeds)
- [Migrations](#migrations)
- [Golden Files](#golden-files)
- [Test Utils](#test-utils)

The resource folder includes, test fixtures, golden files, migrations and database seeds.

```bash
fixtures # test fixtures
golden # golden files
migrations # migrations
seed # data for development
testutils # test helpers
```

### Fixtures

`res/fixtures`

Test fixtures are only loaded into test databases. Fixtures data is feed into our repository level tests. Fixtures allow
the Go service to be tested against a real database instead of running it against mocks, which may lead to production bugs
not being caught in the tests.

__Before every test, the test database is cleaned and the fixture data is loaded into
the database.__ https://github.com/go-testfixtures/testfixtures

### Seeds

`res/seed`

Seed data should be updated as databases change. Keep it simple and maintain a single seed file for each database used by the service.

### Migrations

`res/migrations`

Migrations are managed via `res.MigrateUp()` and via
`make` commands. [Learn more](/README.md#migration-and-seeding).

### Golden Files

`res/golden`

Goldenfiles are used in tests to compare database responses with previous queries preserved as snapshots in json format.
If a database response changes, the golden file test fails and a new snapshot must be saved for the test to pass.
To update all golden files change the following function argument to `true`:
```
golden := testutils.NewGoldenConfig(true)
```
Alternatively, if you want one golden file to update, comment the corresponding
code block:

```go
// pkg/repository/repository_test.go

goldenFile := "employee.json"

//if golden.ShouldUpdate() {
    testutils.SaveGoldenFile(&actual, goldenFile)
//}
```
Then re-run the  tests.

### Test Utils

Test utils are test helpers.