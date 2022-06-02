package webapp

import (
	"net/http"

	"go.uber.org/zap"
)

func (app *WebApp) Dashboard(w http.ResponseWriter, r *http.Request) {
	if err := app.render.Template(w, r, "dashboard", nil); err != nil {
		app.logger.Error("dashboard", zap.Error(err))
	}
}
