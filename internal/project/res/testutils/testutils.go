// Package testutils contains integration test helpers.
package testutils

import (
	"context"
	"github.com/devpies/saas-core/internal/project/config"
	"github.com/devpies/saas-core/internal/project/db"
	"github.com/devpies/saas-core/internal/project/res"
	"go.uber.org/zap"
	"os"
	"path"
	"runtime"
	"testing"
)

const (
	MockRequirement = "mock-requirement"
	MockUUID        = "ac7b523d-1eb9-43f3-bd33-c3e8106c2e70"

	dbDriver = "postgres"
)

var MockCtx = context.Background()

// AsRole enables database role switching in tests.
type AsRole struct {
	cfg config.Config
}

// DBConnect prepares the test database connection.
func DBConnect() *AsRole {
	prepareEnvironment()
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	return &AsRole{
		cfg: cfg,
	}
}

// AsRoot connects to test database as root.
func (a *AsRole) AsRoot() (*db.PostgresDatabase, func() error) {
	return a.setupDBAsRoot()
}

// AsNonRoot connects to test database as non-root.
func (a *AsRole) AsNonRoot() (*db.PostgresDatabase, func() error) {
	_, Close := a.setupDBAsRoot()
	_ = Close()
	a.cfg.DB.User = "user_a"
	repo, Close, err := db.NewPostgresDatabase(zap.NewNop(), a.cfg)
	if err != nil {
		panic(err)
	}
	return repo, Close
}

// setupDBAsRoot migrates and loads fixtures.
func (a *AsRole) setupDBAsRoot() (*db.PostgresDatabase, func() error) {
	a.cfg.DB.User = "postgres"
	repo, Close, err := db.NewPostgresDatabase(zap.NewNop(), a.cfg)
	if err != nil {
		panic(err)
	}
	err = res.MigrateUp(repo.URL.String())
	if err != nil {
		panic(err)
	}
	err = loadFixtures(repo.TestsOnlyDBConnection())
	if err != nil {
		panic(err)
	}
	return repo, Close
}

// resFile returns a file from the resource directory.
func resFile(pathElem ...string) string {
	_, thisFile, _, _ := runtime.Caller(0)
	parentDir, _ := path.Split(path.Dir(thisFile))
	return path.Join(parentDir, path.Join(pathElem...))
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
			value: "30019",
		},
		{
			key:   "PROJECT_DB_DISABLE_TLS",
			value: "true",
		},
		{
			key:   "PROJECT_COGNITO_USER_POOL_ID",
			value: MockRequirement,
		},
		{
			key:   "PROJECT_COGNITO_REGION",
			value: MockRequirement,
		},
	}
	for _, e := range env {
		err := os.Setenv(e.key, e.value)
		if err != nil {
			panic(err)
		}
	}
}

func Debug[T any](t *testing.T, data T) {
	wrapper := "\n\nDEBUG ================================\n\n"
	t.Logf("%s %+v %s", wrapper, data, wrapper)
}
