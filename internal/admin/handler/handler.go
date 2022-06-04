package handler

import (
	"net/http"

	"github.com/devpies/saas-core/internal/admin/render"
)

type renderer interface {
	Template(w http.ResponseWriter, r *http.Request, page string, td *render.TemplateData, partials ...string) error
}
