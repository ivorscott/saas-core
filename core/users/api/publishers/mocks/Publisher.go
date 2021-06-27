package mocks


import (
	"github.com/devpies/devpie-client-core/users/domain/memberships"
	"github.com/devpies/devpie-client-events/go/events"
	"github.com/stretchr/testify/mock"
)

type Publisher struct {
	mock.Mock
}

func (m *Publisher) ProjectUpdated(nats *events.Client, tid *string, pid, uid string) error {
	args := m.Called(nats, tid, pid, uid)
	return args.Error(0)
}
func (m *Publisher) MembershipCreated(nats *events.Client, mem memberships.Membership, uid string) error {
	args := m.Called(nats, mem, uid)
	return args.Error(0)
}
func (m *Publisher) MembershipCreatedForProject(nats *events.Client, mem memberships.Membership, pid , uid string) error {
	args := m.Called(nats, mem, pid, uid)
	return args.Error(0)
}
func (m *Publisher) MembershipDeleted(nats *events.Client, mid, uid string) error {
	args := m.Called(nats, mid, uid)
	return args.Error(0)
}