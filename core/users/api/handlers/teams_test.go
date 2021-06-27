package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/devpies/devpie-client-core/users/domain/memberships"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	mockPub "github.com/devpies/devpie-client-core/users/api/publishers/mocks"
	mockQuery "github.com/devpies/devpie-client-core/users/domain/mocks"
	"github.com/devpies/devpie-client-core/users/domain/projects"
	"github.com/devpies/devpie-client-core/users/domain/teams"
	mockAuth "github.com/devpies/devpie-client-core/users/platform/auth0/mocks"
	th "github.com/devpies/devpie-client-core/users/platform/testhelpers"
	"github.com/devpies/devpie-client-core/users/platform/web"
	"github.com/devpies/devpie-client-events/go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTeamMocks() *Team {
	return &Team{
		repo:  th.Repo(),
		nats:  &events.Client{},
		auth0: &mockAuth.Auther{},
		query: TeamQueries{
			&mockQuery.TeamQuerier{},
			&mockQuery.ProjectQuerier{},
			&mockQuery.MembershipQuerier{},
			&mockQuery.UserQuerier{},
			&mockQuery.InviteQuerier{},
		},
		publish: &mockPub.Publisher{},
	}
}

func newTeam() teams.NewTeam {
	return teams.NewTeam{
		Name:      "TestTeam",
		ProjectID: "8695a94f-7e0a-4198-8c0a-d3e12727a5ba",
	}
}

func team() teams.Team {
	return teams.Team{
		ID:        "39541c75-ca3e-4e2b-9728-54327772d001",
		Name:      "TestTeam",
		UserID:    "a4b54ec1-57f9-4c39-ab53-d936dbb6c177",
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
}

func teamJson(nt teams.NewTeam) string {
	return fmt.Sprintf(`{ "name": "%s", "projectId": "%s" }`,
		nt.Name, nt.ProjectID)
}

func TestTeams_Create_201(t *testing.T) {
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	nt := newTeam()
	tm := team()
	nm := newMembership(tm)
	m := membership(nm)

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", context.Background()).Return(uid)
	fake.query.project.(*mockQuery.ProjectQuerier).On("Retrieve", context.Background(), fake.repo, nt.ProjectID).Return(projects.ProjectCopy{}, nil)
	fake.query.team.(*mockQuery.TeamQuerier).On("Create", context.Background(), fake.repo, nt, uid, mock.AnythingOfType("time.Time")).Return(tm, nil)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("Create", context.Background(), fake.repo, nm, mock.AnythingOfType("time.Time")).Return(m, nil)

	up := projects.UpdateProjectCopy{
		TeamID: &tm.ID,
	}

	fake.query.project.(*mockQuery.ProjectQuerier).On("Update", context.Background(), fake.repo, nt.ProjectID, up).Return(nil)
	fake.publish.(*mockPub.Publisher).On("MembershipCreatedForProject", fake.nats, m, nt.ProjectID, uid).Return(nil)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		_ = fake.Create(w, r)
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/users/me", strings.NewReader(teamJson(nt)))
	mux.ServeHTTP(writer, request)

	t.Run("Assert Handler Response", func(t *testing.T) {
		assert.Equal(t, http.StatusCreated, writer.Code)
	})

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.project.(*mockQuery.ProjectQuerier).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
		fake.query.membership.(*mockQuery.MembershipQuerier).AssertExpectations(t)
		fake.publish.(*mockPub.Publisher).AssertExpectations(t)
	})
}

func TestTeams_Create_400_Missing_Payload(t *testing.T) {
	//setup mocks
	fake := setupTeamMocks()

	testcases := []struct {
		name string
		arg  string
	}{
		{"empty payload", ""},
		{"empty object", "{}"},
	}
	for _, v := range testcases {
		// setup server
		mux := http.NewServeMux()
		mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
			err := fake.Create(w, r)

			t.Run(fmt.Sprintf("Assert Handler Response/%s", v.name), func(t *testing.T) {
				assert.NotNil(t, err)
			})
		})

		// make request
		writer := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "/users/me", strings.NewReader(v.arg))
		mux.ServeHTTP(writer, request)

		t.Run(fmt.Sprintf("Assert Server Response/%s", v.name), func(t *testing.T) {
			assert.Equal(t, http.StatusBadRequest, writer.Code)
		})
	}
}

func TestTeams_Create_404_Missing_Project(t *testing.T) {
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	nt := newTeam()

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", context.Background()).Return(uid)
	fake.query.project.(*mockQuery.ProjectQuerier).On("Retrieve", context.Background(), fake.repo, nt.ProjectID).Return(projects.ProjectCopy{}, projects.ErrNotFound)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.Create(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, projects.ErrNotFound))
			assert.Equal(t, http.StatusNotFound, err.(*web.Error).Status)
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/users/me", strings.NewReader(teamJson(nt)))
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.project.(*mockQuery.ProjectQuerier).AssertExpectations(t)
	})
}

func TestTeams_Create_400_Invalid_Project_ID(t *testing.T) {
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	nt := newTeam()

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", context.Background()).Return(uid)
	fake.query.project.(*mockQuery.ProjectQuerier).On("Retrieve", context.Background(), fake.repo, nt.ProjectID).Return(projects.ProjectCopy{}, projects.ErrInvalidID)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.Create(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, projects.ErrInvalidID))
			assert.Equal(t, http.StatusBadRequest, err.(*web.Error).Status)
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/users/me", strings.NewReader(teamJson(nt)))
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.project.(*mockQuery.ProjectQuerier).AssertExpectations(t)
	})
}

func TestTeams_Create_500_Uncaught_Error_On_Retrieve(t *testing.T) {
	cause := errors.New("something went wrong")

	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	nt := newTeam()

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", context.Background()).Return(uid)
	fake.query.project.(*mockQuery.ProjectQuerier).On("Retrieve", context.Background(), fake.repo, nt.ProjectID).Return(projects.ProjectCopy{}, cause)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		err := fake.Create(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
			assert.Equal(t, fmt.Sprintf("failed retrieving project: %q : %s", nt.ProjectID, cause.Error()), err.Error())
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/users/me", strings.NewReader(teamJson(nt)))
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.project.(*mockQuery.ProjectQuerier).AssertExpectations(t)
	})
}


func TestTeams_Create_500_Uncaught_Error_On_Create(t *testing.T) {
	cause := errors.New("something went wrong")

	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	nt := newTeam()

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", context.Background()).Return(uid)
	fake.query.project.(*mockQuery.ProjectQuerier).On("Retrieve", context.Background(), fake.repo, nt.ProjectID).Return(projects.ProjectCopy{}, nil)
	fake.query.team.(*mockQuery.TeamQuerier).On("Create", context.Background(), fake.repo, nt, uid, mock.AnythingOfType("time.Time")).Return(teams.Team{}, cause)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		err := fake.Create(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/users/me", strings.NewReader(teamJson(nt)))
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.project.(*mockQuery.ProjectQuerier).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
	})
}

func TestTeams_Create_500_Uncaught_Error_On_Membership_Create(t *testing.T) {
	cause := errors.New("something went wrong")

	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	nt := newTeam()
	tm := team()
	nm := newMembership(tm)

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", context.Background()).Return(uid)
	fake.query.project.(*mockQuery.ProjectQuerier).On("Retrieve", context.Background(), fake.repo, nt.ProjectID).Return(projects.ProjectCopy{}, nil)
	fake.query.team.(*mockQuery.TeamQuerier).On("Create", context.Background(), fake.repo, nt, uid, mock.AnythingOfType("time.Time")).Return(tm, nil)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("Create", context.Background(), fake.repo, nm, mock.AnythingOfType("time.Time")).Return(memberships.Membership{}, cause)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		err := fake.Create(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/users/me", strings.NewReader(teamJson(nt)))
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.project.(*mockQuery.ProjectQuerier).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
		fake.query.membership.(*mockQuery.MembershipQuerier).AssertExpectations(t)
	})
}

func TestTeams_Create_500_Uncaught_Error_On_Project_Update(t *testing.T) {
	cause := errors.New("something went wrong")

	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	nt := newTeam()
	tm := team()
	nm := newMembership(tm)
	m := membership(nm)

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", context.Background()).Return(uid)
	fake.query.project.(*mockQuery.ProjectQuerier).On("Retrieve", context.Background(), fake.repo, nt.ProjectID).Return(projects.ProjectCopy{}, nil)
	fake.query.team.(*mockQuery.TeamQuerier).On("Create", context.Background(), fake.repo, nt, uid, mock.AnythingOfType("time.Time")).Return(tm, nil)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("Create", context.Background(), fake.repo, nm, mock.AnythingOfType("time.Time")).Return(m, nil)

	up := projects.UpdateProjectCopy{
		TeamID: &tm.ID,
	}

	fake.query.project.(*mockQuery.ProjectQuerier).On("Update", context.Background(), fake.repo, nt.ProjectID, up).Return(cause)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		err := fake.Create(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/users/me", strings.NewReader(teamJson(nt)))
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.project.(*mockQuery.ProjectQuerier).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
		fake.query.membership.(*mockQuery.MembershipQuerier).AssertExpectations(t)
	})
}

func TestTeams_Create_500_Uncaught_Error_On_Publish(t *testing.T) {
	cause := errors.New("something went wrong")

	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	nt := newTeam()
	tm := team()
	nm := newMembership(tm)
	m := membership(nm)

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", context.Background()).Return(uid)
	fake.query.project.(*mockQuery.ProjectQuerier).On("Retrieve", context.Background(), fake.repo, nt.ProjectID).Return(projects.ProjectCopy{}, nil)
	fake.query.team.(*mockQuery.TeamQuerier).On("Create", context.Background(), fake.repo, nt, uid, mock.AnythingOfType("time.Time")).Return(tm, nil)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("Create", context.Background(), fake.repo, nm, mock.AnythingOfType("time.Time")).Return(m, nil)

	up := projects.UpdateProjectCopy{
		TeamID: &tm.ID,
	}

	fake.query.project.(*mockQuery.ProjectQuerier).On("Update", context.Background(), fake.repo, nt.ProjectID, up).Return(nil)
	fake.publish.(*mockPub.Publisher).On("MembershipCreatedForProject", fake.nats, m, nt.ProjectID, uid).Return(cause)

	// setup server
	mux := http.NewServeMux()
	mux.HandleFunc("/users/me", func(w http.ResponseWriter, r *http.Request) {
		err := fake.Create(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/users/me", strings.NewReader(teamJson(nt)))
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.project.(*mockQuery.ProjectQuerier).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
		fake.query.membership.(*mockQuery.MembershipQuerier).AssertExpectations(t)
		fake.publish.(*mockPub.Publisher).AssertExpectations(t)
	})
}