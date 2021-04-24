package mauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const changePasswordEndpoint = "/api/v2/tickets/password-change"

func ChangePasswordTicket(token *Token, AuthDomain string, member *InvitedUser, resultUrl string) (string, error) {
	var passTicket struct{ Ticket string }

	baseUrl := "https://" + AuthDomain
	resource := changePasswordEndpoint
	timeToLive := 432000

	connId, err := GetConnectionId(token, AuthDomain)
	if err != nil {
		return "", err
	}

	uri, err := url.ParseRequestURI(baseUrl)
	if err != nil {
		return "", err
	}

	uri.Path = resource
	urlStr := uri.String()

	jsonStr := fmt.Sprintf("{ \"connection_id\":\"%s\",\"email\":\"%s\",\"result_url\":\"%s\",\"ttl_sec\":%d,\"mark_email_as_verified\":true }", connId, member.Email, resultUrl, timeToLive)

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

	if err := json.Unmarshal(body, &passTicket); err != nil {
		return "", err
	}

	return passTicket.Ticket, err
}
