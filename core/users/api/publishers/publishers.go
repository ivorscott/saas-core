package publishers

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"

	"github.com/devpies/devpie-client-core/users/domain/memberships"
	"github.com/devpies/devpie-client-events/go/events"
)

// Publisher describes the behavior required for publishing events
type Publisher interface {
	ProjectUpdated(nats *events.Client, tid *string, pid, uid string) error
	MembershipCreated(nats *events.Client, m memberships.Membership, uid string) error
	MembershipCreatedForProject(nats *events.Client, m memberships.Membership, pid, uid string) error
	MembershipDeleted(nats *events.Client, mid, uid string) error
}

// Publishers defines handlers that trigger events
type Publishers struct{}

// ProjectUpdated publishes a ProjectUpdatedEvent
func (p *Publishers) ProjectUpdated(nats *events.Client, tid *string, pid, uid string) error {
	ue := events.ProjectUpdatedEvent{
		ID:   uuid.New().String(),
		Type: events.TypeProjectUpdated,
		Data: events.ProjectUpdatedEventData{
			TeamID:    tid,
			ProjectID: pid,
			UpdatedAt: time.Now().UTC().String(),
		},
		Metadata: events.Metadata{
			TraceID: uuid.New().String(),
			UserID:  uid,
		},
	}

	bytes, err := json.Marshal(ue)
	if err != nil {
		return err
	}

	nats.Publish(string(events.EventsProjectUpdated), bytes)

	return nil
}

// MembershipCreated publishes a MembershipCreatedEvent
func (p *Publishers) MembershipCreated(nats *events.Client, m memberships.Membership, uid string) error {
	e := events.MembershipCreatedEvent{
		ID:   uuid.New().String(),
		Type: events.TypeMembershipCreated,
		Data: events.MembershipCreatedEventData{
			MembershipID: m.ID,
			TeamID:       m.TeamID,
			Role:         m.Role,
			UserID:       m.UserID,
			UpdatedAt:    m.UpdatedAt.String(),
			CreatedAt:    m.CreatedAt.String(),
		},
		Metadata: events.Metadata{
			TraceID: uuid.New().String(),
			UserID:  uid,
		},
	}

	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}

	nats.Publish(string(events.EventsMembershipCreated), bytes)

	return nil
}

// MembershipCreatedForProject publishes a MembershipCreatedForProjectEvent
func (p *Publishers) MembershipCreatedForProject(nats *events.Client, m memberships.Membership, pid, uid string) error {
	e := events.MembershipCreatedForProjectEvent{
		ID:   uuid.New().String(),
		Type: events.TypeMembershipCreatedForProject,
		Data: events.MembershipCreatedForProjectEventData{
			MembershipID: m.ID,
			TeamID:       m.TeamID,
			Role:         m.Role,
			UserID:       m.UserID,
			ProjectID:    pid,
			UpdatedAt:    m.UpdatedAt.String(),
			CreatedAt:    m.CreatedAt.String(),
		},
		Metadata: events.Metadata{
			TraceID: uuid.New().String(),
			UserID:  uid,
		},
	}

	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}

	nats.Publish(string(events.EventsMembershipCreatedForProject), bytes) // mock

	return nil
}

// MembershipDeleted publishes a MembershipDeletedEvent
func (p *Publishers) MembershipDeleted(nats *events.Client, mid, uid string) error {
	me := events.MembershipDeletedEvent{
		ID:   uuid.New().String(),
		Type: events.TypeMembershipDeleted,
		Data: events.MembershipDeletedEventData{
			MembershipID: mid,
		},
		Metadata: events.Metadata{
			TraceID: uuid.New().String(),
			UserID:  uid,
		},
	}

	bytes, err := json.Marshal(me)
	if err != nil {
		return err
	}

	nats.Publish(string(events.EventsMembershipDeleted), bytes)

	return nil
}
