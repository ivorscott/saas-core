// Auth0 provides authentication and authorization.
package auth0

import (
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Auth0 represents the configuration required for any service to use Auth0.
type Auth0 struct {
	Repo         *database.Repository
	Domain       string
	Audience     string
	M2MClient    string
	M2MSecret    string
	MAPIAudience string
}

// AuthUser represents a freshly created Auth0 user created programmatically.
type AuthUser struct {
	Auth0ID       string  `json:"user_id" `
	Email         string  `json:"email"`
	EmailVerified bool    `json:"email_verified"`
	FirstName     *string `json:"nickname"`
	Picture       *string `json:"picture"`
}

// Token represents a Auth0 management Token persisted in the database.
type Token struct {
	ID          string    `db:"ma_token_id"`
	AccessToken string    `db:"access_token"`
	Scope       string    `db:"scope"`
	ExpiresIn   int       `db:"expires_in"`
	TokenType   string    `db:"token_type"`
	CreatedAt   time.Time `db:"created_at"`
}

// NewToken represents a freshly created Token from the Auth0 management API.
type NewToken struct {
	AccessToken string `db:"access_token" json:"access_token"`
	Scope       string `db:"scope" json:"scope"`
	ExpiresIn   int    `db:"expires_in" json:"expires_in"`
	TokenType   string `db:"token_type" json:"token_type"`
}

// Jwks represents storage for a slice of JSONWebKeys.
type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

// JSONWebKeys represents fields related to the JSON Web Key Set for this API.
// These keys contain the public keys, which will be used to verify JWTs.
type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// CustomClaims extends the Standard JWT Claims with Scope information
type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}
