package adminclient

// TODO: Use withAuth middleware.
//func withAuth(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if !session.Exists(r.Context(), "userID") {
//			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
//		}
//		next.ServeHTTP(w, r)
//	})
//}
