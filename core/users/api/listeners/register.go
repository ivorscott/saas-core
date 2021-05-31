package listeners

import (
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/devpies/devpie-client-events/go/events"
	"github.com/nats-io/stan.go"
	"log"
	"time"
)

type Listeners struct {
	log  *log.Logger
	repo *database.Repository
	dur  time.Duration
}

func NewListeners(log *log.Logger, repo *database.Repository) *Listeners {
	dur, err := time.ParseDuration("5s")
	if err != nil {
		log.Printf("warning: parse duration error: %v", err)
	}
	return &Listeners{log, repo, dur}
}

func (l *Listeners) RegisterAll(nats *events.Client, queueGrp string) {
	nats.Listen(string(events.EventsProjectCreated), queueGrp, l.handleProjectCreated, stan.DeliverAllAvailable(),
		stan.SetManualAckMode(), stan.AckWait(l.dur), stan.DurableName(queueGrp))
	nats.Listen(string(events.EventsProjectUpdated), queueGrp, l.handleProjectUpdated, stan.DeliverAllAvailable(),
		stan.SetManualAckMode(), stan.AckWait(l.dur), stan.DurableName(queueGrp))
	nats.Listen(string(events.EventsProjectDeleted), queueGrp, l.handleProjectDeleted, stan.DeliverAllAvailable(),
		stan.SetManualAckMode(), stan.AckWait(l.dur), stan.DurableName(queueGrp))
}
