package listeners

import (
	"context"
	"github.com/devpies/devpie-client-core/projects/internal/memberships"
	"github.com/devpies/devpie-client-events/go/events"
	"github.com/nats-io/stan.go"
	"time"
)

func (l *Listeners) handleMembershipCreated(m *stan.Msg) {
	msg, err := events.UnmarshalMembershipCreatedEvent(m.Data)
	if err != nil {
		l.log.Printf("warning: failed to unmarshal Command \n %v", err)
	}

	event := msg.Data

	layout := "2006-01-02 15:04:05.999999999 -0700 MST"

	ut, err := time.Parse(layout, event.UpdatedAt)
	if err != nil {
		l.log.Printf("failed to parse time")
	}
	ct, err := time.Parse(layout, event.CreatedAt)
	if err != nil {
		l.log.Printf("failed to parse time")
	}

	mem := memberships.MembershipCopy{
		ID: event.MembershipID,
		UserID: event.UserID,
		TeamID: event.TeamID,
		Role: event.Role,
		UpdatedAt: ut,
		CreatedAt: ct,
	}

	if err := memberships.Create(context.Background(), l.repo, mem); err != nil {
		l.log.Printf("failed to create membership: %s \n %v", event.MembershipID, err)
	}

	if err := m.Ack(); err != nil {
		l.log.Printf("failed to Acknowledge message \n %v", err)
	}
}

func (l *Listeners) handleMembershipUpdated(m *stan.Msg) {
	msg, err := events.UnmarshalMembershipUpdatedEvent(m.Data)
	if err != nil {
		l.log.Printf("warning: failed to unmarshal Command \n %v", err)
	}

	event := msg.Data

	layout := "2006-01-02 15:04:05.999999999 -0700 MST"

	ut, err := time.Parse(layout, event.UpdatedAt)
	if err != nil {
		l.log.Printf("failed to parse time")
	}

	mem := memberships.UpdateMembershipCopy{
		Role: event.Role,
		UpdatedAt: ut,
	}
	if err := memberships.Update(context.Background(), l.repo, event.MembershipID, mem); err != nil {
		l.log.Printf("failed to create membership: %s \n %v", event.MembershipID, err)
	}

	if err := m.Ack(); err != nil {
		l.log.Printf("failed to Acknowledge message \n %v", err)
	}
}


func (l *Listeners) handleMembershipDeleted(m *stan.Msg) {
	msg, err := events.UnmarshalMembershipDeletedEvent(m.Data)
	if err != nil {
		l.log.Printf("warning: failed to unmarshal Command \n %v", err)
	}

	event := msg.Data

	if err := memberships.Delete(context.Background(), l.repo, event.MembershipID); err != nil {
		l.log.Printf("failed to delete membership: %s \n %v", event.MembershipID, err)
	}

	if err := m.Ack(); err != nil {
		l.log.Printf("failed to Acknowledge message \n %v", err)
	}
}