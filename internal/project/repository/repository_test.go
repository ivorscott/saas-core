//go:generate mockery --all --dir . --case snake --output ../mocks --exported
package repository_test

import (
	"context"
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
	testCtx      = context.Background()
	testProjects []model.Project
)

func TestMain(m *testing.M) {
	testutils.LoadGoldenFile(&testProjects, "projects.json")

	os.Exit(m.Run())
}

func TestGoldenFiles(t *testing.T) {
	golden := testutils.NewGoldenConfig(false)

	t.Run("projects golden file", func(t *testing.T) {
		var actual []model.Project
		var expected []model.Project

		db, Close := testutils.DBConnect().AsRoot()
		defer Close()

		repo := repository.NewProjectRepository(zap.NewNop(), db)

		ctx := web.NewContext(testCtx, &web.Values{TenantID: ""})

		actual, err := repo.List(ctx)
		require.NoError(t, err)
		goldenFile := "projects.json"

		if golden.ShouldUpdate() {
			testutils.SaveGoldenFile(&actual, goldenFile)
		}

		testutils.LoadGoldenFile(&expected, goldenFile)
		assert.Equal(t, expected, actual)
	})
}
