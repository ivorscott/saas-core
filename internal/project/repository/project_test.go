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

func TestProjectRepository_Create(t *testing.T) {
	expectedTenantID := testProjects[0].TenantID
	expectedUserID := testProjects[0].UserID

	tests := []struct {
		name         string
		ctx          context.Context
		columnID     string
		expectations func(t *testing.T, ctx context.Context, repo *repository.ProjectRepository, expected model.NewProject, actual model.Project, err error)
	}{
		{
			name: "success",
			ctx:  web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID, UserID: expectedUserID}),
			expectations: func(t *testing.T, ctx context.Context, repo *repository.ProjectRepository, newProject model.NewProject, actual model.Project, err error) {
				assert.Nil(t, err)
				assert.Equal(t, expectedTenantID, actual.TenantID)
				assert.Equal(t, "MYP-", actual.Prefix)
				assert.Equal(t, []string{"column-1", "column-2", "column-3", "column-4"}, actual.ColumnOrder)
				assert.Equal(t, newProject.Name, actual.Name)

				expected, err := repo.Retrieve(ctx, actual.ID)
				assert.Nil(t, err)
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "user id is not UUID",
			ctx:  web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID, UserID: "mock"}),
			expectations: func(t *testing.T, ctx context.Context, repo *repository.ProjectRepository, expected model.NewProject, actual model.Project, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrInvalidID, err)
				assert.NotEqual(t, expected, actual)
			},
		},
		{
			name: "context error",
			ctx:  testutils.MockCtx,
			expectations: func(t *testing.T, ctx context.Context, repo *repository.ProjectRepository, expected model.NewProject, actual model.Project, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err, web.CtxErr())
			},
		},
		{
			name: "no tenant error",
			ctx:  web.NewContext(testutils.MockCtx, &web.Values{TenantID: "", UserID: expectedUserID}),
			expectations: func(t *testing.T, ctx context.Context, repo *repository.ProjectRepository, expected model.NewProject, actual model.Project, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrNoTenant, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, Close := dbConnect.AsNonRoot()
			defer Close()

			repo := repository.NewProjectRepository(zap.NewNop(), db)
			np := model.NewProject{
				Name: "M123y Project",
			}
			newProject, err := repo.Create(tc.ctx, np, time.Now())
			tc.expectations(t, tc.ctx, repo, np, newProject, err)
		})
	}
}

func TestProjectRepository_Retrieve(t *testing.T) {
	expectedTenantID := testProjects[0].TenantID
	expectedProject := testProjects[0]

	tests := []struct {
		name         string
		ctx          context.Context
		projectID    string
		expectations func(t *testing.T, ctx context.Context, expected model.Project, actual model.Project, err error)
	}{
		{
			name:      "success",
			projectID: expectedProject.ID,
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			expectations: func(t *testing.T, ctx context.Context, expected model.Project, actual model.Project, err error) {
				assert.Nil(t, err)
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:      "project id not UUID",
			projectID: "mock",
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			expectations: func(t *testing.T, ctx context.Context, expected model.Project, actual model.Project, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrInvalidID, err)
				assert.NotEqual(t, expected, actual)
			},
		},
		{
			name:      "data isolation between tenants",
			projectID: expectedProject.ID,
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: testutils.MockUUID}),
			expectations: func(t *testing.T, ctx context.Context, expected model.Project, actual model.Project, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrNotFound, err)
				assert.NotEqual(t, expected, actual)
			},
		},
		{
			name:      "not found",
			projectID: testutils.MockUUID,
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			expectations: func(t *testing.T, ctx context.Context, expected model.Project, actual model.Project, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrNotFound, err)
				assert.NotEqual(t, expected, actual)
			},
		},
		{
			name:      "context error",
			projectID: expectedProject.ID,
			ctx:       testutils.MockCtx,
			expectations: func(t *testing.T, ctx context.Context, expected model.Project, actual model.Project, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err, web.CtxErr())
			},
		},
		{
			name:      "no tenant error",
			projectID: expectedProject.ID,
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: ""}),
			expectations: func(t *testing.T, ctx context.Context, expected model.Project, actual model.Project, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrNoTenant, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, Close := dbConnect.AsNonRoot()
			defer Close()

			repo := repository.NewProjectRepository(zap.NewNop(), db)
			newProject, err := repo.Retrieve(tc.ctx, tc.projectID)
			tc.expectations(t, tc.ctx, expectedProject, newProject, err)
		})
	}
}

func TestProjectRepository_List(t *testing.T) {
	expectedTenantID := testProjects[0].TenantID
	expectedProjects := testProjects

	tests := []struct {
		name         string
		ctx          context.Context
		expectations func(t *testing.T, ctx context.Context, expected []model.Project, actual []model.Project, err error)
	}{
		{
			name: "success",
			ctx:  web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			expectations: func(t *testing.T, ctx context.Context, expected []model.Project, actual []model.Project, err error) {
				assert.Nil(t, err)
				assert.ElementsMatch(t, expected, actual)
			},
		},
		{
			name: "data isolation between tenants",
			ctx:  web.NewContext(testutils.MockCtx, &web.Values{TenantID: testutils.MockUUID}),
			expectations: func(t *testing.T, ctx context.Context, expected []model.Project, actual []model.Project, err error) {
				assert.Nil(t, err)
				assert.NotEqual(t, expected, actual)
				assert.Equal(t, 0, len(actual))
			},
		},
		{
			name: "context error",
			ctx:  testutils.MockCtx,
			expectations: func(t *testing.T, ctx context.Context, expected []model.Project, actual []model.Project, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, web.CtxErr(), err)
			},
		},
		{
			name: "no tenant error",
			ctx:  web.NewContext(testutils.MockCtx, &web.Values{TenantID: ""}),
			expectations: func(t *testing.T, ctx context.Context, expected []model.Project, actual []model.Project, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrNoTenant, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, Close := dbConnect.AsNonRoot()
			defer Close()

			repo := repository.NewProjectRepository(zap.NewNop(), db)
			projects, err := repo.List(tc.ctx)
			tc.expectations(t, tc.ctx, expectedProjects, projects, err)
		})
	}
}

func TestProjectRepository_Update(t *testing.T) {
	expectedTenantID := testProjects[0].TenantID
	expectedProject := testProjects[0]

	tests := []struct {
		name         string
		ctx          context.Context
		projectID    string
		update       model.UpdateProject
		expectations func(t *testing.T, ctx context.Context, expected model.Project, actual model.Project, err error)
	}{
		{
			name:      "successfully updated project name",
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			projectID: expectedProject.ID,
			update: model.UpdateProject{
				Name: aws.String("Updated"),
			},
			expectations: func(t *testing.T, ctx context.Context, expected model.Project, actual model.Project, err error) {
				assert.Nil(t, err)
				assert.NotEqual(t, expected, actual)
				expected.Name = actual.Name
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:      "successfully updated project description",
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			projectID: expectedProject.ID,
			update: model.UpdateProject{
				Description: aws.String("Updated"),
			},
			expectations: func(t *testing.T, ctx context.Context, expected model.Project, actual model.Project, err error) {
				assert.Nil(t, err)
				assert.NotEqual(t, expected, actual)
				expected.Description = actual.Description
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:      "successfully updated project active field",
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			projectID: expectedProject.ID,
			update: model.UpdateProject{
				Active: aws.Bool(false),
			},
			expectations: func(t *testing.T, ctx context.Context, expected model.Project, actual model.Project, err error) {
				assert.Nil(t, err)
				assert.NotEqual(t, expectedProject, actual)
				expected.Active = actual.Active
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:      "successfully updated project public field",
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			projectID: expectedProject.ID,
			update: model.UpdateProject{
				Public: aws.Bool(true),
			},
			expectations: func(t *testing.T, ctx context.Context, expected model.Project, actual model.Project, err error) {
				assert.Nil(t, err)
				assert.NotEqual(t, expectedProject, actual)
				expected.Public = actual.Public
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:      "successfully updated project column order",
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			projectID: expectedProject.ID,
			update: model.UpdateProject{
				ColumnOrder: []string{"column-4", "column-3", "column-2", "column-1"},
			},
			expectations: func(t *testing.T, ctx context.Context, expected model.Project, actual model.Project, err error) {
				assert.Nil(t, err)
				assert.NotEqual(t, expectedProject, actual)
				expected.ColumnOrder = actual.ColumnOrder
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:      "project id not UUID",
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			projectID: "mock",
			update:    model.UpdateProject{},
			expectations: func(t *testing.T, ctx context.Context, expected model.Project, actual model.Project, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrInvalidID, err)
			},
		},
		{
			name:      "context error",
			ctx:       testutils.MockCtx,
			projectID: expectedProject.ID,
			update:    model.UpdateProject{},
			expectations: func(t *testing.T, ctx context.Context, expected model.Project, actual model.Project, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, web.CtxErr(), err)
			},
		},
		{
			name:      "no tenant error",
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: ""}),
			projectID: expectedProject.ID,
			update:    model.UpdateProject{},
			expectations: func(t *testing.T, ctx context.Context, expected model.Project, actual model.Project, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrNoTenant, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, Close := dbConnect.AsNonRoot()
			defer Close()
			expectedProjectCopy := expectedProject
			repo := repository.NewProjectRepository(zap.NewNop(), db)
			project, err := repo.Update(tc.ctx, tc.projectID, tc.update, time.Now())
			tc.expectations(t, tc.ctx, expectedProjectCopy, project, err)
		})
	}
}

func TestProjectRepository_Delete(t *testing.T) {
	expectedTenantID := testProjects[0].TenantID
	expectedProject := testProjects[0]

	tests := []struct {
		name      string
		ctx       context.Context
		projectID string

		expectations func(t *testing.T, ctx context.Context, repo *repository.ProjectRepository, err error)
	}{
		{
			name:      "success",
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID, UserID: expectedProject.UserID}),
			projectID: expectedProject.ID,
			expectations: func(t *testing.T, ctx context.Context, repo *repository.ProjectRepository, err error) {
				assert.Nil(t, err)
				_, err = repo.Retrieve(ctx, expectedProject.ID)
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrNotFound, err)
			},
		},
		{
			name:      "user id not UUID",
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID, UserID: "mock"}),
			projectID: expectedProject.ID,
			expectations: func(t *testing.T, ctx context.Context, repo *repository.ProjectRepository, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrInvalidID, err)
			},
		},
		{
			name:      "project id not UUID",
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID, UserID: expectedProject.UserID}),
			projectID: "mock",
			expectations: func(t *testing.T, ctx context.Context, repo *repository.ProjectRepository, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrInvalidID, err)
			},
		},
		{
			name:      "context error",
			ctx:       testutils.MockCtx,
			projectID: expectedProject.ID,
			expectations: func(t *testing.T, ctx context.Context, repo *repository.ProjectRepository, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, web.CtxErr(), err)
			},
		},
		{
			name:      "no tenant error",
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: "", UserID: expectedProject.UserID}),
			projectID: expectedProject.ID,
			expectations: func(t *testing.T, ctx context.Context, repo *repository.ProjectRepository, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrNoTenant, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, Close := dbConnect.AsNonRoot()
			defer Close()

			repo := repository.NewProjectRepository(zap.NewNop(), db)
			err := repo.Delete(tc.ctx, tc.projectID)
			tc.expectations(t, tc.ctx, repo, err)
		})
	}
}
