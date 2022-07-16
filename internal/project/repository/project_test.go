package repository_test

import (
	"github.com/aws/aws-sdk-go-v2/aws"
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

		actual, err := repo.Retrieve(ctx, "mock")
		assert.NotNil(t, err)

		assert.Equal(t, fail.ErrInvalidID, err)
		assert.NotEqual(t, expectedProject, actual)
	})

	t.Run("data isolation between tenants", func(t *testing.T) {
		db, Close := testutils.DBConnect().AsNonRoot()
		defer Close()

		repo := repository.NewProjectRepository(zap.NewNop(), db)
		ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: testutils.MockUUID})

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

func TestProjectRepository_Update(t *testing.T) {
	expectedTenantID := testProjects[0].TenantID

	t.Run("success", func(t *testing.T) {
		t.Run("update project name", func(t *testing.T) {
			expectedProject := testProjects[0]

			db, Close := testutils.DBConnect().AsNonRoot()
			defer Close()

			repo := repository.NewProjectRepository(zap.NewNop(), db)
			ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID})

			update := model.UpdateProject{
				Name: aws.String("Updated"),
			}

			actual, err := repo.Update(ctx, expectedProject.ID, update, time.Now())
			assert.Nil(t, err)

			assert.NotEqual(t, expectedProject, actual)

			expectedProject.Name = *update.Name

			assert.Equal(t, expectedProject, actual)
		})

		t.Run("update project description", func(t *testing.T) {
			expectedProject := testProjects[0]

			db, Close := testutils.DBConnect().AsNonRoot()
			defer Close()

			repo := repository.NewProjectRepository(zap.NewNop(), db)
			ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID})

			update := model.UpdateProject{
				Description: aws.String("Updated"),
			}

			actual, err := repo.Update(ctx, expectedProject.ID, update, time.Now())
			assert.Nil(t, err)

			assert.NotEqual(t, expectedProject, actual)

			expectedProject.Description = *update.Description

			assert.Equal(t, expectedProject, actual)
		})

		t.Run("update project active field", func(t *testing.T) {
			expectedProject := testProjects[0]

			db, Close := testutils.DBConnect().AsNonRoot()
			defer Close()

			repo := repository.NewProjectRepository(zap.NewNop(), db)
			ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID})

			update := model.UpdateProject{
				Active: aws.Bool(false),
			}

			actual, err := repo.Update(ctx, expectedProject.ID, update, time.Now())
			assert.Nil(t, err)

			assert.NotEqual(t, expectedProject, actual)

			expectedProject.Active = *update.Active

			assert.Equal(t, expectedProject, actual)
		})

		t.Run("update project public field", func(t *testing.T) {
			expectedProject := testProjects[0]

			db, Close := testutils.DBConnect().AsNonRoot()
			defer Close()

			repo := repository.NewProjectRepository(zap.NewNop(), db)
			ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID})

			update := model.UpdateProject{
				Public: aws.Bool(true),
			}

			actual, err := repo.Update(ctx, expectedProject.ID, update, time.Now())
			assert.Nil(t, err)

			assert.NotEqual(t, expectedProject, actual)

			expectedProject.Public = *update.Public

			assert.Equal(t, expectedProject, actual)
		})

		t.Run("update project column order", func(t *testing.T) {
			expectedProject := testProjects[0]

			db, Close := testutils.DBConnect().AsNonRoot()
			defer Close()

			repo := repository.NewProjectRepository(zap.NewNop(), db)
			ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID})

			update := model.UpdateProject{
				ColumnOrder: []string{"column-4", "column-3", "column-2", "column-1"},
			}

			actual, err := repo.Update(ctx, expectedProject.ID, update, time.Now())
			assert.Nil(t, err)

			assert.NotEqual(t, expectedProject, actual)

			expectedProject.ColumnOrder = update.ColumnOrder

			assert.Equal(t, expectedProject, actual)
		})
	})

	t.Run("id not UUID", func(t *testing.T) {
		db, Close := testutils.DBConnect().AsNonRoot()
		defer Close()

		repo := repository.NewProjectRepository(zap.NewNop(), db)
		ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID})

		update := model.UpdateProject{}

		_, err := repo.Update(ctx, "mock", update, time.Now())
		testutils.Debug(t, err)
		assert.NotNil(t, err)

		assert.Equal(t, fail.ErrInvalidID, err)
	})

	t.Run("context error", func(t *testing.T) {
		expectedProject := testProjects[0]

		db, Close := testutils.DBConnect().AsNonRoot()
		defer Close()

		repo := repository.NewProjectRepository(zap.NewNop(), db)

		update := model.UpdateProject{}

		_, err := repo.Update(testutils.MockCtx, expectedProject.ID, update, time.Now())
		assert.NotNil(t, err)

		assert.Equal(t, err, web.CtxErr())
	})

	t.Run("unauthorized error", func(t *testing.T) {
		expectedProject := testProjects[0]

		db, Close := testutils.DBConnect().AsNonRoot()
		defer Close()

		repo := repository.NewProjectRepository(zap.NewNop(), db)
		ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: ""})

		update := model.UpdateProject{}

		_, err := repo.Update(ctx, expectedProject.ID, update, time.Now())
		assert.NotNil(t, err)

		assert.Equal(t, err.Error(), fail.ErrNotAuthorized.Error())
	})
}
