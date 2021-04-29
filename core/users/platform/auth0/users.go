package auth0

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const usersEndpoint = "/api/v2/users"

type User struct{
	ID string
	Email string
}

var ErrInvalidID = errors.New("id provided was not a valid UUID")

func (a0 *Auth0) CreateUser(token Token, email string) (User, error) {
	var u User

	connectionType := DatabaseConnection
	defaultPassword := uuid.New().String()

	baseUrl := "https://" + a0.Domain
	resource := usersEndpoint

	uri, err := url.ParseRequestURI(baseUrl)
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

// Update auth0 user account with user_id from internal database
func (a0 *Auth0) UpdateUserAppMetaData(token Token, subject, userId string) error {
	if _, err := uuid.Parse(userId); err != nil {
		return ErrInvalidID
	}

	baseUrl := "https://" + a0.Domain
	resource := "/api/v2/users/" + subject

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
	log.Println("------------------",token , subject, userId)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	log.Println("------------------",res)
	defer res.Body.Close()

	return nil
}