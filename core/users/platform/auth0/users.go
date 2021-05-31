package auth0

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const usersEndpoint = "/api/v2/users"

// ErrInvalidID represents an error when a user id is not a valid uuid.
var ErrInvalidID = errors.New("id provided was not a valid UUID")

func (a0 *Auth0) CreateUser(token Token, email string) (AuthUser, error) {
	var u AuthUser

	connectionType := DatabaseConnection
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
