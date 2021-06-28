package mocks

import (
	"context"
	"github.com/devpies/devpie-client-core/users/domain/invites"
	"github.com/devpies/devpie-client-core/users/domain/memberships"
	"github.com/devpies/devpie-client-core/users/domain/projects"
	"github.com/devpies/devpie-client-core/users/domain/teams"
	"github.com/devpies/devpie-client-core/users/domain/users"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/stretchr/testify/mock"
	"time"
)

type UserQuerier struct {
	mock.Mock
}

func (m *UserQuerier) Create(ctx context.Context, repo database.Storer, nu users.NewUser, now time.Time) (users.User, error) {
	args := m.Called(ctx, repo, nu, now)
	return args.Get(0).(users.User), args.Error(1)
}

func (m *UserQuerier) RetrieveByEmail(repo database.Storer, email string) (users.User, error) {
	args := m.Called(repo, email)
	return args.Get(0).(users.User), args.Error(1)
}

func (m *UserQuerier) RetrieveMe(ctx context.Context, repo database.Storer, uid string) (users.User, error) {
	args := m.Called(ctx, repo, uid)
	return args.Get(0).(users.User), args.Error(1)
}

func (m *UserQuerier) RetrieveMeByAuthID(ctx context.Context, repo database.Storer, aid string) (users.User, error) {
	args := m.Called(ctx, repo, aid)
	return args.Get(0).(users.User), args.Error(1)
}

type TeamQuerier struct {
	mock.Mock
}

func (m *TeamQuerier) Create(ctx context.Context, repo database.Storer, nt teams.NewTeam, uid string, now time.Time) (teams.Team, error) {
	args := m.Called(ctx, repo, nt, uid, now)
	return args.Get(0).(teams.Team), args.Error(1)
}

func (m *TeamQuerier) Retrieve(ctx context.Context, repo database.Storer, tid string) (teams.Team, error) {
	args := m.Called(ctx, repo, tid)
	return args.Get(0).(teams.Team), args.Error(1)
}

func (m *TeamQuerier) List(ctx context.Context, repo database.Storer, uid string) ([]teams.Team, error) {
	args := m.Called(ctx, repo, uid)
	return args.Get(0).([]teams.Team), args.Error(1)
}

type ProjectQuerier struct {
	mock.Mock
}

func (m *ProjectQuerier) Create(ctx context.Context, repo *database.Repository, p projects.ProjectCopy) error {
	args := m.Called(ctx, repo, p)
	return args.Error(0)
}
func (m *ProjectQuerier) Retrieve(ctx context.Context, repo database.Storer, pid string) (projects.ProjectCopy, error) {
	args := m.Called(ctx, repo, pid)
	return args.Get(0).(projects.ProjectCopy), args.Error(1)

}
func (m *ProjectQuerier) Update(ctx context.Context, repo database.Storer, pid string, update projects.UpdateProjectCopy) error {
	args := m.Called(ctx, repo, pid, update)
	return args.Error(0)

}
func (m *ProjectQuerier) Delete(ctx context.Context, repo database.Storer, pid string) error {
	args := m.Called(ctx, repo, pid)
	return args.Error(0)
}

type InviteQuerier struct {
	mock.Mock
}

func (m *InviteQuerier) Create(ctx context.Context, repo database.Storer, ni invites.NewInvite, now time.Time) (invites.Invite, error) {
	args := m.Called(ctx, repo, ni, now)
	return args.Get(0).(invites.Invite), args.Error(1)
}
func (m *InviteQuerier) RetrieveInvite(ctx context.Context, repo database.Storer, uid string, iid string) (invites.Invite, error) {
	args := m.Called(ctx, repo, uid, iid)
	return args.Get(0).(invites.Invite), args.Error(1)
}
func (m *InviteQuerier) RetrieveInvites(ctx context.Context, repo database.Storer, uid string) ([]invites.Invite, error) {
	args := m.Called(ctx, repo, uid)
	return args.Get(0).([]invites.Invite), args.Error(1)
}
func (m *InviteQuerier) Update(ctx context.Context, repo database.Storer, update invites.UpdateInvite, uid, iid string, now time.Time) (invites.Invite, error) {
	args := m.Called(ctx, repo, update, uid, iid, now)
	return args.Get(0).(invites.Invite), args.Error(1)
}

type MembershipQuerier struct {
	mock.Mock
}

func (m *MembershipQuerier) Create(ctx context.Context, repo database.Storer, nm memberships.NewMembership, now time.Time) (memberships.Membership, error) {
	args := m.Called(ctx, repo, nm, now)
	return args.Get(0).(memberships.Membership), args.Error(1)
}
func (m *MembershipQuerier) RetrieveMemberships(ctx context.Context, repo database.Storer, uid, tid string) ([]memberships.MembershipEnhanced, error) {
	args := m.Called(ctx, repo, uid, tid)
	return args.Get(0).([]memberships.MembershipEnhanced), args.Error(1)
}
func (m *MembershipQuerier) RetrieveMembership(ctx context.Context, repo database.Storer, uid, tid string) (memberships.Membership, error) {
	args := m.Called(ctx, repo, uid, tid)
	return args.Get(0).(memberships.Membership), args.Error(1)
}
func (m *MembershipQuerier) Update(ctx context.Context, repo database.Storer, tid string, update memberships.UpdateMembership, uid string, now time.Time) error {
	args := m.Called(ctx, repo, tid, update, uid, now)
	return args.Error(0)
}
func (m *MembershipQuerier) Delete(ctx context.Context, repo database.Storer, tid, uid string) (string, error) {
	args := m.Called(ctx, repo, tid, uid)
	return args.String(0), args.Error(1)
}
