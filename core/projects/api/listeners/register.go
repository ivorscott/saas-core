package listeners

import (
	"log"
	"time"

	"github.com/nats-io/stan.go"

	"github.com/devpies/devpie-client-core/projects/platform/database"
	"github.com/devpies/devpie-client-events/go/events"
)

type Listeners struct {
	log  *log.Logger
	repo *database.Repository
	dur  time.Duration
}

func NewListeners(log *log.Logger, repo *database.Repository) *Listeners {
	dur, err := time.ParseDuration("5s")
	if err != nil {
		log.Printf("parse duration error: %v", err)
	}
	return &Listeners{log, repo, dur}
}

func (l *Listeners) RegisterAll(nats *events.Client, queueGrp string) {
	nats.Listen(string(events.EventsMembershipCreated), queueGrp, l.handleMembershipCreated, stan.DeliverAllAvailable(),
		stan.SetManualAckMode(), stan.AckWait(l.dur), stan.DurableName(queueGrp))

	nats.Listen(string(events.EventsMembershipUpdated), queueGrp, l.handleMembershipUpdated, stan.DeliverAllAvailable(),
		stan.SetManualAckMode(), stan.AckWait(l.dur), stan.DurableName(queueGrp))

	nats.Listen(string(events.EventsMembershipDeleted), queueGrp, l.handleMembershipDeleted, stan.DeliverAllAvailable(),
		stan.SetManualAckMode(), stan.AckWait(l.dur), stan.DurableName(queueGrp))

	nats.Listen(string(events.EventsMembershipCreatedForProject), queueGrp, l.handleMembershipCreatedForProject, stan.DeliverAllAvailable(),
		stan.SetManualAckMode(), stan.AckWait(l.dur), stan.DurableName(queueGrp))

	nats.Listen(string(events.EventsProjectUpdated), queueGrp, l.handleProjectUpdated, stan.DeliverAllAvailable(),
		stan.SetManualAckMode(), stan.AckWait(l.dur), stan.DurableName(queueGrp))
}
