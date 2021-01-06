package main

import (
	"fmt"
	"log"
	"time"
	"github.com/google/uuid"

	"github.com/ivorscott/devpie-client-events/go/events"
	"github.com/nats-io/stan.go"
)

type Handlers struct {
	Repo *Repository
	Client *events.Client
}

func NewHandlers(repo *Repository, client *events.Client) *Handlers {
	return &Handlers{Repo: repo, Client: client}
}

func (h Handlers) handleAddUser(m *stan.Msg) {
	var nu NewUser

	msg, err := events.UnmarshalAddUserCommand(m.Data)
	if err != nil {
		log.Fatal(err)
	}

	u := msg.Data

	if _, err := RetrieveMeByAuth0ID(h.Repo, u.Auth0ID); err != nil {
		nu = NewUser{
			ID:            u.ID,
			Auth0ID:       u.Auth0ID,
			Email:         u.Email,
			EmailVerified: u.EmailVerified,
			FirstName:     &u.FirstName,
			LastName:      &u.LastName,
			Picture:       &u.Picture,
			Locale:        &u.Locale,
		}

		_, err = CreateUser(h.Repo, nu, time.Now())
		if err != nil {
			log.Fatal(err)
		}

		err = m.Ack()
		if err != nil {
			log.Fatal(err)
		}

		id, err := uuid.NewRandom()
		if err != nil {
			log.Fatal(err)
		}

		ne := events.UserAddedEvent{
			ID: id.String(),
			Type: events.TypeUserAdded,
			Data: events.UserAddedEventData{
				ID:            u.ID,
				Auth0ID:       u.Auth0ID,
				Email:         u.Email,
				EmailVerified: u.EmailVerified,
				FirstName:     u.FirstName,
				LastName:      u.LastName,
				Picture:       u.Picture,
				Locale:        u.Locale},
			Metadata: events.Metadata{
				TraceID: msg.Metadata.TraceID,
				UserID: msg.Metadata.UserID,
			},
		}

		bytes, err := ne.Marshal()
		if err !=nil {
			log.Fatal(err)
		}

		cat := events.Identity
		sn := fmt.Sprintf("%s.%s", cat, nu.ID)
		h.Client.Publish(sn, bytes)
	}
}