package auth0

import (
	"database/sql"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const oauthEndpoint = "/oauth/token"

var (
	ErrNotFound = errors.New("token not found")
)

func (a0 *Auth0) GetOrCreateToken() (Token, error) {
	var t Token

	t, err := a0.RetrieveToken()
	if err == ErrNotFound || a0.IsExpired(t) {
		nt, err := a0.NewManagementToken()
		if err != nil {
			return t, err
		}
		// clean table before persisting
		if err := a0.DeleteToken(); err != nil {
			return t, err
		}
		if err := a0.PersistToken(nt, time.Now()); err != nil {
			return t, err
		}
	}
	return t, nil
}

func (a0 *Auth0) NewManagementToken() (NewToken, error) {
	var t NewToken

	baseUrl := "https://" + a0.Domain
	resource := oauthEndpoint

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", a0.M2MClient)
	data.Set("client_secret", a0.M2MSecret)
	data.Set("audience", a0.MAPIAudience)

	uri, err := url.ParseRequestURI(baseUrl)
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

func (a0 *Auth0) RetrieveToken() (Token, error) {
	var t Token

	stmt := a0.Repo.SQ.Select(
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

	if err := a0.Repo.DB.Get(&t, q); err != nil {
		if err == sql.ErrNoRows {
			return t, ErrNotFound
		}
		return t, err
	}

	return t, nil
}

func (a0 *Auth0) PersistToken(t NewToken, now time.Time) error {

	stmt := a0.Repo.SQ.Insert(
		"ma_token",
	).SetMap(map[string]interface{}{
		"ma_token_id":  uuid.New().String(),
		"scope":        t.Scope,
		"expires_in":   t.ExpiresIn,
		"access_token": t.AccessToken,
		"token_type":   t.TokenType,
		"created_at":   now.UTC(),
	})
	if _, err := stmt.Exec(); err != nil {
		return errors.Wrapf(err, "inserting token: %v", t)
	}
	return nil
}

func (a0 *Auth0) DeleteToken() error {
	stmt := a0.Repo.SQ.Delete("ma_token")
	if _, err := stmt.Exec(); err != nil {
		return errors.Wrapf(err, "deleting previous token")
	}
	return nil
}
