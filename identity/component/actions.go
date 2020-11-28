package main

import (
	"github.com/ivorscott/devpie-client-events/go/events"
	"github.com/nats-io/stan.go"
	"log"
	"time"
)

func handleAddUserCommand(m *stan.Msg) {
	var nu NewUser

	u, err := events.UnmarshalAddUserCommand(m.Data)
	if err != nil {
		log.Fatal(err)
	}

	ud := u.Data

	nu = NewUser{
		Auth0ID: ud.Auth0ID,
		Email: ud.Email,
		EmailVerified: ud.EmailVerified,
		FirstName: &ud.FirstName,
		LastName: &ud.LastName,
		Picture: &ud.Picture,
		Locale: &ud.Locale,
	}

	_, err = CreateUser(repo, nu, nu.Auth0ID, time.Now())
	if err !=nil {
		log.Fatal(err)
	}

	m.Ack()
}