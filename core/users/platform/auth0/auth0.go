// Auth0 provides authentication and authorization.
package auth0

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/devpies/devpie-client-core/users/platform/web"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Auth0 represents the configuration required for any service to use Auth0.
type Auth0 struct {
	Repo         database.DataStorer
	Domain       string
	Audience     string
	M2MClient    string
	M2MSecret    string
	MAPIAudience string
}

type Auther interface {
	GetPemCert(token *jwt.Token) (string, error)
	GetUserByID(r context.Context) string
	GetUserBySubject(ctx context.Context) string
	GetOrCreateToken() (Token, error)
	GetConnectionID(token Token) (string, error)
	CheckScope(scope, tokenString string) (bool, error)
	ChangePasswordTicket(token Token, user AuthUser, resultURL string) (string, error)
	NewManagementToken() (NewToken, error)
	UpdateUserAppMetaData(token Token, subject, userID string) error
	IsExpired(token Token) bool
	RetrieveToken() (Token, error)
	PersistToken(nt NewToken, now time.Time) (Token, error)
	DeleteToken() error
	CreateUser(token Token, email string) (AuthUser, error)
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

const (
	oauthEndpoint = "/oauth/token"
	usersEndpoint = "/api/v2/users"
	changePasswordEndpoint = "/api/v2/tickets/password-change"
	connectionEndpoint = "/api/v2/connections"
	databaseConnection = "Username-Password-Authentication"
)

// Error codes returned by failures to handle tokens.
var (
	ErrNotFound = errors.New("token not found")
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
						return token, errors.New("invalid audience")
					}
					iss := "https://" + a0.Domain + "/"
					checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
					if !checkIss {
						return token, errors.New("invalid issuer")
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

// GetPemCert takes a token and returns the associated certificate in pem format so it can be parsed.
// It works by grabbing the json web key set and returning the public key certificate.
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

func (a0 *Auth0) GetUserBySubject(ctx context.Context) string {
	claims := ctx.Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
	return fmt.Sprintf("%v", claims["sub"])
}

func (a0 *Auth0) GetUserByID(ctx context.Context) string {
	claims := ctx.Value("user").(*jwt.Token).Claims.(jwt.MapClaims)
	if _, ok := claims["https://client.devpie.io/claims/user_id"]; !ok {
		return ""
	}
	return fmt.Sprintf("%v", claims["https://client.devpie.io/claims/user_id"])
}


// GetOrCreateToken creates a new Token if one does not exist or it returns an existing one.
func (a0 *Auth0) GetOrCreateToken() (Token, error) {
	var t Token

	t, err := a0.RetrieveToken()
	if err == ErrNotFound || a0.IsExpired(t) {
		nt, err := a0.NewManagementToken()
		if err != nil {
			return t, err
		}
		// clean table before persisting
		if err = a0.DeleteToken(); err != nil {
			return t, err
		}
		t, err = a0.PersistToken(nt, time.Now())
		if err != nil {
			return t, err
		}
	}

	return t, nil
}

// NewManagementToken generates a new Auth0 management token and returns it.
func (a0 *Auth0) NewManagementToken() (NewToken, error) {
	var t NewToken

	baseURL := "https://" + a0.Domain
	resource := oauthEndpoint

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", a0.M2MClient)
	data.Set("client_secret", a0.M2MSecret)
	data.Set("audience", a0.MAPIAudience)

	uri, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return t, err
	}

	uri.Path = resource
	urlStr := uri.String()

	req, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return t, err
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return t, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return t, err
	}

	err = json.Unmarshal(body, &t)
	if err != nil {
		return t, err
	}

	return t, nil
}

// IsExpired determines whether or not a Token is expired.
func (a0 *Auth0) IsExpired(token Token) bool {
	parsedToken, err := jwt.ParseWithClaims(token.AccessToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := a0.GetPemCert(token)
		if err != nil {
			return true, err
		}
		return jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	})

	if err != nil {
		// error parsing with claims
		return true
	}

	claims, ok := parsedToken.Claims.(CustomClaims)
	if !ok || !parsedToken.Valid {
		// not ok or not valid
		return true
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		// expired
		return true
	}

	return false
}

// RetrieveToken returns the persisted Token if any exists.
func (a0 *Auth0) RetrieveToken() (Token, error) {
	var t Token

	stmt := a0.Repo.Select(
		"ma_token_id",
		"scope",
		"expires_in",
		"access_token",
		"token_type",
		"created_at",
	).From(
		"ma_token",
	).Limit(1)

	q, args, err := stmt.ToSql()
	if err != nil {
		return t, errors.Wrapf(err, "building query: %v", args)
	}

	if err := a0.Repo.Get(&t, q); err != nil {
		if err == sql.ErrNoRows {
			return t, ErrNotFound
		}
		return t, err
	}

	return t, nil
}

// PersistToken persists a new Token and returns it.
func (a0 *Auth0) PersistToken(nt NewToken, now time.Time) (Token, error) {
	t := Token{
		ID:          uuid.New().String(),
		Scope:       nt.Scope,
		ExpiresIn:   nt.ExpiresIn,
		AccessToken: nt.AccessToken,
		TokenType:   nt.TokenType,
		CreatedAt:   now.UTC(),
	}

	stmt := a0.Repo.Insert(
		"ma_token",
	).SetMap(map[string]interface{}{
		"ma_token_id":  uuid.New().String(),
		"scope":        t.Scope,
		"expires_in":   t.ExpiresIn,
		"access_token": t.AccessToken,
		"token_type":   t.TokenType,
		"created_at":   t.CreatedAt,
	})
	if _, err := stmt.Exec(); err != nil {
		return t, errors.Wrapf(err, "inserting token: %v", t)
	}

	return t, nil
}

// DeleteToken deletes a persisted Token.
func (a0 *Auth0) DeleteToken() error {
	stmt := a0.Repo.Delete("ma_token")
	if _, err := stmt.Exec(); err != nil {
		return errors.Wrapf(err, "deleting previous token")
	}
	return nil
}

// ErrInvalidID represents an error when a user id is not a valid uuid.
var ErrInvalidID = errors.New("id provided was not a valid UUID")

func (a0 *Auth0) CreateUser(token Token, email string) (AuthUser, error) {
	var u AuthUser

	connectionType := databaseConnection
	defaultPassword := uuid.New().String()

	baseURL := "https://" + a0.Domain
	resource := usersEndpoint

	uri, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return u, err
	}

	uri.Path = resource
	urlStr := uri.String()

	jsonStr := fmt.Sprintf("{ \"email\": \"%s\",\"connection\":\"%s\",\"password\":\"%s\", \"email_verified\": false, \"verify_email\": false }", email, connectionType, defaultPassword)

	req, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(jsonStr))
	if err != nil {
		return u, err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return u, err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, &u); err != nil {
		return u, err
	}

	return u, nil
}

// UpdateUserAppMetaData updates the auth0 user account with user_id from internal database
func (a0 *Auth0) UpdateUserAppMetaData(token Token, subject, userID string) error {
	if _, err := uuid.Parse(userID); err != nil {
		return ErrInvalidID
	}

	baseURL := "https://" + a0.Domain
	resource := "/api/v2/users/" + subject

	uri, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return err
	}

	uri.Path = resource
	urlStr := uri.String()

	jsonStr := fmt.Sprintf("{\"app_metadata\": { \"id\": \"%s\" }}", userID)

	req, err := http.NewRequest(http.MethodPatch, urlStr, strings.NewReader(jsonStr))
	if err != nil {
		return err
	}

	bearer := fmt.Sprintf("Bearer %s", token.AccessToken)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", bearer)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (a0 *Auth0) ChangePasswordTicket(token Token, user AuthUser, resultURL string) (string, error) {
	var baseURL = "https://" + a0.Domain
	var fiveDays = 432000 // 5 days in seconds
	var passTicket struct {
		Ticket string
	}

	connID, err := a0.GetConnectionID(token)
	if err != nil {
		return "", err
	}

	uri, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return "", err
	}
	uri.Path = changePasswordEndpoint
	urlStr := uri.String()

	jsonStr := fmt.Sprintf("{ \"connection_id\":\"%s\",\"email\":\"%s\",\"result_url\":\"%s\",\"ttl_sec\":%d,\"mark_email_as_verified\":true }", connID, user.Email, resultURL, fiveDays)
	req, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(jsonStr))
	if err != nil {
		return "", err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	if err = json.Unmarshal(body, &passTicket); err != nil {
		return "", err
	}

	ticket := passTicket.Ticket + "invite"

	return ticket, err
}

func (a0 *Auth0) GetConnectionID(token Token) (string, error) {
	var conn []struct {
		ID   string
		Name string
	}

	urlStr := "https://" + a0.Domain + connectionEndpoint

	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	q := req.URL.Query()
	q.Add("strategy", "auth0")
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if err = json.Unmarshal(body, &conn); err != nil {
		return "", err
	}

	for _, v := range conn {
		if v.Name == databaseConnection {
			return v.ID, nil
		}
	}

	return "", err
}
