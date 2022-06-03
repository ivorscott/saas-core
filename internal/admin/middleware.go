package admin

import (
	"net/http"

	"github.com/devpies/core/pkg/web"
)

func loadSession(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func withSession() web.Middleware {
	f := func(before web.Handler) web.Handler {
		h := func(w http.ResponseWriter, r *http.Request) error {
			if !session.Exists(r.Context(), "UserID") {
				web.Redirect(w, r, "/")
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
				web.Redirect(w, r, "/admin")
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
				web.Redirect(w, r, "/")
				return nil
			}
			err := before(w, r)
			return err
		}
		return h
	}
	return f
}
