package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/devpies/devpie-client-events/go/events"
	"github.com/nats-io/stan.go"
)

type Handlers struct {
	client *events.Client
	log *log.Logger
}

func NewHandlers(client *events.Client, log *log.Logger) *Handlers {
	return &Handlers{client: client, log: log}
}

func (h Handlers) handleAddUser(m *stan.Msg) {

	cmd, err := events.UnmarshalAddUserCommand(m.Data)
	if err != nil {
		log.Fatal(err)
	}

	u := cmd.Data

	log.Printf("message data: %v", u)

	// search for user added event 
	// if any previous entry matches last message user was already added 

	eid, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}

	ne := events.UserAddedEvent{
		ID: eid.String(),
		Type: events.TypeUserAdded,
		Data: events.UserAddedEventData{
			ID:            u.ID,
			Auth0ID:       u.Auth0ID,
			Email:         u.Email,
			EmailVerified: u.EmailVerified,
			FirstName:     u.FirstName,
			LastName:      u.LastName,
			Picture:       u.Picture,
			Locale:        u.Locale,
		},
		Metadata: events.Metadata{
			TraceID: cmd.Metadata.TraceID,
			UserID: cmd.Metadata.UserID,
		},
	}

	bytes, err := ne.Marshal()
	if err !=nil {
		log.Fatal(err)
	}

	cat := events.Identity
	sn := fmt.Sprintf("%s.%s", cat, u.ID)
	h.client.Publish(sn, bytes)

	err = m.Ack()
	if err != nil {
		log.Fatal(err)
	}
}