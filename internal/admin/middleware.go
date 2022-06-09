package admin

import (
	"net/http"
	"strings"

	"github.com/devpies/saas-core/pkg/web"

	"go.uber.org/zap"
)

func loadSession(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func withSession() web.Middleware {
	f := func(before web.Handler) web.Handler {
		h := func(w http.ResponseWriter, r *http.Request) error {
			if !session.Exists(r.Context(), "UserID") {
				web.Redirect(w, r, "/", http.StatusSeeOther)
				return nil
			}
			err := before(w, r)
			return err
		}
		return h
	}
	return f
}

func withNoSession() web.Middleware {
	f := func(before web.Handler) web.Handler {
		h := func(w http.ResponseWriter, r *http.Request) error {
			if session.Exists(r.Context(), "UserID") {
				web.Redirect(w, r, "/admin", http.StatusSeeOther)
				return nil
			}
			err := before(w, r)
			return err
		}
		return h
	}
	return f
}

func withPasswordChallengeSession() web.Middleware {
	f := func(before web.Handler) web.Handler {
		h := func(w http.ResponseWriter, r *http.Request) error {
			if !session.Exists(r.Context(), "PasswordChallenge") {
				web.Redirect(w, r, "/", http.StatusSeeOther)
				return nil
			}
			err := before(w, r)
			return err
		}
		return h
	}
	return f
}

func withAuth(log *zap.Logger, region string, UserPoolClientID string) web.Middleware {
	f := func(after web.Handler) web.Handler {
		h := func(w http.ResponseWriter, r *http.Request) error {
			if strings.Contains(r.URL.Path, "/admin/api/") {
				r, err := web.Authenticate(log, r, region, UserPoolClientID)
				if err != nil {
					web.Redirect(w, r, "/admin/logout", http.StatusSeeOther)
					return web.NewRequestError(err, http.StatusUnauthorized)
				}
			}
			return after(w, r)
		}
		return h
	}
	return f
}
