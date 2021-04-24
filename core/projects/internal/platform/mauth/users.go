package mauth

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const usersEndpoint = "/api/v2/users"

func CreateUser(token *Token, AuthDomain, email string) (*InvitedUser, error) {
	var iu InvitedUser

	connectionType := DatabaseConnection
	defaultPassword := uuid.New().String()

	baseUrl := "https://" + AuthDomain
	resource := usersEndpoint

	uri, err := url.ParseRequestURI(baseUrl)
	if err != nil {
		return nil, err
	}

	uri.Path = resource
	urlStr := uri.String()

	jsonStr := fmt.Sprintf("{ \"email\": \"%s\",\"connection\":\"%s\",\"password\":\"%s\", \"email_verified\": false }", email, connectionType, defaultPassword)
	log.Println(jsonStr)

	req, err := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(jsonStr))
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, &iu); err != nil {
		return nil, err
	}

	log.Printf("response=========== %s", string(body))
	log.Printf("response=========== %+v", iu)

	return &iu, nil
}
