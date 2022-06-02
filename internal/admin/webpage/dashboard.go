package webpage

import (
	"net/http"

	"go.uber.org/zap"
)

func (page *WebPage) Dashboard(w http.ResponseWriter, r *http.Request) {
	if err := page.render.Template(w, r, "dashboard", nil); err != nil {
		page.logger.Error("dashboard", zap.Error(err))
	}
}
