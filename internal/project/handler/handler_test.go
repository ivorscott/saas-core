//go:generate mockery --all --dir . --case snake --output ../mocks --exported
package handler_test

import (
	"os"
	"testing"

	"github.com/devpies/saas-core/internal/project/model"
	"github.com/devpies/saas-core/internal/project/res/testutils"
)

var (
	testProjects []model.Project
)

func TestMain(m *testing.M) {
	testutils.LoadGoldenFile(&testProjects, "projects.json")
	os.Exit(m.Run())
}
