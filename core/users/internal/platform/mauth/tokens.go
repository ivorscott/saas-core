package mauth

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/ivorscott/devpie-client-core/users/internal/platform/database"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type PemHandler func(token *jwt.Token) (string, error)

var (
	ErrNotFound = errors.New("token not found")
)

const oauthEndpoint = "/oauth/token"

func NewManagementToken(Domain, M2MClient, M2MSecret, MAPIAudience string) (*Token, error) {
	baseUrl := "https://" + Domain
	resource := oauthEndpoint

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", M2MClient)
	data.Set("client_secret", M2MSecret)
	data.Set("audience", MAPIAudience)

	uri, err := url.ParseRequestURI(baseUrl)
	if err != nil {
		return nil, err
	}

	uri.Path = resource
	urlStr := uri.String()

	req, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	token := Token{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}


func IsExpired(mt *Token, getPemCert PemHandler) bool {
	token, err := jwt.ParseWithClaims(mt.AccessToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := getPemCert(token)
		if err != nil {
			return true, err
		}
		return jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	})

	if err != nil {
		// error parsing with claims
		return true
	}

	claims, ok := token.Claims.(CustomClaims)
	if !ok || !token.Valid {
		// not ok or not valid
		return true
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		// expired
		return true
	}

	return false
}

func Retrieve(ctx context.Context, repo *database.Repository) (*Token, error) {
	var t Token

	stmt := repo.SQ.Select(
		"ma_token_id",
		"token",
		"created",
	).From(
		"ma_token",
	).Limit(1)

	q, args, err := stmt.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "building query: %v", args)
	}

	if err := repo.DB.GetContext(ctx, &t, q); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &t, nil
}

func Persist(ctx context.Context, repo *database.Repository, nt *Token, now time.Time) error {
	t := Token{
		ID:          uuid.New().String(),
		AccessToken: nt.AccessToken,
		Created:     now.UTC(),
	}

	stmt := repo.SQ.Insert(
		"ma_token",
	).SetMap(map[string]interface{}{
		"ma_token_id": t.ID,
		"token":       t.AccessToken,
		"created":     t.Created,
	})

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrapf(err, "inserting token: %v", nt)
	}

	return nil
}

func Delete(ctx context.Context, repo *database.Repository) error {
	stmt := repo.SQ.Delete("ma_token")

	if _, err := stmt.ExecContext(ctx); err != nil {
		return errors.Wrapf(err, "deleting previous token")
	}

	return nil
}
