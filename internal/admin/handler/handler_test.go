//go:generate mockery --quiet --all --dir . --case snake --output ../mocks --exported
package handler_test

import "github.com/devpies/core/internal/admin/render"

var (
	testTemplateData *render.TemplateData
)
