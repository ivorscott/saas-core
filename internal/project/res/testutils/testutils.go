// Package testutils contains integration test helpers.
package testutils

import (
	"context"
	"database/sql"
	"os"
	"path"
	"runtime"

	"github.com/devpies/saas-core/internal/project/config"
	"github.com/devpies/saas-core/internal/project/db"
	"github.com/devpies/saas-core/internal/project/res"

	"github.com/go-testfixtures/testfixtures/v3"
	"go.uber.org/zap"
)

const (
	// MockUUID represents a mock uuid.
	MockUUID    = "ac7b523d-1eb9-43f3-bd33-c3e8106c2e70"
	rootUser    = "postgres"
	dbDriver    = "postgres"
	fixturesDir = "fixtures"
)

// MockCtx represents an empty context.
var MockCtx = context.Background()

// DatabaseClient sets up a database and enables role switching.
type DatabaseClient struct {
	cfg      config.Config
	fixtures *testfixtures.Loader
	user     string
	root     string
}

// NewDatabaseClient returns a DatabaseClient and a close method to clean up the setup connection.
func NewDatabaseClient() (*DatabaseClient, func() error) {
	prepareEnvironment()

	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	c := DatabaseClient{
		cfg:  cfg,
		root: rootUser,
		user: cfg.DB.User,
	}

	return &c, c.setupDB()
}

// setupDB runs migrations and fixtures as root.
func (client *DatabaseClient) setupDB() func() error {
	client.setRole(client.root)
	pg, Close := client.connect()

	err := res.MigrateUp(pg.URL.String())
	if err != nil {
		panic(err)
	}
	client.setFixtures(pg.TestsOnlyDBConnection())
	return Close
}

// AsRoot connects as root.
func (client *DatabaseClient) AsRoot() (*db.PostgresDatabase, func() error) {
	client.setRole(client.root)
	client.loadFixtures()
	return client.connect()
}

// AsNonRoot connects as non-root.
func (client *DatabaseClient) AsNonRoot() (*db.PostgresDatabase, func() error) {
	client.setRole(client.user)
	client.loadFixtures()
	return client.connect()
}

// setRole sets the user role.
func (client *DatabaseClient) setRole(role string) {
	client.cfg.DB.User = role
}

// connect creates a database connection.
func (client *DatabaseClient) connect() (*db.PostgresDatabase, func() error) {
	repo, Close, err := db.NewPostgresDatabase(zap.NewNop(), client.cfg)
	if err != nil {
		panic(err)
	}
	return repo, Close
}

func (client *DatabaseClient) loadFixtures() {
	err := client.fixtures.Load()
	if err != nil {
		panic(err)
	}
}

func (client *DatabaseClient) setFixtures(db *sql.DB) {
	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect(dbDriver),
		testfixtures.Directory(resFile(fixturesDir)),
	)
	if err != nil {
		panic(err)
	}
	client.fixtures = fixtures
}

// prepareEnvironment prepares environment for integration testing with postgres while mocking cognito.
func prepareEnvironment() {
	env := []struct {
		key   string
		value string
	}{
		{
			key:   "PROJECT_DB_NAME",
			value: "project_test",
		},
		{
			key:   "PROJECT_DB_PORT",
			value: "5432",
		},
		{
			key:   "PROJECT_DB_DISABLE_TLS",
			value: "true",
		},
		{
			key:   "PROJECT_COGNITO_USER_POOL_ID",
			value: "mock",
		},
		{
			key:   "PROJECT_COGNITO_REGION",
			value: "mock",
		},
	}
	for _, e := range env {
		err := os.Setenv(e.key, e.value)
		if err != nil {
			panic(err)
		}
	}
}

// resFile returns a file from the resource directory.
func resFile(pathElem ...string) string {
	_, thisFile, _, _ := runtime.Caller(0)
	parentDir, _ := path.Split(path.Dir(thisFile))
	return path.Join(parentDir, path.Join(pathElem...))
}
