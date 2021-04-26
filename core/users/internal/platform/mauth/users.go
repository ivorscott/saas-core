package mauth

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const usersEndpoint = "/api/v2/users"

type User struct{
	ID string
	Email string
}

func CreateUser(token *Token, AuthDomain, email string) (User, error) {
	var u User

	connectionType := DatabaseConnection
	defaultPassword := uuid.New().String()

	baseUrl := "https://" + AuthDomain
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
