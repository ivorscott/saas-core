package auth0

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const changePasswordEndpoint = "/api/v2/tickets/password-change"

func (a0 *Auth0) ChangePasswordTicket(token Token, user AuthUser, resultUrl string) (string, error) {
	var baseUrl = "https://" + a0.Domain
	var fiveDays = 432000 // 5 days in seconds
	var passTicket struct{
		Ticket string
	}

	connId, err := a0.GetConnectionId(token)
	if err != nil {
		return "", err
	}

	uri, err := url.ParseRequestURI(baseUrl)
	if err != nil {
		return "", err
	}
	uri.Path = changePasswordEndpoint
	urlStr := uri.String()

	jsonStr := fmt.Sprintf("{ \"connection_id\":\"%s\",\"email\":\"%s\",\"result_url\":\"%s\",\"ttl_sec\":%d,\"mark_email_as_verified\":true }", connId, user.Email, resultUrl, fiveDays)
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

	ticket := passTicket.Ticket + "invite"
	log.Println("==========jsonStr",ticket)

	return ticket, err
}

