package auth0

import (
	"database/sql"
	"encoding/json"
	"fmt"
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

func (a0 *Auth0) GetToken() (Token, error) {
	var t Token
	t, err := a0.Retrieve()
	if err == ErrNotFound || a0.IsExpired(t) {
		// create new token
		t, err = a0.NewManagementToken()
		if err != nil {
			return t, err
		}
		// clean table before persisting
		if err := a0.Delete(); err != nil {
			return t, err
		}
		if err := a0.Persist(t, time.Now()); err != nil {
			return t, err
		}
	}
	return t, nil
}

func (a0 *Auth0) NewManagementToken() (Token, error) {
	var t Token
	fmt.Print("creating new token===============", t)

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
	fmt.Print("token===============", string(body))

	err = json.Unmarshal(body, &t)
	if err != nil {
		return t, err
	}
	fmt.Print("token===============", t)
	return t, nil
}

func (a0 *Auth0) IsExpired(token Token) bool {
	parsedToken, err := jwt.ParseWithClaims(token.AccessToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := a0.CertHandler(token)
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

func (a0 *Auth0) Retrieve() (Token, error) {
	var t Token

	stmt := a0.Repo.SQ.Select(
		"ma_token_id",
		"token",
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

func (a0 *Auth0) Persist(nt Token, now time.Time) error {
	t := Token{
		ID:          uuid.New().String(),
		AccessToken: nt.AccessToken,
	}
	stmt := a0.Repo.SQ.Insert(
		"ma_token",
	).SetMap(map[string]interface{}{
		"ma_token_id": t.ID,
		"token":       t.AccessToken,
		"created_at":  now.UTC(),
	})
	if _, err := stmt.Exec(); err != nil {
		return errors.Wrapf(err, "inserting token: %v", nt)
	}
	return nil
}

func (a0 *Auth0) Delete() error {
	stmt := a0.Repo.SQ.Delete("ma_token")
	if _, err := stmt.Exec(); err != nil {
		return errors.Wrapf(err, "deleting previous token")
	}
	return nil
}
