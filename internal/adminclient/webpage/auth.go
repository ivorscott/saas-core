package handler

import (
	"github.com/devpies/core/internal/adminclient/render"
	"go.uber.org/zap"
	"net/http"
)

type AuthPages struct {
	logger *zap.Logger
	config adminapi.ClientConfig
	render *render.Render
}

func (ah *AuthPages) Login(w http.ResponseWriter, r *http.Request) {
	if err := page.render.Template(w, r, "login", nil); err != nil {
		page.logger.Error("login", zap.Error(err))
	}
}

func (ah *AuthPages) Logout(w http.ResponseWriter, r *http.Request) {}

func (ah *AuthPages) ForceNewPassword(w http.ResponseWriter, r *http.Request) {
	if err := page.render.Template(w, r, "new-password", nil); err != nil {
		page.logger.Error("new-password", zap.Error(err))
	}
}

func (ah *AuthPages) Dashboard(w http.ResponseWriter, r *http.Request) {
	if err := page.render.Template(w, r, "dashboard", nil); err != nil {
		page.logger.Error("dashboard", zap.Error(err))
	}
}
