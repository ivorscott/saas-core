package repository_test

import (
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
	var testTenant = testProjects[0].TenantID

	t.Run("success", func(t *testing.T) {
		db, Close := testutils.DBConnect().AsNonRoot()
		defer Close()

		p := model.NewProject{
			Name: "My Project",
		}

		repo := repository.NewProjectRepository(zap.NewNop(), db)

		ctx := web.NewContext(testCtx, &web.Values{TenantID: testTenant})

		actual, err := repo.Create(ctx, p, time.Now())
		assert.Nil(t, err)

		assert.Equal(t, testTenant, actual.TenantID)
		assert.Equal(t, p.Name, actual.Name)
	})
}
