package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/devpies/devpie-client-core/users/api/handlers"
	"github.com/devpies/devpie-client-core/users/domain/projects"
	"github.com/devpies/devpie-client-core/users/domain/teams"
	"github.com/google/go-cmp/cmp"
)

func TestTeams(t *testing.T) {
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

	tm := TeamsTests{
		app: handlers.API(shutdown, repo, logger, cfg.Web.CorsOrigins,
			cfg.Web.AuthAudience, cfg.Web.AuthDomain, cfg.Web.AuthMAPIAudience, cfg.Web.AuthM2MClient,
			cfg.Web.AuthM2MSecret, cfg.Web.SendgridAPIKey, nil),

		userToken: test.token("testuser@devpie.io", "devpie12345!"),
	}

	t.Run("getTeams200", tm.getTeams200)
	t.Run("postTeam201", tm.postTeam201)
	t.Run("postTeam400", tm.postTeam400)
	t.Run("postTeam404", tm.postTeam404)
}

type TeamsTests struct {
	app       http.Handler
	userToken string
}

func (tm *TeamsTests) getTeams200(t *testing.T) {
	teamIDSeed := "39541c75-ca3e-4e2b-9728-54327772d001"
	teamNameSeed := "TestTeam"

	req := httptest.NewRequest("GET", "/api/v1/users/teams", nil)
	resp := httptest.NewRecorder()

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tm.userToken))

	tm.app.ServeHTTP(resp, req)

	var got []teams.Team
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("error while decoding payload %s", err)
	}

	exp := got[0]
	exp.ID = teamIDSeed
	exp.Name = teamNameSeed

	if diff := cmp.Diff(got[0], exp); diff != "" {
		t.Fatalf("Should get the expected result. Diff %s", diff)
	}
	if resp.Code != http.StatusOK {
		t.Errorf("Should have received status 200 but got %d", resp.Code)
	}
}

func (tm *TeamsTests) postTeam201(t *testing.T) {
	teamName := "My Team"
	projectIDSeed := "8695a94f-7e0a-4198-8c0a-d3e12727a5ba"
	userID := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"

	jsonStr := strings.NewReader(fmt.Sprintf(`{ "name":"%s", "projectId":"%s"}`, teamName, projectIDSeed))
	req := httptest.NewRequest("POST", "/api/v1/users/teams", jsonStr)
	resp := httptest.NewRecorder()

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tm.userToken))

	tm.app.ServeHTTP(resp, req)

	var got teams.Team
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("error while decoding payload %s", err)
	}

	exp := got
	exp.Name = teamName
	exp.UserID = userID

	if diff := cmp.Diff(got, exp); diff != "" {
		t.Fatalf("Should get the expected result. Diff %s", diff)
	}

	if resp.Code != http.StatusCreated {
		t.Errorf("Should have received status 201 but got %d", resp.Code)
	}
}

func (tm *TeamsTests) postTeam404(t *testing.T) {
	teamName := "My Team"
	fakeProjectID := "7695a94f-7e0a-4198-8c0a-d3e12727a5bb"

	jsonStr := strings.NewReader(fmt.Sprintf(`{ "name":"%s", "projectId":"%s"}`, teamName, fakeProjectID))
	req := httptest.NewRequest("POST", "/api/v1/users/teams", jsonStr)
	resp := httptest.NewRecorder()

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tm.userToken))

	tm.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusNotFound {
		t.Errorf("Should have received status 404 but got %d", resp.Code)
	}
}

func (tm *TeamsTests) postTeam400(t *testing.T) {
	teamName := "My Team"
	fakeProjectID := "123"
	want := projects.ErrInvalidID
	jsonStr := strings.NewReader(fmt.Sprintf(`{ "name":"%s", "projectId":"%s"}`, teamName, fakeProjectID))
	req := httptest.NewRequest("POST", "/api/v1/users/teams", jsonStr)
	resp := httptest.NewRecorder()

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tm.userToken))

	tm.app.ServeHTTP(resp, req)

	var e struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
		t.Fatalf("error while decoding payload %s", err)
	}
	if resp.Code != http.StatusBadRequest {
		t.Errorf("Should have received status 400 but got %d", resp.Code)
	}
	if e.Error != want.Error() {
		t.Errorf("Should have received %s, but got %s", want, e.Error)
	}
}
