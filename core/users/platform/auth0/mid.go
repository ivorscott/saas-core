package auth0

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/devpies/devpie-client-core/users/domain/users"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"

	"github.com/devpies/devpie-client-core/users/platform/web"
)

// Authenticate middleware verifies the access token sent from auth0
func (a0 *Auth0) Authenticate() web.Middleware {
	// this is the actual middleware function to be executed.
	f := func(after web.Handler) web.Handler {
		// create the handler that will be attached in the middleware chain.
		h := func(w http.ResponseWriter, r *http.Request) error {

			jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
				Debug: false,
				ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
					checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(a0.Audience, false)
					if !checkAud {
						return token, errors.New("invalid audience.")
					}
					iss := "https://" + a0.Domain + "/"
					checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
					if !checkIss {
						return token, errors.New("invalid issuer.")
					}
					cert, err := a0.GetPemCert(token)
					if err != nil {
						return nil, err
					}
					return jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
				},
				SigningMethod: jwt.SigningMethodRS256,
			})

			if err := jwtMiddleware.CheckJWT(w, r); err != nil {
				return web.NewRequestError(err, http.StatusForbidden)
			}

			return after(w, r)
		}

		return h
	}

	return f
}

// CheckScope middleware verifies the access token has the correct scope before returning a successful response
func (a0 *Auth0) CheckScope(scope, tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := a0.GetPemCert(token)
		if err != nil {
			return nil, err
		}

		// Parse pem encoded pkcs1 or pkcs8 public key
		return jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	})
	if err != nil {
		return false, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	hasScope := false

	if ok && token.Valid {
		result := strings.Split(claims.Scope, " ")
		for i := range result {
			if result[i] == scope {
				hasScope = true
			}
		}
	}
	return hasScope, nil
}

// You need to create a function that grabs the json web key set and returns the public key certificate.
// GetPemCert takes a token and returns the associated certificate in pem format so it can be parsed.
func (a0 *Auth0) GetPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://" + a0.Domain + "/.well-known/jwks.json")
	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}

	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}
	return cert, nil
}

type Lookup func(ctx context.Context, repo *database.Repository, aid string)  (users.User,error)

func (a0 *Auth0) GetUser(r *http.Request, cb Lookup) string {
	uid := a0.GetUserById(r)
	if uid == "" {
		user, err := cb(r.Context(), a0.Repo, a0.GetUserBySubject(r))
		if err != nil {
			return ""
		}
		return user.ID
	}
	return uid
}

func (a0 *Auth0) GetUserBySubject(r *http.Request) string {
	claims := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
	return fmt.Sprintf("%v", claims["sub"])
}

func (a0 *Auth0) GetUserById(r *http.Request) string {
	claims := r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
	if _, ok := claims["https://client.devpie.io/claims/user_id"]; !ok {
		return ""
	}
	return fmt.Sprintf("%v", claims["https://client.devpie.io/claims/user_id"])
}
