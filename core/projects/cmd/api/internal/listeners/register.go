package listeners

import (
	"github.com/devpies/devpie-client-events/go/events"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/nats-io/stan.go"
	"log"
	"time"
)

type Listeners struct {
	log *log.Logger
	repo *database.Repository
	dur time.Duration
}

func NewListeners(log *log.Logger, repo *database.Repository) *Listeners {
	dur, err := time.ParseDuration("5s")
	if err != nil {
		log.Printf("warning: parse duration error: %v", err)
	}
	return &Listeners{log, repo, dur}
}

func(l *Listeners) RegisterAll(nats *events.Client, queueGrp string) {
	nats.ListenQ(string(events.EventsUserAdded), queueGrp, l.handleUserAdded, stan.DeliverAllAvailable(),
		stan.SetManualAckMode(), stan.AckWait(l.dur), stan.DurableName(queueGrp))
}
