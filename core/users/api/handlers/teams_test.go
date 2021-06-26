package handlers

import (
	"context"
	"github.com/devpies/devpie-client-core/users/domain/invites"
	"github.com/devpies/devpie-client-core/users/domain/projects"
	"github.com/devpies/devpie-client-core/users/domain/teams"
	mockAuth "github.com/devpies/devpie-client-core/users/platform/auth0/mocks"
	"github.com/devpies/devpie-client-core/users/platform/database"
	th "github.com/devpies/devpie-client-core/users/platform/testhelpers"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type teamQueryMock struct {
	mock.Mock
}

func (m *teamQueryMock) Create(ctx context.Context, repo database.Storer, nt teams.NewTeam, uid string, now time.Time) (teams.Team, error) {
	args := m.Called(ctx, repo, nt, uid, now)
	return args.Get(0).(teams.Team), args.Error(1)
}

func (m *teamQueryMock) Retrieve(ctx context.Context, repo database.Storer, tid string) (teams.Team, error) {
	args := m.Called(ctx, repo, tid)
	return args.Get(0).(teams.Team), args.Error(1)
}

func (m *teamQueryMock) List(ctx context.Context, repo database.Storer, uid string) ([]teams.Team, error) {
	args := m.Called(ctx, repo, uid)
	return args.Get(0).([]teams.Team), args.Error(1)
}

type projectQueryMock struct {
	mock.Mock
}

func (m *projectQueryMock) Create(ctx context.Context, repo *database.Repository, p projects.ProjectCopy) error {
	args := m.Called(ctx, repo, p)
	return args.Error(0)
}
func (m *projectQueryMock) Retrieve(ctx context.Context, repo database.Storer, pid string) (projects.ProjectCopy, error) {
	args := m.Called(ctx, repo, pid)
	return args.Get(0).(projects.ProjectCopy), args.Error(1)

}
func (m *projectQueryMock) Update(ctx context.Context, repo database.Storer, pid string, update projects.UpdateProjectCopy) error {
	args := m.Called(ctx, repo, pid, update)
	return args.Error(0)

}
func (m *projectQueryMock) Delete(ctx context.Context, repo database.Storer, pid string) error {
	args := m.Called(ctx, repo, pid)
	return args.Error(0)

}

type inviteQueryMock struct {
	mock.Mock
}

func (m *inviteQueryMock) Create(ctx context.Context, repo database.Storer, ni invites.NewInvite, now time.Time) (invites.Invite, error) {
	args := m.Called(ctx, repo, ni, now)
	return args.Get(0).(invites.Invite), args.Error(1)
}
func (m *inviteQueryMock) RetrieveInvite(ctx context.Context, repo database.Storer, uid string, iid string) (invites.Invite, error) {
	args := m.Called(ctx, repo)
	return args.Get(0).(invites.Invite), args.Error(1)
}
func (m *inviteQueryMock) RetrieveInvites(ctx context.Context, repo database.Storer, uid string) ([]invites.Invite, error) {
	args := m.Called(ctx, repo, uid)
	return args.Get(0).([]invites.Invite), args.Error(1)
}
func (m *inviteQueryMock) Update(ctx context.Context, repo database.Storer, update invites.UpdateInvite, uid, iid string, now time.Time) (invites.Invite, error) {
	args := m.Called(ctx, repo, uid, iid, now)
	return args.Get(0).(invites.Invite), args.Error(1)
}

type TeamDeps struct {
	service *Team
	repo    *database.Repository
	auth0   *mockAuth.Auther
	query   TeamQueries
}

func setupTeamMocks() *TeamDeps {
	mockRepo := th.Repo()
	mockAuth0 := &mockAuth.Auther{}
	mockTeamQueries := &teamQueryMock{}
	mockProjectQueries := &projectQueryMock{}
	mockMembershipQueries := &membershipQueryMock{}
	mockUserQueries := &userQueryMock{}
	mockInviteQueries := &inviteQueryMock{}

	tq := TeamQueries{mockTeamQueries, mockProjectQueries, mockMembershipQueries, mockUserQueries, mockInviteQueries}

	return &TeamDeps{
		repo:  mockRepo,
		auth0: mockAuth0,
		query: tq,
		service: &Team{
			repo:  mockRepo,
			auth0: mockAuth0,
			query: tq,
		},
	}
}

func TestTeams_Create_200(t *testing.T) {
	t.Skip()
	setupTeamMocks()
}