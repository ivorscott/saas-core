package mauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const changePasswordEndpoint = "/api/v2/tickets/password-change"

func ChangePasswordTicket(token *Token, AuthDomain string, member User, ttl time.Duration, resultUrl string) (string, error) {
	var passTicket struct{ Ticket string }

	baseUrl := "https://" + AuthDomain

	connId, err := GetConnectionId(token, AuthDomain)
	if err != nil {
		return "", err
	}

	uri, err := url.ParseRequestURI(baseUrl)
	if err != nil {
		return "", err
	}

	uri.Path = changePasswordEndpoint
	urlStr := uri.String()

	jsonStr := fmt.Sprintf("{ \"connection_id\":\"%s\",\"email\":\"%s\",\"result_url\":\"%s\",\"ttl_sec\":%d,\"mark_email_as_verified\":true }", connId, member.Email, resultUrl, ttl.Seconds())

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

	passTicket.Ticket = passTicket.Ticket + "invite"

	return passTicket.Ticket, err
}

