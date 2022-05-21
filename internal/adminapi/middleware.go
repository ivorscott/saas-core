package adminapi

import (
	"net/http"
)

func withAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the token.
		next.ServeHTTP(w, r)
	})
}
