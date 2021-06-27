package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	mockPub "github.com/devpies/devpie-client-core/users/api/publishers/mocks"
	"github.com/devpies/devpie-client-core/users/domain/memberships"
	mockQuery "github.com/devpies/devpie-client-core/users/domain/mocks"
	"github.com/devpies/devpie-client-core/users/domain/projects"
	"github.com/devpies/devpie-client-core/users/domain/teams"
	mockAuth "github.com/devpies/devpie-client-core/users/platform/auth0/mocks"
	th "github.com/devpies/devpie-client-core/users/platform/testhelpers"
	"github.com/devpies/devpie-client-core/users/platform/web"
	"github.com/devpies/devpie-client-events/go/events"
	"github.com/go-chi/chi"
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
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.project.(*mockQuery.ProjectQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, nt.ProjectID).Return(projects.ProjectCopy{}, nil)
	fake.query.team.(*mockQuery.TeamQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, nt, uid, mock.AnythingOfType("time.Time")).Return(tm, nil)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, nm, mock.AnythingOfType("time.Time")).Return(m, nil)

	up := projects.UpdateProjectCopy{
		TeamID: &tm.ID,
	}

	fake.query.project.(*mockQuery.ProjectQuerier).On("Update", mock.AnythingOfType("*context.valueCtx"), fake.repo, nt.ProjectID, up).Return(nil)
	fake.publish.(*mockPub.Publisher).On("MembershipCreatedForProject", fake.nats, m, nt.ProjectID, uid).Return(nil)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ = fake.Create(w, r)
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(teamJson(nt)))
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
		mux := chi.NewMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			err := fake.Create(w, r)

			t.Run(fmt.Sprintf("Assert Handler Response/%s", v.name), func(t *testing.T) {
				assert.NotNil(t, err)
			})
		})

		// make request
		writer := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(v.arg))
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
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.project.(*mockQuery.ProjectQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, nt.ProjectID).Return(projects.ProjectCopy{}, projects.ErrNotFound)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(teamJson(nt)))
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
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.project.(*mockQuery.ProjectQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, nt.ProjectID).Return(projects.ProjectCopy{}, projects.ErrInvalidID)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(teamJson(nt)))
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
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.project.(*mockQuery.ProjectQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, nt.ProjectID).Return(projects.ProjectCopy{}, cause)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := fake.Create(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
			assert.Equal(t, fmt.Sprintf("failed to retrieve project: %s", cause), err.Error())
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(teamJson(nt)))
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
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.project.(*mockQuery.ProjectQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, nt.ProjectID).Return(projects.ProjectCopy{}, nil)
	fake.query.team.(*mockQuery.TeamQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, nt, uid, mock.AnythingOfType("time.Time")).Return(teams.Team{}, cause)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := fake.Create(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(teamJson(nt)))
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
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.project.(*mockQuery.ProjectQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, nt.ProjectID).Return(projects.ProjectCopy{}, nil)
	fake.query.team.(*mockQuery.TeamQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, nt, uid, mock.AnythingOfType("time.Time")).Return(tm, nil)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, nm, mock.AnythingOfType("time.Time")).Return(memberships.Membership{}, cause)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := fake.Create(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(teamJson(nt)))
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
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.project.(*mockQuery.ProjectQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, nt.ProjectID).Return(projects.ProjectCopy{}, nil)
	fake.query.team.(*mockQuery.TeamQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, nt, uid, mock.AnythingOfType("time.Time")).Return(tm, nil)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, nm, mock.AnythingOfType("time.Time")).Return(m, nil)

	up := projects.UpdateProjectCopy{
		TeamID: &tm.ID,
	}

	fake.query.project.(*mockQuery.ProjectQuerier).On("Update", mock.AnythingOfType("*context.valueCtx"), fake.repo, nt.ProjectID, up).Return(cause)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := fake.Create(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(teamJson(nt)))
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
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.project.(*mockQuery.ProjectQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, nt.ProjectID).Return(projects.ProjectCopy{}, nil)
	fake.query.team.(*mockQuery.TeamQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, nt, uid, mock.AnythingOfType("time.Time")).Return(tm, nil)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, nm, mock.AnythingOfType("time.Time")).Return(m, nil)

	up := projects.UpdateProjectCopy{
		TeamID: &tm.ID,
	}

	fake.query.project.(*mockQuery.ProjectQuerier).On("Update", mock.AnythingOfType("*context.valueCtx"), fake.repo, nt.ProjectID, up).Return(nil)
	fake.publish.(*mockPub.Publisher).On("MembershipCreatedForProject", fake.nats, m, nt.ProjectID, uid).Return(cause)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := fake.Create(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(teamJson(nt)))
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.project.(*mockQuery.ProjectQuerier).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
		fake.query.membership.(*mockQuery.MembershipQuerier).AssertExpectations(t)
		fake.publish.(*mockPub.Publisher).AssertExpectations(t)
	})
}

func TestTeams_AssignExistingTeam_200(t *testing.T) {
	pid := "40541c75-ed7a-4a2b-8788-88322221c000"
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	tm := team()

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.team.(*mockQuery.TeamQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, tm.ID).Return(tm, nil)

	fake.query.project.(*mockQuery.ProjectQuerier).
		On("Update", mock.AnythingOfType("*context.valueCtx"), fake.repo, pid, mock.AnythingOfType("projects.UpdateProjectCopy")).
		Return(nil)
	fake.publish.(*mockPub.Publisher).On("ProjectUpdated", fake.nats, mock.AnythingOfType("*string"), pid, uid).Return(nil)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}/{pid}", func(w http.ResponseWriter, r *http.Request) {
		_ = fake.AssignExistingTeam(w, r)
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", tm.ID, pid), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Handler Response", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.project.(*mockQuery.ProjectQuerier).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
		fake.publish.(*mockPub.Publisher).AssertExpectations(t)
	})
}

func TestTeams_AssignExistingTeam_400(t *testing.T) {
	pid := "40541c75-ed7a-4a2b-8788-88322221c000"
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	tm := team()

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.team.(*mockQuery.TeamQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, tm.ID).Return(tm, nil)

	fake.query.project.(*mockQuery.ProjectQuerier).
		On("Update", mock.AnythingOfType("*context.valueCtx"), fake.repo, pid, mock.AnythingOfType("projects.UpdateProjectCopy")).
		Return(nil)
	fake.publish.(*mockPub.Publisher).On("ProjectUpdated", fake.nats, mock.AnythingOfType("*string"), pid, uid).Return(nil)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}/{pid}", func(w http.ResponseWriter, r *http.Request) {
		_ = fake.AssignExistingTeam(w, r)
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", tm.ID, pid), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Handler Response", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.project.(*mockQuery.ProjectQuerier).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
		fake.publish.(*mockPub.Publisher).AssertExpectations(t)
	})
}
