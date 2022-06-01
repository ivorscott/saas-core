package webpage

import (
	"net/http"

	"go.uber.org/zap"
)

func (page *WebPage) Login(w http.ResponseWriter, r *http.Request) {
	if err := page.render.Template(w, r, "login", nil); err != nil {
		page.logger.Error("login", zap.Error(err))
	}
}

func (page *WebPage) Logout(w http.ResponseWriter, r *http.Request) {}

func (page *WebPage) ForceNewPassword(w http.ResponseWriter, r *http.Request) {
	if err := page.render.Template(w, r, "new-password", nil); err != nil {
		page.logger.Error("new-password", zap.Error(err))
	}
}
