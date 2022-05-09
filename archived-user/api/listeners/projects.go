package listeners

import (
	"context"

	"github.com/devpies/devpie-client-core/users/domain/projects"
	"github.com/devpies/devpie-client-events/go/events"
	"github.com/nats-io/stan.go"
)

// handleProjectCreated listens for a ProjectCreatedEvent
func (l *Listener) handleProjectCreated(m *stan.Msg) {
	msg, err := events.UnmarshalProjectCreatedEvent(m.Data)
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

	update := projects.ProjectCopy{
		ID:          event.ProjectID,
		Name:        event.Name,
		Prefix:      event.Prefix,
		Description: event.Description,
		UserID:      event.UserID,
		TeamID:      event.TeamID,
		Active:      event.Active,
		Public:      event.Public,
		ColumnOrder: event.ColumnOrder,
		UpdatedAt:   updatedtime,
		CreatedAt:   createdtime,
	}

	if err = l.query.project.Create(context.Background(), l.repo, update); err != nil {
		l.log.Printf("failed to update project: %s \n %v", event.ProjectID, err)
	}

	err = m.Ack()
	if err != nil {
		l.log.Printf("failed to Acknowledge message \n %v", err)
	}
}

// handleProjectUpdated listens for a ProjectUpdatedEvent
func (l *Listener) handleProjectUpdated(m *stan.Msg) {
	msg, err := events.UnmarshalProjectUpdatedEvent(m.Data)
	if err != nil {
		l.log.Printf("warning: failed to unmarshal Command \n %v", err)
	}

	event := msg.Data

	updatedtime, err := events.ParseTime(event.UpdatedAt)
	if err != nil {
		l.log.Printf("failed to parse time")
	}

	update := projects.UpdateProjectCopy{
		Name:        event.Name,
		Description: event.Description,
		Active:      event.Active,
		Public:      event.Public,
		TeamID:      event.TeamID,
		ColumnOrder: event.ColumnOrder,
		UpdatedAt:   updatedtime,
	}

	if err = l.query.project.Update(context.Background(), l.repo, event.ProjectID, update); err != nil {
		l.log.Printf("failed to update project: %s \n %v", event.ProjectID, err)
	}

	err = m.Ack()
	if err != nil {
		l.log.Printf("failed to Acknowledge message \n %v", err)
	}
}

// handleProjectDeleted listens for a ProjectDeletedEvent
func (l *Listener) handleProjectDeleted(m *stan.Msg) {
	msg, err := events.UnmarshalProjectDeletedEvent(m.Data)
	if err != nil {
		l.log.Printf("warning: failed to unmarshal Command \n %v", err)
	}

	event := msg.Data

	if err = l.query.project.Delete(context.Background(), l.repo, event.ProjectID); err != nil {
		l.log.Printf("failed to delete project: %s \n %v", event.ProjectID, err)
	}

	err = m.Ack()
	if err != nil {
		l.log.Printf("failed to Acknowledge message \n %v", err)
	}
}
