package repository_test

import (
	"flag"
	"os"
	"testing"

	"github.com/devpies/saas-core/internal/project/model"
	"github.com/devpies/saas-core/internal/project/res/testutils"
)

var (
	dbConnect    *testutils.DatabaseClient
	testProjects []model.Project
	testProject  model.Project
	testColumns  []model.Column
	testColumn   model.Column
	testTasks    []model.Task
	testTask     model.Task
)

var port = flag.String("port", "5432", "database port to use in integration tests")

func TestMain(m *testing.M) {
	flag.Parse()
	db, dbClose := testutils.NewDatabaseClient(*port)
	dbConnect = db
	defer dbClose()

	testutils.LoadGoldenFile(&testProjects, "projects.json")
	testutils.LoadGoldenFile(&testProject, "project.json")

	testutils.LoadGoldenFile(&testColumns, "columns.json")
	testutils.LoadGoldenFile(&testColumn, "column.json")

	testutils.LoadGoldenFile(&testTasks, "tasks.json")
	testutils.LoadGoldenFile(&testTask, "task.json")

	os.Exit(m.Run())
}

//func TestGoldenFiles(t *testing.T) {
//	golden := testutils.NewGoldenConfig(false)
//
//	t.Run("task golden files", func(t *testing.T) {
//		t.Run("list", func(t *testing.T) {
//			var actual []model.Task
//			var expected []model.Task
//
//			db, Close := dbConnect.AsRoot()
//			defer Close()
//
//			repo := repository.NewTaskRepository(zap.NewNop(), db)
//
//			ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: testutils.MockUUID})
//
//			for _, v := range testProjects {
//				list, err := repo.List(ctx, v.ID)
//				require.NoError(t, err)
//				actual = append(actual, list...)
//			}
//
//			goldenFile := "tasks.json"
//
//			if golden.ShouldUpdate() {
//				testutils.SaveGoldenFile(&actual, goldenFile)
//			}
//
//			testutils.LoadGoldenFile(&expected, goldenFile)
//			assert.Equal(t, expected, actual)
//		})
//
//		t.Run("retrieve by id", func(t *testing.T) {
//			var actual model.Task
//			var expected model.Task
//
//			db, Close := dbConnect.AsRoot()
//			defer Close()
//
//			repo := repository.NewTaskRepository(zap.NewNop(), db)
//
//			ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: testutils.MockUUID})
//
//			actual, err := repo.Retrieve(ctx, testTasks[0].ID)
//			require.NoError(t, err)
//			goldenFile := "task.json"
//
//			if golden.ShouldUpdate() {
//				testutils.SaveGoldenFile(&actual, goldenFile)
//			}
//
//			testutils.LoadGoldenFile(&expected, goldenFile)
//			assert.Equal(t, expected, actual)
//		})
//	})
//
//	t.Run("column golden files", func(t *testing.T) {
//		t.Run("list", func(t *testing.T) {
//			var actual []model.Column
//			var expected []model.Column
//
//			db, Close := dbConnect.AsRoot()
//			defer Close()
//
//			repo := repository.NewColumnRepository(zap.NewNop(), db)
//
//			ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: testutils.MockUUID})
//
//			for _, v := range testProjects {
//				list, err := repo.List(ctx, v.ID)
//				require.NoError(t, err)
//				actual = append(actual, list...)
//			}
//
//			goldenFile := "columns.json"
//
//			if golden.ShouldUpdate() {
//				testutils.SaveGoldenFile(&actual, goldenFile)
//			}
//
//			testutils.LoadGoldenFile(&expected, goldenFile)
//			assert.Equal(t, expected, actual)
//		})
//
//		t.Run("retrieve by id", func(t *testing.T) {
//			var actual model.Column
//			var expected model.Column
//
//			db, Close := dbConnect.AsRoot()
//			defer Close()
//
//			repo := repository.NewColumnRepository(zap.NewNop(), db)
//
//			ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: testutils.MockUUID})
//
//			actual, err := repo.Retrieve(ctx, testColumns[0].ID)
//			require.NoError(t, err)
//			goldenFile := "column.json"
//
//			if golden.ShouldUpdate() {
//				testutils.SaveGoldenFile(&actual, goldenFile)
//			}
//
//			testutils.LoadGoldenFile(&expected, goldenFile)
//			assert.Equal(t, expected, actual)
//		})
//	})
//
//	t.Run("project golden files", func(t *testing.T) {
//		t.Run("list", func(t *testing.T) {
//			var actual []model.Project
//			var expected []model.Project
//
//			db, Close := dbConnect.AsRoot()
//			defer Close()
//
//			repo := repository.NewProjectRepository(zap.NewNop(), db)
//
//			ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: testutils.MockUUID})
//
//			actual, err := repo.List(ctx)
//			require.NoError(t, err)
//			goldenFile := "projects.json"
//
//			if golden.ShouldUpdate() {
//				testutils.SaveGoldenFile(&actual, goldenFile)
//			}
//
//			testutils.LoadGoldenFile(&expected, goldenFile)
//			assert.Equal(t, expected, actual)
//		})
//
//		t.Run("retrieve by id", func(t *testing.T) {
//			var actual model.Project
//			var expected model.Project
//
//			db, Close := dbConnect.AsRoot()
//			defer Close()
//
//			repo := repository.NewProjectRepository(zap.NewNop(), db)
//
//			ctx := web.NewContext(testutils.MockCtx, &web.Values{TenantID: testutils.MockUUID})
//
//			actual, err := repo.Retrieve(ctx, testProjects[0].ID)
//			require.NoError(t, err)
//			goldenFile := "project.json"
//
//			if golden.ShouldUpdate() {
//				testutils.SaveGoldenFile(&actual, goldenFile)
//			}
//
//			testutils.LoadGoldenFile(&expected, goldenFile)
//			assert.Equal(t, expected, actual)
//		})
//	})
//}
