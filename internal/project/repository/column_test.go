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

func TestColumnRepository_Create(t *testing.T) {
	expectedTenantID := testProjects[0].TenantID
	expectedProject := testProjects[0]

	tests := []struct {
		name         string
		ctx          context.Context
		columnID     string
		expectations func(t *testing.T, ctx context.Context, repo *repository.ColumnRepository, actual model.Column, err error)
	}{
		{
			name: "success",
			ctx:  web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID, UserID: expectedProject.UserID}),
			expectations: func(t *testing.T, ctx context.Context, repo *repository.ColumnRepository, actual model.Column, err error) {
				assert.Nil(t, err)
				expected, err := repo.Retrieve(ctx, actual.ID)
				assert.Nil(t, err)
				assert.Equal(t, expected, actual)
			},
		},
		{
			name: "context error",
			ctx:  testutils.MockCtx,
			expectations: func(t *testing.T, ctx context.Context, repo *repository.ColumnRepository, actual model.Column, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err, web.CtxErr())
			},
		},
		{
			name: "no tenant error",
			ctx:  web.NewContext(testutils.MockCtx, &web.Values{TenantID: ""}),
			expectations: func(t *testing.T, ctx context.Context, repo *repository.ColumnRepository, actual model.Column, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrNoTenant, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, Close := dbConnect.AsNonRoot()
			defer Close()

			repo := repository.NewColumnRepository(zap.NewNop(), db)
			nc := model.NewColumn{
				Title:      "Testing",
				ColumnName: "column-a",
				ProjectID:  expectedProject.ID,
			}
			newColumn, err := repo.Create(tc.ctx, nc, time.Now())
			tc.expectations(t, tc.ctx, repo, newColumn, err)
		})
	}
}

func TestColumnRepository_Retrieve(t *testing.T) {
	expectedColumn := testColumns[0]

	tests := []struct {
		name         string
		ctx          context.Context
		columnID     string
		expectations func(t *testing.T, expected model.Column, actual model.Column, err error)
	}{
		{
			name:     "success",
			columnID: expectedColumn.ID,
			ctx:      web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedColumn.TenantID}),
			expectations: func(t *testing.T, expected model.Column, actual model.Column, err error) {
				assert.Nil(t, err)
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:     "column id not UUID",
			columnID: "mock",
			ctx:      web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedColumn.TenantID}),
			expectations: func(t *testing.T, expected model.Column, actual model.Column, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrInvalidID, err)
				assert.NotEqual(t, expected, actual)
			},
		},
		{
			name:     "data isolation between tenants",
			columnID: expectedColumn.ID,
			ctx:      web.NewContext(testutils.MockCtx, &web.Values{TenantID: testutils.MockUUID}),
			expectations: func(t *testing.T, expected model.Column, actual model.Column, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrNotFound, err)
				assert.NotEqual(t, expected, actual)
			},
		},
		{
			name:     "not found",
			columnID: testutils.MockUUID,
			ctx:      web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedColumn.TenantID}),
			expectations: func(t *testing.T, expected model.Column, actual model.Column, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrNotFound, err)
				assert.NotEqual(t, expectedColumn, actual)
			},
		},
		{
			name:     "context error",
			columnID: expectedColumn.ID,
			ctx:      testutils.MockCtx,
			expectations: func(t *testing.T, expected model.Column, actual model.Column, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err, web.CtxErr())
			},
		},
		{
			name:     "no tenant error",
			columnID: expectedColumn.ID,
			ctx:      web.NewContext(testutils.MockCtx, &web.Values{TenantID: ""}),
			expectations: func(t *testing.T, expected model.Column, actual model.Column, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrNoTenant, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, Close := dbConnect.AsNonRoot()
			defer Close()

			repo := repository.NewColumnRepository(zap.NewNop(), db)
			column, err := repo.Retrieve(tc.ctx, tc.columnID)
			tc.expectations(t, expectedColumn, column, err)
		})
	}
}

func TestColumnRepository_List(t *testing.T) {
	expectedTenantID := testColumns[0].TenantID
	expectedProjectID := testColumns[0].ProjectID
	expectedColumns := make([]model.Column, 0)

	for _, v := range testColumns {
		if expectedProjectID == v.ProjectID {
			expectedColumns = append(expectedColumns, v)
		}
	}

	tests := []struct {
		name         string
		ctx          context.Context
		projectID    string
		expectations func(t *testing.T, expected []model.Column, actual []model.Column, err error)
	}{
		{
			name:      "success",
			projectID: expectedProjectID,
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			expectations: func(t *testing.T, expected []model.Column, actual []model.Column, err error) {
				assert.Nil(t, err)
				assert.ElementsMatch(t, expected, actual)
			},
		},
		{
			name:      "project id not UUID",
			projectID: "mock",
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			expectations: func(t *testing.T, expected []model.Column, actual []model.Column, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrInvalidID, err)
				assert.NotEqual(t, expected, actual)
			},
		},
		{
			name:      "data isolation between tenants",
			projectID: expectedProjectID,
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: testutils.MockUUID}),
			expectations: func(t *testing.T, expected []model.Column, actual []model.Column, err error) {
				assert.Nil(t, err)
				assert.Equal(t, 0, len(actual))
			},
		},
		{
			name:      "context error",
			projectID: expectedProjectID,
			ctx:       testutils.MockCtx,
			expectations: func(t *testing.T, expected []model.Column, actual []model.Column, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, err, web.CtxErr())
			},
		},
		{
			name:      "no tenant error",
			projectID: expectedProjectID,
			ctx:       web.NewContext(testutils.MockCtx, &web.Values{TenantID: ""}),
			expectations: func(t *testing.T, expected []model.Column, actual []model.Column, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrNoTenant, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, Close := dbConnect.AsNonRoot()
			defer Close()

			repo := repository.NewColumnRepository(zap.NewNop(), db)
			list, err := repo.List(tc.ctx, tc.projectID)
			tc.expectations(t, expectedColumns, list, err)
		})
	}
}

func TestColumnRepository_Update(t *testing.T) {
	expectedTenantID := testColumns[0].TenantID
	expectedColumn := testColumns[0]

	tests := []struct {
		name         string
		ctx          context.Context
		columnID     string
		update       model.UpdateColumn
		expectations func(t *testing.T, ctx context.Context, expected model.Column, actual model.Column, err error)
	}{
		{
			name:     "successfully updated column title",
			ctx:      web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			columnID: expectedColumn.ID,
			update: model.UpdateColumn{
				Title: aws.String("Updated"),
			},
			expectations: func(t *testing.T, ctx context.Context, expected model.Column, actual model.Column, err error) {
				assert.Nil(t, err)
				assert.NotEqual(t, expected, actual)
				expected.Title = actual.Title
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:     "successfully updated task ids",
			ctx:      web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			columnID: expectedColumn.ID,
			update: model.UpdateColumn{
				TaskIDS: []string{"Updated", "Updated"},
			},
			expectations: func(t *testing.T, ctx context.Context, expected model.Column, actual model.Column, err error) {
				assert.Nil(t, err)
				assert.NotEqual(t, expected, actual)
				expected.TaskIDS = actual.TaskIDS
				assert.Equal(t, expected, actual)
			},
		},
		{
			name:     "column id not UUID",
			ctx:      web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID}),
			columnID: "mock",
			update:   model.UpdateColumn{},
			expectations: func(t *testing.T, ctx context.Context, expected model.Column, actual model.Column, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrInvalidID, err)
			},
		},
		{
			name:     "context error",
			ctx:      testutils.MockCtx,
			columnID: expectedColumn.ID,
			update:   model.UpdateColumn{},
			expectations: func(t *testing.T, ctx context.Context, expected model.Column, actual model.Column, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, web.CtxErr(), err)
			},
		},
		{
			name:     "no tenant error",
			ctx:      web.NewContext(testutils.MockCtx, &web.Values{TenantID: ""}),
			columnID: expectedColumn.ID,
			update:   model.UpdateColumn{},
			expectations: func(t *testing.T, ctx context.Context, expected model.Column, actual model.Column, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrNoTenant, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, Close := dbConnect.AsNonRoot()
			defer Close()

			expectedColumnCopy := expectedColumn
			repo := repository.NewColumnRepository(zap.NewNop(), db)
			column, err := repo.Update(tc.ctx, tc.columnID, tc.update, time.Now())
			tc.expectations(t, tc.ctx, expectedColumnCopy, column, err)
		})
	}
}

func TestColumnRepository_Delete(t *testing.T) {
	expectedTenantID := testColumns[0].TenantID
	expectedColumnID := testColumns[0].ID

	tests := []struct {
		name         string
		ctx          context.Context
		columnID     string
		expectations func(t *testing.T, ctx context.Context, repo *repository.ColumnRepository, err error)
	}{
		{
			name:     "success",
			ctx:      web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID, UserID: testutils.MockUUID}),
			columnID: expectedColumnID,
			expectations: func(t *testing.T, ctx context.Context, repo *repository.ColumnRepository, err error) {
				assert.Nil(t, err)
				_, err = repo.Retrieve(ctx, expectedColumnID)
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrNotFound, err)
			},
		},
		{
			name:     "column id not UUID",
			ctx:      web.NewContext(testutils.MockCtx, &web.Values{TenantID: expectedTenantID, UserID: testutils.MockUUID}),
			columnID: "mock",
			expectations: func(t *testing.T, ctx context.Context, repo *repository.ColumnRepository, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrInvalidID, err)
			},
		},
		{
			name:     "context error",
			ctx:      testutils.MockCtx,
			columnID: expectedColumnID,
			expectations: func(t *testing.T, ctx context.Context, repo *repository.ColumnRepository, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, web.CtxErr(), err)
			},
		},
		{
			name:     "no tenant error",
			ctx:      web.NewContext(testutils.MockCtx, &web.Values{TenantID: "", UserID: testutils.MockUUID}),
			columnID: expectedColumnID,
			expectations: func(t *testing.T, ctx context.Context, repo *repository.ColumnRepository, err error) {
				assert.NotNil(t, err)
				assert.Equal(t, fail.ErrNoTenant, err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, Close := dbConnect.AsNonRoot()
			defer Close()

			repo := repository.NewColumnRepository(zap.NewNop(), db)
			err := repo.Delete(tc.ctx, tc.columnID)
			tc.expectations(t, tc.ctx, repo, err)
		})
	}
}
