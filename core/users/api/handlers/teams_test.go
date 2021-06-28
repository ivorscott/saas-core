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
	"github.com/devpies/devpie-client-core/users/domain/invites"
	"github.com/devpies/devpie-client-core/users/domain/memberships"
	mockQuery "github.com/devpies/devpie-client-core/users/domain/mocks"
	"github.com/devpies/devpie-client-core/users/domain/projects"
	"github.com/devpies/devpie-client-core/users/domain/teams"
	"github.com/devpies/devpie-client-core/users/domain/users"
	"github.com/devpies/devpie-client-core/users/platform/auth0"
	mockAuth "github.com/devpies/devpie-client-core/users/platform/auth0/mocks"
	th "github.com/devpies/devpie-client-core/users/platform/testhelpers"
	"github.com/devpies/devpie-client-core/users/platform/web"
	"github.com/devpies/devpie-client-events/go/events"
	"github.com/go-chi/chi"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTeamMocks() *Team {
	return &Team{
		repo:    th.Repo(),
		nats:    &events.Client{},
		auth0:   &mockAuth.Auther{},
		origins: "http://localhost, https://devpie.io",
		sender: func(email *mail.SGMailV3) (*rest.Response, error) {
			var resp *rest.Response
			return resp, nil
		},
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

func authUser() auth0.AuthUser {
	return auth0.AuthUser{
		Auth0ID:   "auth0|60a666916089a00069b2a773",
		Email:     "testuser@devpie.io",
		FirstName: th.StringPointer("testuser"),
		Picture:   th.StringPointer("https://s.gravatar.com/avatar/xxxxxxxxxxxx.png"),
	}
}

func invite() invites.Invite {
	return invites.Invite{
		ID:         "29231d75-ca3e-7a2c-9728-12324445c009",
		TeamID:     "39541c75-ca3e-4e2b-9728-54327772d001",
		UserID:     "a4b54ec1-57f9-4c39-ab53-d936dbb6c177",
		Expiration: time.Now().AddDate(0, 0, 5),
		UpdatedAt:  time.Now(),
		CreatedAt:  time.Now(),
	}
}

func teamJSON(nt teams.NewTeam) string {
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
	request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(teamJSON(nt)))
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
	request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(teamJSON(nt)))
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
	request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(teamJSON(nt)))
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
	request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(teamJSON(nt)))
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
	request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(teamJSON(nt)))
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
	request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(teamJSON(nt)))
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
	request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(teamJSON(nt)))
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
	request, _ := http.NewRequest(http.MethodGet, "/", strings.NewReader(teamJSON(nt)))
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
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
		fake.query.project.(*mockQuery.ProjectQuerier).AssertExpectations(t)
		fake.publish.(*mockPub.Publisher).AssertExpectations(t)
	})
}

func TestTeams_AssignExistingTeam_400_Invalid_Team_ID(t *testing.T) {
	pid := "40541c75-ed7a-4a2b-8788-88322221c000"
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	tid := "123"
	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.team.(*mockQuery.TeamQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, tid).Return(teams.Team{}, teams.ErrInvalidID)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}/{pid}", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.AssignExistingTeam(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, teams.ErrInvalidID))
			assert.Equal(t, http.StatusBadRequest, err.(*web.Error).Status)
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", tid, pid), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
	})
}

func TestTeams_AssignExistingTeam_404_Missing_Team(t *testing.T) {
	pid := "40541c75-ed7a-4a2b-8788-88322221c000"
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	tm := team()

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.team.(*mockQuery.TeamQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, tm.ID).Return(teams.Team{}, teams.ErrNotFound)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}/{pid}", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.AssignExistingTeam(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, teams.ErrNotFound))
			assert.Equal(t, http.StatusNotFound, err.(*web.Error).Status)
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", tm.ID, pid), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
	})
}

func TestTeams_AssignExistingTeam_500_Uncaught_Error_On_Retrieve(t *testing.T) {
	cause := errors.New("something went wrong")

	pid := "40541c75-ed7a-4a2b-8788-88322221c000"
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	tm := team()

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.team.(*mockQuery.TeamQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, tm.ID).Return(teams.Team{}, cause)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}/{pid}", func(w http.ResponseWriter, r *http.Request) {
		err := fake.AssignExistingTeam(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
			assert.Equal(t, fmt.Sprintf("failed to retrieve team: %s", cause), err.Error())
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", tm.ID, pid), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
	})
}

func TestTeams_AssignExistingTeam_400_Invalid_Project_ID(t *testing.T) {
	pid := "40541c75-ed7a-4a2b-8788-88322221c000"
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	tm := team()

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.team.(*mockQuery.TeamQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, tm.ID).Return(tm, nil)

	fake.query.project.(*mockQuery.ProjectQuerier).
		On("Update", mock.AnythingOfType("*context.valueCtx"), fake.repo, pid, mock.AnythingOfType("projects.UpdateProjectCopy")).
		Return(projects.ErrInvalidID)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}/{pid}", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.AssignExistingTeam(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, projects.ErrInvalidID))
			assert.Equal(t, http.StatusBadRequest, err.(*web.Error).Status)
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", tm.ID, pid), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
		fake.query.project.(*mockQuery.ProjectQuerier).AssertExpectations(t)
	})
}

func TestTeams_AssignExistingTeam_400_Missing_Project(t *testing.T) {
	pid := "40541c75-ed7a-4a2b-8788-88322221c000"
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	tm := team()

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.team.(*mockQuery.TeamQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, tm.ID).Return(tm, nil)

	fake.query.project.(*mockQuery.ProjectQuerier).
		On("Update", mock.AnythingOfType("*context.valueCtx"), fake.repo, pid, mock.AnythingOfType("projects.UpdateProjectCopy")).
		Return(projects.ErrNotFound)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}/{pid}", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.AssignExistingTeam(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, projects.ErrNotFound))
			assert.Equal(t, http.StatusNotFound, err.(*web.Error).Status)
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", tm.ID, pid), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
		fake.query.project.(*mockQuery.ProjectQuerier).AssertExpectations(t)
	})
}

func TestTeams_AssignExistingTeam_500_Uncaught_Error_On_Publish(t *testing.T) {
	cause := errors.New("something went wrong")

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
	fake.publish.(*mockPub.Publisher).On("ProjectUpdated", fake.nats, mock.AnythingOfType("*string"), pid, uid).Return(cause)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}/{pid}", func(w http.ResponseWriter, r *http.Request) {
		err := fake.AssignExistingTeam(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", tm.ID, pid), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
		fake.query.project.(*mockQuery.ProjectQuerier).AssertExpectations(t)
		fake.publish.(*mockPub.Publisher).AssertExpectations(t)
	})
}

func TestTeams_LeaveTeam_200(t *testing.T) {
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	mid := "085cb8a0-b221-4a6d-95be-592eb68d5670"
	tm := team()

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("Delete", mock.AnythingOfType("*context.valueCtx"), fake.repo, tm.ID, uid).Return(mid, nil)
	fake.publish.(*mockPub.Publisher).On("MembershipDeleted", fake.nats, mid, uid).Return(nil)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		_ = fake.LeaveTeam(w, r)
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tm.ID), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Handler Response", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
		fake.query.project.(*mockQuery.ProjectQuerier).AssertExpectations(t)
		fake.publish.(*mockPub.Publisher).AssertExpectations(t)
	})
}

func TestTeams_LeaveTeam_400_Invalid_ID(t *testing.T) {
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	tid := "123"

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("Delete", mock.AnythingOfType("*context.valueCtx"), fake.repo, tid, uid).Return("", memberships.ErrInvalidID)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.LeaveTeam(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, memberships.ErrInvalidID))
			assert.Equal(t, http.StatusBadRequest, err.(*web.Error).Status)
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tid), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.membership.(*mockQuery.MembershipQuerier).AssertExpectations(t)
	})
}

func TestTeams_LeaveTeam_404_Missing_Membership(t *testing.T) {
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	tid := "123"

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("Delete", mock.AnythingOfType("*context.valueCtx"), fake.repo, tid, uid).Return("", memberships.ErrInvalidID)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.LeaveTeam(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, memberships.ErrInvalidID))
			assert.Equal(t, http.StatusBadRequest, err.(*web.Error).Status)
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tid), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.membership.(*mockQuery.MembershipQuerier).AssertExpectations(t)
	})
}

func TestTeams_LeaveTeam_500_Uncaught_Error_On_Delete(t *testing.T) {
	cause := errors.New("something went wrong")

	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	mid := "085cb8a0-b221-4a6d-95be-592eb68d5670"
	tm := team()

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("Delete", mock.AnythingOfType("*context.valueCtx"), fake.repo, tm.ID, uid).Return(mid, cause)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		err := fake.LeaveTeam(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
			assert.Equal(t, fmt.Sprintf("failed to delete membership: %s", cause), err.Error())
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tm.ID), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.membership.(*mockQuery.MembershipQuerier).AssertExpectations(t)
	})
}

func TestTeams_LeaveTeam_500_Uncaught_Error_On_Publish(t *testing.T) {
	cause := errors.New("something went wrong")

	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	mid := "085cb8a0-b221-4a6d-95be-592eb68d5670"
	tm := team()

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("Delete", mock.AnythingOfType("*context.valueCtx"), fake.repo, tm.ID, uid).Return(mid, nil)
	fake.publish.(*mockPub.Publisher).On("MembershipDeleted", fake.nats, mid, uid).Return(cause)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		err := fake.LeaveTeam(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tm.ID), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
		fake.query.project.(*mockQuery.ProjectQuerier).AssertExpectations(t)
		fake.publish.(*mockPub.Publisher).AssertExpectations(t)
	})
}

func TestTeams_Retrieve_200(t *testing.T) {
	tm := team()

	//setup mocks
	fake := setupTeamMocks()
	fake.query.team.(*mockQuery.TeamQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, tm.ID).Return(tm, nil)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		_ = fake.Retrieve(w, r)
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tm.ID), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Handler Response", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
	})
}

func TestTeams_Retrieve_400_Invalid_ID(t *testing.T) {
	tid := "123"

	//setup mocks
	fake := setupTeamMocks()
	fake.query.team.(*mockQuery.TeamQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, tid).Return(teams.Team{}, teams.ErrInvalidID)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.Retrieve(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, teams.ErrInvalidID))
			assert.Equal(t, http.StatusBadRequest, err.(*web.Error).Status)
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tid), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
	})
}

func TestTeams_Retrieve_404_Missing_Team(t *testing.T) {
	tm := team()

	//setup mocks
	fake := setupTeamMocks()
	fake.query.team.(*mockQuery.TeamQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, tm.ID).Return(teams.Team{}, teams.ErrNotFound)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.Retrieve(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, teams.ErrNotFound))
			assert.Equal(t, http.StatusNotFound, err.(*web.Error).Status)
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tm.ID), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
	})
}

func TestTeams_Retrieve_500_Uncaught_Error_On_Retrieve(t *testing.T) {
	cause := errors.New("something went wrong")

	tm := team()

	//setup mocks
	fake := setupTeamMocks()
	fake.query.team.(*mockQuery.TeamQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, tm.ID).Return(teams.Team{}, cause)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		err := fake.Retrieve(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
			assert.Equal(t, fmt.Sprintf("failed to retrieve team: %s", cause), err.Error())
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tm.ID), nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
	})
}

func TestTeams_List_200(t *testing.T) {
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	ts := []teams.Team{team(), team()}

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.team.(*mockQuery.TeamQuerier).On("List", mock.AnythingOfType("*context.valueCtx"), fake.repo, uid).Return(ts, nil)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ = fake.List(w, r)
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Handler Response", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
	})
}

func TestTeams_List_400_Invalid_ID(t *testing.T) {
	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return("")
	fake.query.team.(*mockQuery.TeamQuerier).On("List", mock.AnythingOfType("*context.valueCtx"), fake.repo, "").Return([]teams.Team{}, teams.ErrInvalidID)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.List(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, teams.ErrInvalidID))
			assert.Equal(t, http.StatusBadRequest, err.(*web.Error).Status)
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
	})
}

func TestTeams_List_404_Missing_Team(t *testing.T) {
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.team.(*mockQuery.TeamQuerier).On("List", mock.AnythingOfType("*context.valueCtx"), fake.repo, uid).Return([]teams.Team{}, teams.ErrNotFound)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.List(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, teams.ErrNotFound))
			assert.Equal(t, http.StatusNotFound, err.(*web.Error).Status)
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
	})
}

func TestTeams_List_500_Uncaught_Error_On_Retrieve(t *testing.T) {
	cause := errors.New("something went wrong")

	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.team.(*mockQuery.TeamQuerier).On("List", mock.AnythingOfType("*context.valueCtx"), fake.repo, uid).Return([]teams.Team{}, cause)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := fake.List(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
			assert.Equal(t, fmt.Sprintf("failed to retrieve teams: %s", cause), err.Error())
		})
	})

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
	})
}

func TestTeam_CreateInvite_200_New_Users(t *testing.T) {
	tm := team()
	au := authUser()
	nu := newUser()
	u := user()
	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("GenerateToken").Return(auth0.Token{}, nil).Once()
	fake.query.user.(*mockQuery.UserQuerier).On("RetrieveByEmail", fake.repo, mock.AnythingOfType("string")).Return(users.User{}, users.ErrNotFound).Twice()
	fake.auth0.(*mockAuth.Auther).On("CreateUser", auth0.Token{}, mock.AnythingOfType("string")).Return(au, nil).Twice()
	fake.query.user.(*mockQuery.UserQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, nu, mock.AnythingOfType("time.Time")).Return(u, nil).Twice()

	fake.auth0.(*mockAuth.Auther).On("UpdateUserAppMetaData", auth0.Token{}, au.Auth0ID, u.ID).Return(nil).Twice()
	fake.auth0.(*mockAuth.Auther).On("ChangePasswordTicket", auth0.Token{}, au, mock.AnythingOfType("string")).Return("link-to-change-password", nil).Twice()
	fake.query.invite.(*mockQuery.InviteQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, mock.AnythingOfType("invites.NewInvite"), mock.AnythingOfType("time.Time")).Return(invites.Invite{}, nil).Twice()

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		_ = fake.CreateInvite(w, r)
	})

	list := strings.NewReader(`{ "emailList": ["example@devpie.io","example@gmail.com"] }`)

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tm.ID), list)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Handler Response", func(t *testing.T) {
		assert.Equal(t, http.StatusCreated, writer.Code)
	})

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.user.(*mockQuery.UserQuerier).AssertExpectations(t)
		fake.query.invite.(*mockQuery.InviteQuerier).AssertExpectations(t)
	})
}

func TestTeam_CreateInvite_200_Existing_Users(t *testing.T) {
	tm := team()
	u := user()

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("GenerateToken").Return(auth0.Token{}, nil).Once()
	fake.query.user.(*mockQuery.UserQuerier).On("RetrieveByEmail", fake.repo, mock.AnythingOfType("string")).Return(u, nil).Twice()
	fake.query.invite.(*mockQuery.InviteQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, mock.AnythingOfType("invites.NewInvite"), mock.AnythingOfType("time.Time")).Return(invites.Invite{}, nil).Twice()

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		_ = fake.CreateInvite(w, r)
	})

	list := strings.NewReader(`{ "emailList": ["example@devpie.io","example@gmail.com"] }`)

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tm.ID), list)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Handler Response", func(t *testing.T) {
		assert.Equal(t, http.StatusCreated, writer.Code)
	})

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.user.(*mockQuery.UserQuerier).AssertExpectations(t)
		fake.query.invite.(*mockQuery.InviteQuerier).AssertExpectations(t)
	})
}

func TestTeam_CreateInvite_400_Missing_Payload(t *testing.T) {
	tm := team()

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
		mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
			err := fake.Create(w, r)

			t.Run(fmt.Sprintf("Assert Handler Response/%s", v.name), func(t *testing.T) {
				assert.NotNil(t, err)
			})
		})

		// make request
		writer := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tm.ID), strings.NewReader(v.arg))
		mux.ServeHTTP(writer, request)

		t.Run(fmt.Sprintf("Assert Server Response/%s", v.name), func(t *testing.T) {
			assert.Equal(t, http.StatusBadRequest, writer.Code)
		})
	}
}

func TestTeam_CreateInvite_500_Uncaught_Error_On_GenerateToken(t *testing.T) {
	cause := errors.New("something went wrong")

	tm := team()

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("GenerateToken").Return(auth0.Token{}, cause).Once()

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		err := fake.CreateInvite(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
			assert.Equal(t, fmt.Sprintf("failure during token generation: %s", cause), err.Error())
		})
	})

	list := strings.NewReader(`{ "emailList": ["example@devpie.io","example@gmail.com"] }`)

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tm.ID), list)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
	})
}

func TestTeam_CreateInvite_500_Uncaught_Error_On_Create_Auth_User(t *testing.T) {
	cause := errors.New("something went wrong")

	tm := team()
	au := authUser()

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("GenerateToken").Return(auth0.Token{}, nil).Once()
	fake.query.user.(*mockQuery.UserQuerier).On("RetrieveByEmail", fake.repo, mock.AnythingOfType("string")).Return(users.User{}, users.ErrNotFound)
	fake.auth0.(*mockAuth.Auther).On("CreateUser", auth0.Token{}, mock.AnythingOfType("string")).Return(au, cause)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		err := fake.CreateInvite(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
		})
	})

	list := strings.NewReader(`{ "emailList": ["example@devpie.io","example@gmail.com"] }`)

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tm.ID), list)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.user.(*mockQuery.UserQuerier).AssertExpectations(t)

	})
}

func TestTeam_CreateInvite_500_Uncaught_Error_On_Create_User(t *testing.T) {
	cause := errors.New("something went wrong")

	tm := team()
	au := authUser()
	nu := newUser()

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("GenerateToken").Return(auth0.Token{}, nil).Once()
	fake.query.user.(*mockQuery.UserQuerier).On("RetrieveByEmail", fake.repo, mock.AnythingOfType("string")).Return(users.User{}, users.ErrNotFound).Once()
	fake.auth0.(*mockAuth.Auther).On("CreateUser", auth0.Token{}, mock.AnythingOfType("string")).Return(au, nil).Once()
	fake.query.user.(*mockQuery.UserQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, nu, mock.AnythingOfType("time.Time")).Return(users.User{}, cause).Once()

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		err := fake.CreateInvite(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
		})
	})

	list := strings.NewReader(`{ "emailList": ["example@devpie.io"] }`)

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tm.ID), list)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.user.(*mockQuery.UserQuerier).AssertExpectations(t)
	})
}

func TestTeam_CreateInvite_500_Uncaught_Error_On_UpdateUserAppMetaData(t *testing.T) {
	cause := errors.New("something went wrong")

	tm := team()
	au := authUser()
	nu := newUser()
	u := user()
	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("GenerateToken").Return(auth0.Token{}, nil).Once()
	fake.query.user.(*mockQuery.UserQuerier).On("RetrieveByEmail", fake.repo, mock.AnythingOfType("string")).Return(users.User{}, users.ErrNotFound).Once()
	fake.auth0.(*mockAuth.Auther).On("CreateUser", auth0.Token{}, mock.AnythingOfType("string")).Return(au, nil).Once()
	fake.query.user.(*mockQuery.UserQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, nu, mock.AnythingOfType("time.Time")).Return(u, nil).Once()
	fake.auth0.(*mockAuth.Auther).On("UpdateUserAppMetaData", auth0.Token{}, au.Auth0ID, u.ID).Return(cause).Once()

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		err := fake.CreateInvite(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
		})
	})

	list := strings.NewReader(`{ "emailList": ["example@devpie.io"] }`)

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tm.ID), list)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.user.(*mockQuery.UserQuerier).AssertExpectations(t)
	})
}

func TestTeam_CreateInvite_500_Uncaught_Error_On_ChangePasswordTicket(t *testing.T) {
	cause := errors.New("something went wrong")

	tm := team()
	au := authUser()
	nu := newUser()
	u := user()

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("GenerateToken").Return(auth0.Token{}, nil).Once()
	fake.query.user.(*mockQuery.UserQuerier).On("RetrieveByEmail", fake.repo, mock.AnythingOfType("string")).Return(users.User{}, users.ErrNotFound).Once()
	fake.auth0.(*mockAuth.Auther).On("CreateUser", auth0.Token{}, mock.AnythingOfType("string")).Return(au, nil).Once()
	fake.query.user.(*mockQuery.UserQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, nu, mock.AnythingOfType("time.Time")).Return(u, nil).Once()
	fake.auth0.(*mockAuth.Auther).On("UpdateUserAppMetaData", auth0.Token{}, au.Auth0ID, u.ID).Return(nil).Once()
	fake.auth0.(*mockAuth.Auther).On("ChangePasswordTicket", auth0.Token{}, au, mock.AnythingOfType("string")).Return("link-to-change-password", cause).Once()

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		err := fake.CreateInvite(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
		})
	})

	list := strings.NewReader(`{ "emailList": ["example@devpie.io"] }`)

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tm.ID), list)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.user.(*mockQuery.UserQuerier).AssertExpectations(t)
	})
}

func TestTeam_CreateInvite_500_Uncaught_Error_On_Invite_Create(t *testing.T) {
	cause := errors.New("something went wrong")

	tm := team()
	au := authUser()
	nu := newUser()
	u := user()

	//setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("GenerateToken").Return(auth0.Token{}, nil).Once()
	fake.query.user.(*mockQuery.UserQuerier).On("RetrieveByEmail", fake.repo, mock.AnythingOfType("string")).Return(users.User{}, users.ErrNotFound).Once()
	fake.auth0.(*mockAuth.Auther).On("CreateUser", auth0.Token{}, mock.AnythingOfType("string")).Return(au, nil).Once()
	fake.query.user.(*mockQuery.UserQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, nu, mock.AnythingOfType("time.Time")).Return(u, nil).Once()
	fake.auth0.(*mockAuth.Auther).On("UpdateUserAppMetaData", auth0.Token{}, au.Auth0ID, u.ID).Return(nil).Once()
	fake.auth0.(*mockAuth.Auther).On("ChangePasswordTicket", auth0.Token{}, au, mock.AnythingOfType("string")).Return("link-to-change-password", nil).Once()
	fake.query.invite.(*mockQuery.InviteQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, mock.AnythingOfType("invites.NewInvite"), mock.AnythingOfType("time.Time")).Return(invites.Invite{}, cause).Once()

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		err := fake.CreateInvite(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
		})
	})

	list := strings.NewReader(`{ "emailList": ["example@devpie.io"] }`)

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tm.ID), list)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.user.(*mockQuery.UserQuerier).AssertExpectations(t)
	})
}

func TestTeam_CreateInvite_500_Uncaught_Error_On_SendMail(t *testing.T) {
	cause := errors.New("something went wrong")

	tm := team()
	au := authUser()
	nu := newUser()
	u := user()

	//setup mocks
	fake := setupTeamMocks()
	fake.sender = func(email *mail.SGMailV3) (*rest.Response, error) {
		var resp *rest.Response
		return resp, cause
	}
	fake.auth0.(*mockAuth.Auther).On("GenerateToken").Return(auth0.Token{}, nil).Once()
	fake.query.user.(*mockQuery.UserQuerier).On("RetrieveByEmail", fake.repo, mock.AnythingOfType("string")).Return(users.User{}, users.ErrNotFound).Once()
	fake.auth0.(*mockAuth.Auther).On("CreateUser", auth0.Token{}, mock.AnythingOfType("string")).Return(au, nil).Once()
	fake.query.user.(*mockQuery.UserQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, nu, mock.AnythingOfType("time.Time")).Return(u, nil).Once()
	fake.auth0.(*mockAuth.Auther).On("UpdateUserAppMetaData", auth0.Token{}, au.Auth0ID, u.ID).Return(nil).Once()
	fake.auth0.(*mockAuth.Auther).On("ChangePasswordTicket", auth0.Token{}, au, mock.AnythingOfType("string")).Return("link-to-change-password", nil).Once()

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}", func(w http.ResponseWriter, r *http.Request) {
		err := fake.CreateInvite(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
			assert.Equal(t, fmt.Sprintf("failed to send email: %s", cause), err.Error())
		})
	})

	list := strings.NewReader(`{ "emailList": ["example@devpie.io"] }`)

	// make request
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s", tm.ID), list)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.user.(*mockQuery.UserQuerier).AssertExpectations(t)
	})
}

func TestTeam_RetrieveInvites_200(t *testing.T) {
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	iv := []invites.Invite{invite(), invite()}
	tm := team()

	// setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.invite.(*mockQuery.InviteQuerier).On("RetrieveInvites", mock.AnythingOfType("*context.valueCtx"), fake.repo, uid).Return(iv, nil).Once()
	fake.query.team.(*mockQuery.TeamQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, iv[0].TeamID).Return(tm, nil).Twice()

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ = fake.RetrieveInvites(w, r)
	})

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Handler Response", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.invite.(*mockQuery.InviteQuerier).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
	})
}

func TestTeam_RetrieveInvites_400_Invalid_User_ID(t *testing.T) {
	uid := "123"

	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.invite.(*mockQuery.InviteQuerier).On("RetrieveInvites", mock.AnythingOfType("*context.valueCtx"), fake.repo, uid).Return([]invites.Invite{}, invites.ErrInvalidID)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.RetrieveInvites(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, invites.ErrInvalidID))
			assert.Equal(t, http.StatusBadRequest, err.(*web.Error).Status)
		})
	})

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.invite.(*mockQuery.InviteQuerier).AssertExpectations(t)
	})
}

func TestTeam_RetrieveInvites_404_Missing_Invites(t *testing.T) {
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"

	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.invite.(*mockQuery.InviteQuerier).On("RetrieveInvites", mock.AnythingOfType("*context.valueCtx"), fake.repo, uid).Return([]invites.Invite{}, invites.ErrNotFound)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.RetrieveInvites(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, invites.ErrNotFound))
			assert.Equal(t, http.StatusNotFound, err.(*web.Error).Status)
		})
	})

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.invite.(*mockQuery.InviteQuerier).AssertExpectations(t)
	})
}

func TestTeam_RetrieveInvites_500_Uncaught_Error_On_RetrieveInvites(t *testing.T) {
	cause := errors.New("something went wrong")

	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"

	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.invite.(*mockQuery.InviteQuerier).On("RetrieveInvites", mock.AnythingOfType("*context.valueCtx"), fake.repo, uid).Return([]invites.Invite{}, cause)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := fake.RetrieveInvites(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
			assert.Equal(t, fmt.Sprintf("failed to retrieve invites: %s", cause), err.Error())
		})
	})

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.invite.(*mockQuery.InviteQuerier).AssertExpectations(t)
	})
}

func TestTeam_RetrieveInvites_400_Invalid_Team_ID(t *testing.T) {
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	iv := []invites.Invite{invite()}

	// setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.invite.(*mockQuery.InviteQuerier).On("RetrieveInvites", mock.AnythingOfType("*context.valueCtx"), fake.repo, uid).Return(iv, nil)
	fake.query.team.(*mockQuery.TeamQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, iv[0].TeamID).Return(teams.Team{}, teams.ErrInvalidID)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.RetrieveInvites(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, teams.ErrInvalidID))
			assert.Equal(t, http.StatusBadRequest, err.(*web.Error).Status)
		})
	})

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.invite.(*mockQuery.InviteQuerier).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
	})
}

func TestTeam_RetrieveInvites_404_Missing_Team(t *testing.T) {
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	iv := []invites.Invite{invite()}

	// setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.invite.(*mockQuery.InviteQuerier).On("RetrieveInvites", mock.AnythingOfType("*context.valueCtx"), fake.repo, uid).Return(iv, nil)
	fake.query.team.(*mockQuery.TeamQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, iv[0].TeamID).Return(teams.Team{}, teams.ErrNotFound)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.RetrieveInvites(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, teams.ErrNotFound))
			assert.Equal(t, http.StatusNotFound, err.(*web.Error).Status)
		})
	})

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.invite.(*mockQuery.InviteQuerier).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
	})
}

func TestTeam_RetrieveInvites_500_Uncaught_Error_On_Retrieve_Team(t *testing.T) {
	cause := errors.New("something went wrong")

	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	iv := []invites.Invite{invite()}

	// setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.invite.(*mockQuery.InviteQuerier).On("RetrieveInvites", mock.AnythingOfType("*context.valueCtx"), fake.repo, uid).Return(iv, nil)
	fake.query.team.(*mockQuery.TeamQuerier).On("Retrieve", mock.AnythingOfType("*context.valueCtx"), fake.repo, iv[0].TeamID).Return(teams.Team{}, cause)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := fake.RetrieveInvites(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
			assert.Equal(t, fmt.Sprintf("failed to retrieve team: %s", cause), err.Error())
		})
	})

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.invite.(*mockQuery.InviteQuerier).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
	})
}

func TestTeam_UpdateInvite_200(t *testing.T) {
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	tm := team()
	nm := newMembership(tm)
	nm.Role = "editor"
	m := membership(nm)
	iv := invite()

	// setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.invite.(*mockQuery.InviteQuerier).On("Update", mock.AnythingOfType("*context.valueCtx"), fake.repo, invites.UpdateInvite{Accepted: true}, uid, iv.ID, mock.AnythingOfType("time.Time")).Return(iv, nil)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, nm, mock.AnythingOfType("time.Time")).Return(m, nil)
	fake.publish.(*mockPub.Publisher).On("MembershipCreated", fake.nats, m, uid).Return(nil)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}/{iid}", func(w http.ResponseWriter, r *http.Request) {
		_ = fake.UpdateInvite(w, r)
	})

	json := strings.NewReader(`{ "accepted": true }`)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", tm.ID, iv.ID), json)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Handler Response", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, writer.Code)
	})

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.invite.(*mockQuery.InviteQuerier).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
		fake.publish.(*mockPub.Publisher).AssertExpectations(t)
	})
}

func TestTeam_UpdateInvite_400_Missing_Payload(t *testing.T) {
	tm := team()
	iv := invite()

	// setup mocks
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
		mux.HandleFunc("/{tid}/{iid}", func(w http.ResponseWriter, r *http.Request) {
			err := fake.UpdateInvite(w, r)

			t.Run(fmt.Sprintf("Assert Handler Response/%s", v.name), func(t *testing.T) {
				assert.NotNil(t, err)
			})
		})

		writer := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", tm.ID, iv.ID), strings.NewReader(v.arg))
		mux.ServeHTTP(writer, request)
	}
}

func TestTeam_UpdateInvite_400_Invalid_ID(t *testing.T) {
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	tm := team()
	iv := invite()

	// setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.invite.(*mockQuery.InviteQuerier).On("Update", mock.AnythingOfType("*context.valueCtx"), fake.repo, invites.UpdateInvite{Accepted: true}, uid, iv.ID, mock.AnythingOfType("time.Time")).Return(iv, invites.ErrInvalidID)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}/{iid}", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.UpdateInvite(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, invites.ErrInvalidID))
			assert.Equal(t, http.StatusBadRequest, err.(*web.Error).Status)
		})
	})

	json := strings.NewReader(`{ "accepted": true }`)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", tm.ID, iv.ID), json)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.invite.(*mockQuery.InviteQuerier).AssertExpectations(t)
	})
}

func TestTeam_UpdateInvite_404_Missing_Invite(t *testing.T) {
	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	tm := team()
	iv := invite()

	// setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.invite.(*mockQuery.InviteQuerier).On("Update", mock.AnythingOfType("*context.valueCtx"), fake.repo, invites.UpdateInvite{Accepted: true}, uid, iv.ID, mock.AnythingOfType("time.Time")).Return(iv, invites.ErrNotFound)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}/{iid}", func(w http.ResponseWriter, r *http.Request) {
		var webErr *web.Error
		err := fake.UpdateInvite(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.As(err, &webErr))
			assert.True(t, errors.Is(err.(*web.Error).Err, invites.ErrNotFound))
			assert.Equal(t, http.StatusNotFound, err.(*web.Error).Status)
		})
	})

	json := strings.NewReader(`{ "accepted": true }`)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", tm.ID, iv.ID), json)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.invite.(*mockQuery.InviteQuerier).AssertExpectations(t)
	})
}

func TestTeam_UpdateInvite_500_Uncaught_Error_On_Update(t *testing.T) {
	cause := errors.New("something went wrong")

	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	tm := team()
	nm := newMembership(tm)
	nm.Role = "editor"
	iv := invite()

	// setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.invite.(*mockQuery.InviteQuerier).On("Update", mock.AnythingOfType("*context.valueCtx"), fake.repo, invites.UpdateInvite{Accepted: true}, uid, iv.ID, mock.AnythingOfType("time.Time")).Return(iv, cause)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}/{iid}", func(w http.ResponseWriter, r *http.Request) {
		err := fake.UpdateInvite(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
			assert.Equal(t, fmt.Sprintf("failed to update invite: %s", cause), err.Error())
		})
	})

	json := strings.NewReader(`{ "accepted": true }`)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", tm.ID, iv.ID), json)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.invite.(*mockQuery.InviteQuerier).AssertExpectations(t)
	})
}

func TestTeam_UpdateInvite_500_Uncaught_Error_On_Membership_Create(t *testing.T) {
	cause := errors.New("something went wrong")

	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	tm := team()
	nm := newMembership(tm)
	nm.Role = "editor"
	m := membership(nm)
	iv := invite()

	// setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.invite.(*mockQuery.InviteQuerier).On("Update", mock.AnythingOfType("*context.valueCtx"), fake.repo, invites.UpdateInvite{Accepted: true}, uid, iv.ID, mock.AnythingOfType("time.Time")).Return(iv, nil)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, nm, mock.AnythingOfType("time.Time")).Return(m, cause)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}/{iid}", func(w http.ResponseWriter, r *http.Request) {
		err := fake.UpdateInvite(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
			assert.Equal(t, fmt.Sprintf("failed to insert membership: %s", cause), err.Error())
		})
	})

	json := strings.NewReader(`{ "accepted": true }`)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", tm.ID, iv.ID), json)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.invite.(*mockQuery.InviteQuerier).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
	})
}

func TestTeam_UpdateInvite_500_Uncaught_Error_On_Publish(t *testing.T) {
	cause := errors.New("something went wrong")

	uid := "a4b54ec1-57f9-4c39-ab53-d936dbb6c177"
	tm := team()
	nm := newMembership(tm)
	nm.Role = "editor"
	m := membership(nm)
	iv := invite()

	// setup mocks
	fake := setupTeamMocks()
	fake.auth0.(*mockAuth.Auther).On("UserByID", mock.AnythingOfType("*context.valueCtx")).Return(uid)
	fake.query.invite.(*mockQuery.InviteQuerier).On("Update", mock.AnythingOfType("*context.valueCtx"), fake.repo, invites.UpdateInvite{Accepted: true}, uid, iv.ID, mock.AnythingOfType("time.Time")).Return(iv, nil)
	fake.query.membership.(*mockQuery.MembershipQuerier).On("Create", mock.AnythingOfType("*context.valueCtx"), fake.repo, nm, mock.AnythingOfType("time.Time")).Return(m, nil)
	fake.publish.(*mockPub.Publisher).On("MembershipCreated", fake.nats, m, uid).Return(cause)

	// setup server
	mux := chi.NewMux()
	mux.HandleFunc("/{tid}/{iid}", func(w http.ResponseWriter, r *http.Request) {
		err := fake.UpdateInvite(w, r)

		t.Run("Assert Handler Response", func(t *testing.T) {
			assert.True(t, errors.Is(err, cause))
		})
	})

	json := strings.NewReader(`{ "accepted": true }`)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s", tm.ID, iv.ID), json)
	mux.ServeHTTP(writer, request)

	t.Run("Assert Mock Expectations", func(t *testing.T) {
		fake.auth0.(*mockAuth.Auther).AssertExpectations(t)
		fake.query.invite.(*mockQuery.InviteQuerier).AssertExpectations(t)
		fake.query.team.(*mockQuery.TeamQuerier).AssertExpectations(t)
		fake.publish.(*mockPub.Publisher).AssertExpectations(t)
	})
}
