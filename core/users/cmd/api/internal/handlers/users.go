package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/ivorscott/devpie-client-core/users/internal/mid"
	"github.com/ivorscott/devpie-client-core/users/internal/platform/database"
	"github.com/ivorscott/devpie-client-core/users/internal/platform/mauth"
	"github.com/ivorscott/devpie-client-core/users/internal/platform/web"
	"github.com/ivorscott/devpie-client-core/users/internal/users"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Users struct {
	repo           *database.Repository
	log            *log.Logger
	auth0          *mid.Auth0
	origins        string
	sendgridAPIKey string
}

// Retrieve a single user
func (u *Users) RetrieveMe(w http.ResponseWriter, r *http.Request) error {
	var us users.User
	var err error

	id := u.auth0.GetUserById(r)

	if id == "" {
		us, err = users.RetrieveMeByAuthID(r.Context(), u.repo, u.auth0.GetUserBySubject(r))
	} else {
		us, err = users.RetrieveMe(r.Context(), u.repo, id)
	}

	if err != nil {
		switch err {
		case users.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case users.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for user %q", id)
		}
	}

	return web.Respond(r.Context(), w, us, http.StatusOK)
}

// Create a new user
func (u *Users) Create(w http.ResponseWriter, r *http.Request) error {
	var t *mauth.Token
	sub := u.auth0.GetUserBySubject(r)

	var nu users.NewUser
	if err := web.Decode(r, &nu); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	// create user
	us, err := users.Create(r.Context(), u.repo, nu, sub, time.Now())
	if err != nil {
		return err
	}

	// begin auth0 account update

	// try getting existing auth0 management api token
	t, err = mauth.Retrieve(r.Context(), u.repo)
	if err == mauth.ErrNotFound || u.IsExpired(t) {
		// create new management api token
		t, err = u.NewManagementToken()
		if err != nil {
			return err
		}
		// clean table before persisting
		if err := mauth.Delete(r.Context(), u.repo); err != nil {
			return err
		}
		// persist management api token
		if err := mauth.Persist(r.Context(), u.repo, t, time.Now()); err != nil {
			return err
		}
	}

	// add user_id to app_metadata
	if err := u.UpdateUserAppMetaData(t, sub, us.ID); err != nil {
		return err
	}

	return web.Respond(r.Context(), w, us, http.StatusCreated)
}

// Update auth0 user account with user_id from database
func (u *Users) UpdateUserAppMetaData(token *mauth.Token, aid, userId string) error {

	if _, err := uuid.Parse(userId); err != nil {
		return users.ErrInvalidID
	}

	baseUrl := "https://" + u.auth0.Domain
	resource := "/api/v2/users/" + aid

	uri, err := url.ParseRequestURI(baseUrl)
	if err != nil {
		return err
	}

	uri.Path = resource
	urlStr := uri.String()

	jsonStr := fmt.Sprintf("{\"app_metadata\": { \"id\": \"%s\" }}", userId)

	req, err := http.NewRequest(http.MethodPatch, urlStr, strings.NewReader(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

// Create auth0 management token
func (u *Users) NewManagementToken() (*mauth.Token, error) {
	baseUrl := "https://" + u.auth0.Domain
	resource := "/oauth/token"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", u.auth0.M2MClient)
	data.Set("client_secret", u.auth0.M2MSecret)
	data.Set("audience", u.auth0.MAPIAudience)

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

	token := mauth.Token{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// Check management api token for expiration
func (u *Users) IsExpired(t *mauth.Token) bool {
	token, err := jwt.ParseWithClaims(t.AccessToken, &mauth.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := u.auth0.GetPemCert(token)
		if err != nil {
			return true, err
		}
		return jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	})
	if err != nil {
		u.log.Print("error parsing with claims")
		return true
	}

	claims, ok := token.Claims.(*mauth.CustomClaims)
	if !ok || !token.Valid {
		u.log.Print("not ok or not valid")
		return true
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		u.log.Print("token expired")
		return true
	}

	return false
}
