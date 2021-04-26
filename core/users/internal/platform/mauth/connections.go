package mauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const connectionEndpoint = "/api/v2/connections"

func GetConnectionId(token *Token, AuthDomain string) (string, error) {
	var conn []struct {
		ID   string
		Name string
	}

	urlStr := "https://" + AuthDomain + connectionEndpoint

	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", token.AccessToken))
	q := req.URL.Query()
	q.Add("strategy", "auth0")
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &conn); err != nil {
		return "", err
	}

	for _, v := range conn {
		if v.Name == DatabaseConnection {
			return v.ID, nil
		}
	}

	return "", err
}
