package auth0

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/devpies/devpie-client-core/users/internal/platform/database"
	"github.com/pkg/errors"
	"time"
)

type Auth0 struct {
	Repo *database.Repository
	Auth0User string
	Domain string
	Audience string
	M2MClient string
	M2MSecret string
	MAPIAudience string
	CertHandler PemHandler
}

type Token struct {
	ID          string    `db:"ma_token_id" json:"id"`
	AccessToken string    `db:"token" json:"accessToken"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
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

type PemHandler func(token *jwt.Token) (string, error)

var (
	ErrNotFound = errors.New("token not found")
)
