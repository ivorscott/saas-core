package main

import "net/http"

func withAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !session.Exists(r.Context(), "userID") {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		}
		next.ServeHTTP(w, r)
	})
}

// withSession loads and saves session data for the current request.
func withSession(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
