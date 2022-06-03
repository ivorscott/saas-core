package admin

import (
	"net/http"

	"github.com/devpies/core/pkg/web"
)

func loadSession(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func withSession() web.Middleware {
	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {
		// Create the handler that will be attached in the middleware chain.
		h := func(w http.ResponseWriter, r *http.Request) error {
			if !session.Exists(r.Context(), "userID") {
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			}

			// Return the error so it can be handled further up the chain.
			err := before(w, r)

			return err
		}

		return h
	}

	return f
}

func withNoSession() web.Middleware {
	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {
		// Create the handler that will be attached in the middleware chain.
		h := func(w http.ResponseWriter, r *http.Request) error {
			if session.Exists(r.Context(), "userID") {
				http.Redirect(w, r, "/admin", http.StatusTemporaryRedirect)
			}

			// Return the error so it can be handled further up the chain.
			err := before(w, r)

			return err
		}

		return h
	}

	return f
}
