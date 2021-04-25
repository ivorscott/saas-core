package listeners

import (
	"github.com/devpies/devpie-client-events/go/events"
	"github.com/ivorscott/devpie-client-backend-go/internal/user"
	"github.com/nats-io/stan.go"
	"time"
)

func (l *Listeners) handleUserAdded(m *stan.Msg) {
	// TODO: update team member
	var nu user.NewUser

	msg, err := events.UnmarshalUserAddedEvent(m.Data)
	if err != nil {
		l.log.Printf("warning: failed to unmarshal Command \n %v", err)
	}
	u := msg.Data

	if _, err := user.Retrieve(l.repo, u.Auth0ID); err != nil {
		nu = user.NewUser{
			Auth0ID:       u.Auth0ID,
			Email:         u.Email,
		}
		_, err = user.Create(l.repo, nu, time.Now())
		if err != nil {
			l.log.Printf("warning: failed to copy user \n %v", err)
		}
	}

	err = m.Ack()
	if err != nil {
		l.log.Printf("warning: failed to Acknowledge message \n %v", err)
	}
}
