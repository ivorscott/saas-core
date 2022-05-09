package e2e

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/devpies/devpie-client-core/users/api/handlers"
)

func TestUsers(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	cfg, repo, rClose, logger := setupTests(t)
	defer rClose()

	test := Test{
		t:                 t,
		Auth0Domain:       cfg.Web.AuthDomain,
		Auth0Audience:     cfg.Web.AuthAudience,
		Auth0ClientID:     cfg.Web.AuthTestClientID,
		Auth0ClientSecret: cfg.Web.AuthTestClientSecret,
	}

	shutdown := make(chan os.Signal, 1)

	ut := UsersTests{
		app: handlers.API(shutdown, repo, logger, cfg.Web.CorsOrigins,
			cfg.Web.AuthAudience, cfg.Web.AuthDomain, cfg.Web.AuthMAPIAudience, cfg.Web.AuthM2MClient,
			cfg.Web.AuthM2MSecret, cfg.Web.SendgridAPIKey, nil),
		userToken:    test.token("testuser@devpie.io", "devpie12345!"),
		newUserToken: test.token("test@example.com", "devpie12345!"),
	}

	t.Run("getUser200", ut.getUser200)
	t.Run("postUser201", ut.postUser201)
	t.Run("postUser202", ut.postUser202)
}

type UsersTests struct {
	app          http.Handler
	userToken    string
	newUserToken string
}

func (u *UsersTests) getUser200(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/v1/users/me", nil)
	resp := httptest.NewRecorder()

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", u.userToken))

	u.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Should have received status 200 but got %d", resp.Code)
	}
}

func (u *UsersTests) postUser201(t *testing.T) {
	json := strings.NewReader(`{ "auth0Id":"auth0|60a7f8dfa7ebc9006a4c6af4", "email":"test@example.com","emailVerified":false,"firstName":"Test","lastName":"User","picture":"","locale":""}`)
	req := httptest.NewRequest("POST", "/api/v1/users", json)
	resp := httptest.NewRecorder()

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", u.newUserToken))

	u.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Should have received status 201 but got %d", resp.Code)
	}
}

func (u *UsersTests) postUser202(t *testing.T) {
	json := strings.NewReader(`{ "auth0Id":"auth0|60a666916089a00069b2a773", "email":"testuser@devpie.io","emailVerified":false,"firstName":"Test","lastName":"User","picture":"","locale":""}`)
	req := httptest.NewRequest("POST", "/api/v1/users", json)
	resp := httptest.NewRecorder()

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", u.userToken))

	u.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusAccepted {
		t.Errorf("Should have received status 202 but got %d", resp.Code)
	}
}
