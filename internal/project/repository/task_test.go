package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/devpies/saas-core/internal/project/fail"
	"github.com/devpies/saas-core/internal/project/model"
	"github.com/devpies/saas-core/internal/project/repository"
	"github.com/devpies/saas-core/internal/project/res/testutils"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestTaskRepository_Create(t *testing.T) {
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
			name:      "project id is not UUID",
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID, UserID: expectedProject.UserID}),
			projectID: "mock",
			expectations: func(t *testing.T, ctx context.Context, repo *repository.TaskRepository, actual model.Task, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrInvalidID, err)
			},
		},
		{
			name:      "user id is not UUID",
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID, UserID: "mock"}),
			projectID: expectedProject.ID,
			expectations: func(t *testing.T, ctx context.Context, repo *repository.TaskRepository, actual model.Task, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrInvalidID, err)
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

func TestTaskRepository_Retrieve(t *testing.T) {
	expectedTenantID := testProjects[0].TenantID
	expectedTask := testTasks[0]

	tests := []struct {
		name         string
		ctx          context.Context
		taskID       string
		expectations func(t *testing.T, ctx context.Context, actual model.Task, err error)
	}{
		{
			name:   "success",
			ctx:    web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID, UserID: testutils.MockUUID}),
			taskID: expectedTask.ID,
			expectations: func(t *testing.T, ctx context.Context, actual model.Task, err error) {
				assert.Nil(t, err)
				assert.Equal(t, expectedTask, actual)
			},
		},
		{
			name:   "task id is not UUID",
			ctx:    web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID, UserID: testutils.MockUUID}),
			taskID: "mock",
			expectations: func(t *testing.T, ctx context.Context, actual model.Task, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrInvalidID, err)
				assert.NotEqual(t, expectedTask, actual)
			},
		},
		{
			name:   "context error",
			ctx:    testutils.MockCtx,
			taskID: expectedTask.ID,
			expectations: func(t *testing.T, ctx context.Context, actual model.Task, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err, web.CtxErr())
			},
		},
		{
			name:   "no tenant error",
			ctx:    web.NewContext(testutils.MockCtx, &web.Values{TenantID: "", UserID: testutils.MockUUID}),
			taskID: expectedTask.ID,
			expectations: func(t *testing.T, ctx context.Context, actual model.Task, err error) {
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
			newTask, err := repo.Retrieve(tc.ctx, tc.taskID)
			tc.expectations(t, tc.ctx, newTask, err)
		})
	}
}

func TestTaskRepository_List(t *testing.T) {
	expectedTenantID := testProjects[1].TenantID
	expectedProject := testProjects[1]

	expectedTasks := make([]model.Task, 0)
	for _, v := range testTasks {
		if v.ProjectID == expectedProject.ID {
			expectedTasks = append(expectedTasks, v)
		}
	}

	tests := []struct {
		name         string
		ctx          context.Context
		projectID    string
		expectations func(t *testing.T, ctx context.Context, actual []model.Task, err error)
	}{
		{
			name:      "success",
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID, UserID: testutils.MockUUID}),
			projectID: expectedProject.ID,
			expectations: func(t *testing.T, ctx context.Context, actual []model.Task, err error) {
				assert.Nil(t, err)
				assert.Equal(t, expectedTasks, actual)
			},
		},
		{
			name:      "project id is not UUID",
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID, UserID: testutils.MockUUID}),
			projectID: "mock",
			expectations: func(t *testing.T, ctx context.Context, actual []model.Task, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrInvalidID, err)
				assert.NotEqual(t, expectedTasks, actual)
			},
		},
		{
			name:      "context error",
			ctx:       testutils.MockCtx,
			projectID: expectedProject.ID,
			expectations: func(t *testing.T, ctx context.Context, actual []model.Task, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err, web.CtxErr())
			},
		},
		{
			name:      "no tenant error",
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: "", UserID: testutils.MockUUID}),
			projectID: expectedProject.ID,
			expectations: func(t *testing.T, ctx context.Context, actual []model.Task, err error) {
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
			list, err := repo.List(tc.ctx, tc.projectID)
			tc.expectations(t, tc.ctx, list, err)
		})
	}
}

func TestTaskRepository_Update(t *testing.T) {
	expectedTenantID := testProjects[0].TenantID
	expectedTask := testTasks[0]

	tests := []struct {
		name         string
		ctx          context.Context
		taskID       string
		update       model.UpdateTask
		expectations func(t *testing.T, ctx context.Context, expected model.Task, actual model.Task, err error)
	}{
		{
			name:   "successfully updated task title",
			ctx:    web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			taskID: expectedTask.ID,
			update: model.UpdateTask{
				Title: aws.String("Updated"),
			},
			expectations: func(t *testing.T, ctx context.Context, expected model.Task, actual model.Task, err error) {
				assert.Nil(t, err)
				assert.NotEqual(t, expected, actual)
				expected.Title = actual.Title
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:   "successfully updated task points",
			ctx:    web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			taskID: expectedTask.ID,
			update: model.UpdateTask{
				Points: aws.Int(3),
			},
			expectations: func(t *testing.T, ctx context.Context, expected model.Task, actual model.Task, err error) {
				assert.Nil(t, err)
				assert.NotEqual(t, expected, actual)
				expected.Points = actual.Points
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:   "successfully updated task content",
			ctx:    web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			taskID: expectedTask.ID,
			update: model.UpdateTask{
				Content: aws.String("Updated"),
			},
			expectations: func(t *testing.T, ctx context.Context, expected model.Task, actual model.Task, err error) {
				assert.Nil(t, err)
				assert.NotEqual(t, expectedTask, actual)
				expected.Content = actual.Content
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:   "successfully updated task assigned to field",
			ctx:    web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			taskID: expectedTask.ID,
			update: model.UpdateTask{
				AssignedTo: aws.String("Updated"),
			},
			expectations: func(t *testing.T, ctx context.Context, expected model.Task, actual model.Task, err error) {
				assert.Nil(t, err)
				assert.NotEqual(t, expected, actual)
				expected.AssignedTo = actual.AssignedTo
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:   "successfully updated task attachments",
			ctx:    web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			taskID: expectedTask.ID,
			update: model.UpdateTask{
				Attachments: []string{"Updated", "Updated"},
			},
			expectations: func(t *testing.T, ctx context.Context, expected model.Task, actual model.Task, err error) {
				assert.Nil(t, err)
				assert.NotEqual(t, expected, actual)
				expected.Attachments = actual.Attachments
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:   "task id not UUID",
			ctx:    web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			taskID: "mock",
			update: model.UpdateTask{},
			expectations: func(t *testing.T, ctx context.Context, expected model.Task, actual model.Task, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrInvalidID, err)
			},
		},
		{
			name:   "context error",
			ctx:    testutils.MockCtx,
			taskID: expectedTask.ID,
			update: model.UpdateTask{},
			expectations: func(t *testing.T, ctx context.Context, expected model.Task, actual model.Task, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, web.CtxErr(), err)
			},
		},
		{
			name:   "no tenant error",
			ctx:    web.NewContext(testutils.MockCtx, &web.Values{TenantID: ""}),
			taskID: expectedTask.ID,
			update: model.UpdateTask{},
			expectations: func(t *testing.T, ctx context.Context, expected model.Task, actual model.Task, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrNoTenant, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, Close := dbConnect.AsNonRoot()
			defer Close()

			expectedTaskCopy := expectedTask
			repo := repository.NewTaskRepository(zap.NewNop(), db)
			task, err := repo.Update(tc.ctx, tc.taskID, tc.update, time.Now())
			tc.expectations(t, tc.ctx, expectedTaskCopy, task, err)
		})
	}
}

func TestTaskRepository_Delete(t *testing.T) {
	expectedTenantID := testProjects[0].TenantID
	expectedTask := testTasks[0]

	tests := []struct {
		name   string
		ctx    context.Context
		taskID string

		expectations func(t *testing.T, ctx context.Context, repo *repository.TaskRepository, err error)
	}{
		{
			name:   "success",
			ctx:    web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID, UserID: testutils.MockUUID}),
			taskID: expectedTask.ID,
			expectations: func(t *testing.T, ctx context.Context, repo *repository.TaskRepository, err error) {
				assert.Nil(t, err)
				_, err = repo.Retrieve(ctx, expectedTask.ID)
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrNotFound, err)
			},
		},
		{
			name:   "task id not UUID",
			ctx:    web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID, UserID: testutils.MockUUID}),
			taskID: "mock",
			expectations: func(t *testing.T, ctx context.Context, repo *repository.TaskRepository, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrInvalidID, err)
			},
		},
		{
			name:   "context error",
			ctx:    testutils.MockCtx,
			taskID: expectedTask.ID,
			expectations: func(t *testing.T, ctx context.Context, repo *repository.TaskRepository, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, web.CtxErr(), err)
			},
		},
		{
			name:   "no tenant error",
			ctx:    web.NewContext(testutils.MockCtx, &web.Values{TenantID: "", UserID: testutils.MockUUID}),
			taskID: expectedTask.ID,
			expectations: func(t *testing.T, ctx context.Context, repo *repository.TaskRepository, err error) {
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
			err := repo.Delete(tc.ctx, tc.taskID)
			tc.expectations(t, tc.ctx, repo, err)
		})
	}
}
