package listeners

import (
	"context"

	"github.com/devpies/devpie-client-core/projects/domain/memberships"
	"github.com/devpies/devpie-client-core/projects/domain/projects"
	"github.com/devpies/devpie-client-events/go/events"
	"github.com/nats-io/stan.go"
)

func (l *Listeners) handleMembershipCreated(m *stan.Msg) {
	msg, err := events.UnmarshalMembershipCreatedEvent(m.Data)
	if err != nil {
		l.log.Printf("warning: failed to unmarshal Command \n %v", err)
	}

	event := msg.Data

	updatedtime, err := events.ParseTime(event.UpdatedAt)
	if err != nil {
		l.log.Printf("failed to parse time")
	}
	createdtime, err := events.ParseTime(event.CreatedAt)
	if err != nil {
		l.log.Printf("failed to parse time")
	}

	mem := memberships.MembershipCopy{
		ID: event.MembershipID,
		UserID: event.UserID,
		TeamID: event.TeamID,
		Role: event.Role,
		UpdatedAt: updatedtime,
		CreatedAt: createdtime,
	}

	if err := memberships.Create(context.Background(), l.repo, mem); err != nil {
		l.log.Printf("failed to create membership: %s \n %v", event.MembershipID, err)
	}

	if err := m.Ack(); err != nil {
		l.log.Printf("failed to Acknowledge message \n %v", err)
	}
}

func (l *Listeners) handleMembershipCreatedForProject(m *stan.Msg) {
	msg, err := events.UnmarshalMembershipCreatedForProjectEvent(m.Data)
	if err != nil {
		l.log.Printf("warning: failed to unmarshal Command \n %v", err)
	}

	event := msg.Data

	updatedtime, err := events.ParseTime(event.UpdatedAt)
	if err != nil {
		l.log.Printf("failed to parse time")
	}
	createdtime, err := events.ParseTime(event.CreatedAt)
	if err != nil {
		l.log.Printf("failed to parse time")
	}

	mem := memberships.MembershipCopy{
		ID: event.MembershipID,
		UserID: event.UserID,
		TeamID: event.TeamID,
		Role: event.Role,
		UpdatedAt: updatedtime,
		CreatedAt: createdtime,
	}

	if err := memberships.Create(context.Background(), l.repo, mem); err != nil {
		l.log.Printf("failed to create membership: %s \n %v", event.MembershipID, err)
	}

	update := projects.UpdateProject{
		TeamID: &event.TeamID,
	}

	if _, err := projects.Update(context.Background(), l.repo, event.ProjectID, event.UserID, update, updatedtime); err != nil {
		l.log.Printf("failed to update projects: %s \n %v", event.ProjectID, err)
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

	updatedtime, err := events.ParseTime(event.UpdatedAt)
	if err != nil {
		l.log.Printf("failed to parse time")
	}

	mem := memberships.UpdateMembershipCopy{
		Role: event.Role,
		UpdatedAt: updatedtime,
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

func (l *Listeners) handleProjectUpdated(m *stan.Msg) {
	msg, err := events.UnmarshalProjectUpdatedEvent(m.Data)
	if err != nil {
		l.log.Printf("warning: failed to unmarshal Command \n %v", err)
	}

	event := msg.Data
	updatedtime, err := events.ParseTime(event.UpdatedAt)
	if err != nil {
		l.log.Printf("failed to parse time")
	}

	update := projects.UpdateProject{
		Name:        event.Name,
		Description: event.Description,
		Active:      event.Active,
		Public:      event.Public,
		TeamID:      event.TeamID,
		ColumnOrder: event.ColumnOrder,
	}

	if _, err = projects.Update(context.Background(), l.repo, event.ProjectID, msg.Metadata.UserID, update, updatedtime); err != nil {
		l.log.Printf("failed to update project: %s \n %v", event.ProjectID, err)
	}

	err = m.Ack()
	if err != nil {
		l.log.Printf("failed to Acknowledge message \n %v", err)
	}
}