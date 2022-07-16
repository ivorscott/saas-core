//go:generate mockery --all --dir . --case snake --output ../mocks --exported
package repository_test

import (
	"github.com/devpies/saas-core/internal/project/repository"
	"github.com/devpies/saas-core/internal/project/res/testutils"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"os"
	"testing"

	"github.com/devpies/saas-core/internal/project/model"
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

func TestGoldenFiles(t *testing.T) {
	golden := testutils.NewGoldenConfig(false)

	t.Run("column golden files", func(t *testing.T) {
		t.Run("list", func(t *testing.T) {
			var actual []model.Column
			var expected []model.Column

			db, Close := dbConnect.AsRoot()
			defer Close()

			repo := repository.NewColumnRepository(zap.NewNop(), db)

			ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: testutils.MockUUID})

			actual, err := repo.List(ctx, testProjects[0].ID)
			require.NoError(t, err)
			goldenFile := "columns.json"

			if golden.ShouldUpdate() {
				testutils.SaveGoldenFile(&actual, goldenFile)
			}

			testutils.LoadGoldenFile(&expected, goldenFile)
			assert.Equal(t, expected, actual)
		})

		t.Run("retrieve by id", func(t *testing.T) {
			var actual model.Column
			var expected []model.Column

			db, Close := dbConnect.AsRoot()
			defer Close()

			repo := repository.NewColumnRepository(zap.NewNop(), db)

			ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: testutils.MockUUID})

			actual, err := repo.Retrieve(ctx, testColumns[0].ID)
			require.NoError(t, err)
			goldenFile := "columns.json"

			if golden.ShouldUpdate() {
				testutils.SaveGoldenFile(&actual, goldenFile)
			}

			testutils.LoadGoldenFile(&expected, goldenFile)
			assert.Equal(t, expected[0], actual)
		})
	})

	t.Run("project golden files", func(t *testing.T) {
		t.Run("list", func(t *testing.T) {
			var actual []model.Project
			var expected []model.Project

			db, Close := dbConnect.AsRoot()
			defer Close()

			repo := repository.NewProjectRepository(zap.NewNop(), db)

			ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: testutils.MockUUID})

			actual, err := repo.List(ctx)
			require.NoError(t, err)
			goldenFile := "projects.json"

			if golden.ShouldUpdate() {
				testutils.SaveGoldenFile(&actual, goldenFile)
			}

			testutils.LoadGoldenFile(&expected, goldenFile)
			assert.Equal(t, expected, actual)
		})

		t.Run("retrieve by id", func(t *testing.T) {
			var actual model.Project
			var expected []model.Project

			db, Close := dbConnect.AsRoot()
			defer Close()

			repo := repository.NewProjectRepository(zap.NewNop(), db)

			ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: testutils.MockUUID})

			actual, err := repo.Retrieve(ctx, testProjects[0].ID)
			require.NoError(t, err)
			goldenFile := "projects.json"

			if golden.ShouldUpdate() {
				testutils.SaveGoldenFile(&actual, goldenFile)
			}

			testutils.LoadGoldenFile(&expected, goldenFile)
			assert.Equal(t, expected[0], actual)
		})
	})
}
