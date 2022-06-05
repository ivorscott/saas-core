package mid

import (
	"fmt"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// APIAuth middleware verifies the id_token.
func APIAuth(log *zap.Logger, region string, userPoolClientID string) web.Middleware {
	// this is the actual middleware function to be executed.
	f := func(after web.Handler) web.Handler {
		// create the handler that will be attached in the middleware chain.
		h := func(w http.ResponseWriter, r *http.Request) error {
			if strings.Contains(r.URL.Path, "/api/") {
				err := verifyToken(w, r, region, userPoolClientID)
				if err != nil {
					log.Info("api authentication failed", zap.Error(err))
					// If verification fails log out the user. This is ok to do because the admin app
					// frontend is the only client supported. In other words, we don't expect
					// additional clients to make requests to the admin api.
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

func verifyToken(w http.ResponseWriter, r *http.Request, region string, userPoolClientID string) error {
	authHeader := r.Header.Get("Authorization")
	splitAuthHeader := strings.Split(authHeader, " ")

	if len(splitAuthHeader) != 2 {
		return fmt.Errorf("missing or invalid authorization header")
	}

	pubKeyURL := "https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json"
	formattedURL := fmt.Sprintf(pubKeyURL, region, userPoolClientID)

	keySet, err := jwk.Fetch(r.Context(), formattedURL)
	if err != nil {
		return err
	}

	_, err = jwt.Parse(
		[]byte(splitAuthHeader[1]),
		jwt.WithKeySet(keySet),
		jwt.WithValidate(true),
	)

	// Add user object to context
	return err
}
