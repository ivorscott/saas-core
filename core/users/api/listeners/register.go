package listeners

import (
	"context"
	"github.com/nats-io/stan.go"
	"log"
	"time"

	"github.com/devpies/devpie-client-core/users/domain/projects"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/devpies/devpie-client-events/go/events"
)

// Listener defines subscription handlers and their dependencies
type Listener struct {
	log   *log.Logger
	repo  *database.Repository
	dur   time.Duration
	query ListenerQueries
}

// ListenerQueries defines queries required by subscription handlers
type ListenerQueries struct {
	project ProjectQuerier
}

// ProjectQuerier describes behavior required for executing project related queries
type ProjectQuerier interface {
	Create(ctx context.Context, repo *database.Repository, p projects.ProjectCopy) error
	Update(ctx context.Context, repo database.Storer, pid string, update projects.UpdateProjectCopy) error
	Delete(ctx context.Context, repo database.Storer, pid string) error
}

// NewListener creates a new Listener object
func NewListener(log *log.Logger, repo *database.Repository) *Listener {
	dur, err := time.ParseDuration("5s")
	if err != nil {
		log.Printf("warning: parse duration error: %v", err)
	}
	return &Listener{log, repo, dur, ListenerQueries{&projects.Queries{}}}
}

// RegisterAll registers all subscription handlers defined in the implementation body
func (l *Listener) RegisterAll(nats *events.Client, queueGrp string) {
	nats.Listen(string(events.EventsProjectCreated), queueGrp, l.handleProjectCreated, stan.DeliverAllAvailable(),
		stan.SetManualAckMode(), stan.AckWait(l.dur), stan.DurableName(queueGrp))
	nats.Listen(string(events.EventsProjectUpdated), queueGrp, l.handleProjectUpdated, stan.DeliverAllAvailable(),
		stan.SetManualAckMode(), stan.AckWait(l.dur), stan.DurableName(queueGrp))
	nats.Listen(string(events.EventsProjectDeleted), queueGrp, l.handleProjectDeleted, stan.DeliverAllAvailable(),
		stan.SetManualAckMode(), stan.AckWait(l.dur), stan.DurableName(queueGrp))
}
