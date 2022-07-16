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

func TestNewColumnRepository_Create(t *testing.T) {
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
				column, err := repo.Retrieve(ctx, actual.ID)
				assert.Nil(t, err)
				assert.Equal(t, column, actual)
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
			name:     "id not UUID",
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
