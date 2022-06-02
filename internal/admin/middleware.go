package adminclient

import (
	"net/http"
)

func loadSession(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func withSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !session.Exists(r.Context(), "userID") {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		}
		next.ServeHTTP(w, r)
	})
}

//func withAuth(log *zap.Logger, cfg config.Config, next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		err := verifyToken(w, r, cfg)
//		if err != nil {
//			err = web.Respond(r.Context(), w, nil, http.StatusUnauthorized)
//			if err != nil {
//				log.Error(err.Error())
//			}
//		}
//		next.ServeHTTP(w, r)
//	})
//}
//
//func verifyToken(w http.ResponseWriter, r *http.Request, cfg config.Config) error {
//	authHeader := r.Header.Get("Authorization")
//	splitAuthHeader := strings.Split(authHeader, " ")
//
//	if len(splitAuthHeader) != 2 {
//		return fmt.Errorf("missing or invalid authorization header")
//	}
//
//	pubKeyURL := "https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json"
//	formattedURL := fmt.Sprintf(pubKeyURL, cfg.Cognito.Region, cfg.Cognito.UserPoolClientID)
//
//	keySet, err := jwk.Fetch(r.Context(), formattedURL)
//	if err != nil {
//		return err
//	}
//
//	_, err = jwt.Parse(
//		[]byte(splitAuthHeader[1]),
//		jwt.WithKeySet(keySet),
//		jwt.WithValidate(true),
//	)
//	return err
//}
