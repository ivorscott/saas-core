package auth0

import (
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Auth0 struct {
	Repo         *database.Repository
	Domain       string
	Audience     string
	M2MClient    string
	M2MSecret    string
	MAPIAudience string
}

type Token struct {
	ID          string    `db:"ma_token_id"`
	AccessToken string    `db:"access_token"`
	Scope       string    `db:"scope"`
	ExpiresIn   int       `db:"expires_in"`
	TokenType   string    `db:"token_type"`
	CreatedAt   time.Time `db:"created_at"`
}

type NewToken struct {
	AccessToken string `db:"access_token" json:"access_token"`
	Scope       string `db:"scope" json:"scope"`
	ExpiresIn   int    `db:"expires_in" json:"expires_in"`
	TokenType   string `db:"token_type" json:"token_type"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type CustomClaims struct {
	Scope string `json:"scope"`
	jwt.StandardClaims
}
