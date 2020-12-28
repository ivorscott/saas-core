package main

import (
	"fmt"
	"log"
	"time"

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

	fmt.Printf("Sequence: %d, Subject: %s, Message: %+v,", m.Sequence, m.Subject, msg)

	if _, err := RetrieveMeByAuth0ID(h.Repo, u.Auth0ID); err != nil {
		nu = NewUser{
			ID: u.ID,
			Auth0ID: u.Auth0ID,
			Email: u.Email,
			EmailVerified: u.EmailVerified,
			FirstName: &u.FirstName,
			LastName: &u.LastName,
			Picture: &u.Picture,
			Locale: &u.Locale,
		}

		fmt.Printf("New user: %+v", nu)
		_, err = CreateUser(h.Repo, nu, time.Now());
		if err !=nil {
			log.Fatal(err)
		}

		err = m.Ack()
		if err != nil {
			log.Fatal(err)
		}

		bytes, err := msg.Marshal()
		if err !=nil {
			log.Fatal(err)
		}

		cat := events.Identity
		sn := fmt.Sprintf("%s.%s", cat, nu.ID)
		h.Client.Publish(sn, bytes)
	}
}