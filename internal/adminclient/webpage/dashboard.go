package webpage

import (
	"go.uber.org/zap"
	"net/http"
)

func (page *WebPage) Dashboard(w http.ResponseWriter, r *http.Request) {
	if err := page.render.Template(w, r, "dashboard", nil); err != nil {
		page.logger.Error("dashboard", zap.Error(err))
	}
}
