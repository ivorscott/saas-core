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

func TestNewColumnRepository_Create(t *testing.T) {
	expectedTenantID := testProjects[0].TenantID
	expectedProject := testProjects[0]

	t.Run("success", func(t *testing.T) {
		db, Close := dbConnect.AsNonRoot()
		defer Close()

		repo := repository.NewColumnRepository(zap.NewNop(), db)
		ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID, UserID: expectedProject.UserID})

		nc := model.NewColumn{
			Title:      "Testing",
			ColumnName: "column-a",
			ProjectID:  expectedProject.ID,
		}
		newColumn, err := repo.Create(ctx, nc, time.Now())
		assert.Nil(t, err)

		column, err := repo.Retrieve(ctx, newColumn.ID)
		assert.Nil(t, err)

		assert.Equal(t, column, newColumn)
	})

	t.Run("context error", func(t *testing.T) {
		db, Close := dbConnect.AsNonRoot()
		defer Close()

		nc := model.NewColumn{}

		repo := repository.NewColumnRepository(zap.NewNop(), db)

		_, err := repo.Create(testutils.MockCtx, nc, time.Now())
		assert.NotNil(t, err)

		assert.Equal(t, err, web.CtxErr())
	})

	t.Run("no tenant error", func(t *testing.T) {
		db, Close := dbConnect.AsNonRoot()
		defer Close()

		nc := model.NewColumn{}

		repo := repository.NewColumnRepository(zap.NewNop(), db)
		ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: ""})

		_, err := repo.Create(ctx, nc, time.Now())
		assert.NotNil(t, err)

		assert.Equal(t, fail.ErrNoTenant, err)
	})
}
