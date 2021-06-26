package handlers

import (
	"context"
	"github.com/devpies/devpie-client-core/users/domain/memberships"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/stretchr/testify/mock"
	"time"
)

type membershipQueryMock struct {
	mock.Mock
}

func (m *membershipQueryMock) Create(ctx context.Context, repo database.Storer, nm memberships.NewMembership, now time.Time) (memberships.Membership, error) {
	args := m.Called(ctx, repo, nm, now)
	return args.Get(0).(memberships.Membership), args.Error(1)
}
func (m *membershipQueryMock) RetrieveMemberships(ctx context.Context, repo database.Storer, uid, tid string) ([]memberships.MembershipEnhanced, error) {
	args := m.Called(ctx, repo, uid, tid)
	return args.Get(0).([]memberships.MembershipEnhanced), args.Error(1)
}
func (m *membershipQueryMock) RetrieveMembership(ctx context.Context, repo database.Storer, uid, tid string) (memberships.Membership, error) {
	args := m.Called(ctx, repo, uid, tid)
	return args.Get(0).(memberships.Membership), args.Error(1)
}
func (m *membershipQueryMock) Update(ctx context.Context, repo database.Storer, tid string, update memberships.UpdateMembership, uid string, now time.Time) error {
	args := m.Called(ctx, repo, tid, update, uid, now)
	return args.Error(0)
}
func (m *membershipQueryMock) Delete(ctx context.Context, repo database.Storer, tid, uid string) (string, error) {
	args := m.Called(ctx, repo, tid, uid)
	return args.String(0), args.Error(1)
}
