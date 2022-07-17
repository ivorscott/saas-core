package repository_test

import (
	"context"
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

func TestNewTaskRepository_Create(t *testing.T) {
	expectedTenantID := testProjects[0].TenantID
	expectedProject := testProjects[0]

	tests := []struct {
		name         string
		ctx          context.Context
		projectID    string
		expectations func(t *testing.T, ctx context.Context, repo *repository.TaskRepository, actual model.Task, err error)
	}{
		{
			name:      "success",
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID, UserID: expectedProject.UserID}),
			projectID: expectedProject.ID,
			expectations: func(t *testing.T, ctx context.Context, repo *repository.TaskRepository, actual model.Task, err error) {
				assert.Nil(t, err)
				task, err := repo.Retrieve(ctx, actual.ID)
				assert.Nil(t, err)
				assert.Equal(t, task, actual)
			},
		},
		{
			name:      "context error",
			ctx:       testutils.MockCtx,
			projectID: expectedProject.ID,
			expectations: func(t *testing.T, ctx context.Context, repo *repository.TaskRepository, actual model.Task, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err, web.CtxErr())
			},
		},
		{
			name:      "no tenant error",
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: "", UserID: expectedProject.UserID}),
			projectID: expectedProject.ID,
			expectations: func(t *testing.T, ctx context.Context, repo *repository.TaskRepository, actual model.Task, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrNoTenant, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, Close := dbConnect.AsNonRoot()
			defer Close()

			repo := repository.NewTaskRepository(zap.NewNop(), db)
			nt := model.NewTask{
				Title: "Testing",
			}
			newTask, err := repo.Create(tc.ctx, nt, tc.projectID, time.Now())
			tc.expectations(t, tc.ctx, repo, newTask, err)
		})
	}
}
