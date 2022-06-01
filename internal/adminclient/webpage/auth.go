package webpage

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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

func (page *WebPage) CreateSession(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	// Store session
	page.session.Put(r.Context(), "userID", userID)
	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}
