package repository_test

import (
	"github.com/devpies/saas-core/internal/project/fail"
	"github.com/devpies/saas-core/internal/project/model"
	"github.com/devpies/saas-core/internal/project/repository"
	"github.com/devpies/saas-core/internal/project/res/testutils"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestProjectRepository_Create(t *testing.T) {
	expectedTenantID := testProjects[0].TenantID

	t.Run("success", func(t *testing.T) {
		db, Close := testutils.DBConnect().AsNonRoot()
		defer Close()

		p := model.NewProject{
			Name: "My Project",
		}

		repo := repository.NewProjectRepository(zap.NewNop(), db)
		ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID})

		actual, err := repo.Create(ctx, p, time.Now())
		assert.Nil(t, err)

		assert.Equal(t, expectedTenantID, actual.TenantID)
		assert.Equal(t, "MYP-", actual.Prefix)
		assert.Equal(t, []string{"column-1", "column-2", "column-3", "column-4"}, actual.ColumnOrder)
		assert.Equal(t, p.Name, actual.Name)
	})

	t.Run("context error", func(t *testing.T) {
		db, Close := testutils.DBConnect().AsNonRoot()
		defer Close()

		p := model.NewProject{
			Name: "My Project",
		}

		repo := repository.NewProjectRepository(zap.NewNop(), db)

		_, err := repo.Create(testutils.MockCtx, p, time.Now())
		assert.NotNil(t, err)

		assert.Equal(t, err, web.CtxErr())
	})

	t.Run("unauthorized error", func(t *testing.T) {
		db, Close := testutils.DBConnect().AsNonRoot()
		defer Close()

		p := model.NewProject{
			Name: "My Project",
		}

		repo := repository.NewProjectRepository(zap.NewNop(), db)
		ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: ""})

		_, err := repo.Create(ctx, p, time.Now())
		assert.NotNil(t, err)

		assert.Equal(t, err.Error(), fail.ErrNotAuthorized.Error())
	})
}

func TestProjectRepository_Retrieve(t *testing.T) {
	expectedTenantID := testProjects[0].TenantID
	expectedProject := testProjects[0]

	t.Run("success", func(t *testing.T) {
		db, Close := testutils.DBConnect().AsNonRoot()
		defer Close()

		repo := repository.NewProjectRepository(zap.NewNop(), db)
		ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID})

		actual, err := repo.Retrieve(ctx, expectedProject.ID)
		assert.Nil(t, err)

		assert.Equal(t, expectedProject, actual)
	})

	t.Run("id not UUID", func(t *testing.T) {
		db, Close := testutils.DBConnect().AsNonRoot()
		defer Close()

		repo := repository.NewProjectRepository(zap.NewNop(), db)
		ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID})

		actual, err := repo.Retrieve(ctx, testutils.MockRequirement)
		assert.NotNil(t, err)

		assert.Equal(t, fail.ErrInvalidID, err)
		assert.NotEqual(t, expectedProject, actual)
	})

	t.Run("data isolation between tenants", func(t *testing.T) {
		db, Close := testutils.DBConnect().AsNonRoot()
		defer Close()

		repo := repository.NewProjectRepository(zap.NewNop(), db)
		ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: testutils.MockRequirement})

		actual, err := repo.Retrieve(ctx, expectedProject.ID)
		assert.NotNil(t, err)

		assert.Equal(t, fail.ErrNotFound, err)
		assert.NotEqual(t, expectedProject, actual)
	})

	t.Run("not found", func(t *testing.T) {
		db, Close := testutils.DBConnect().AsNonRoot()
		defer Close()

		repo := repository.NewProjectRepository(zap.NewNop(), db)
		ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID})

		actual, err := repo.Retrieve(ctx, testutils.MockUUID)
		assert.NotNil(t, err)

		assert.Equal(t, fail.ErrNotFound, err)
		assert.NotEqual(t, expectedProject, actual)
	})

	t.Run("context error", func(t *testing.T) {
		db, Close := testutils.DBConnect().AsNonRoot()
		defer Close()

		repo := repository.NewProjectRepository(zap.NewNop(), db)

		_, err := repo.Retrieve(testutils.MockCtx, expectedProject.ID)
		assert.NotNil(t, err)

		assert.Equal(t, err, web.CtxErr())
	})

	t.Run("unauthorized error", func(t *testing.T) {
		db, Close := testutils.DBConnect().AsNonRoot()
		defer Close()

		repo := repository.NewProjectRepository(zap.NewNop(), db)
		ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: ""})

		_, err := repo.Retrieve(ctx, expectedProject.ID)
		assert.NotNil(t, err)

		assert.Equal(t, err.Error(), fail.ErrNotAuthorized.Error())
	})
}
