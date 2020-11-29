package main

import (
	"fmt"
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

	fmt.Printf("Command user data: %+v", u)

	nu = NewUser{
		ID: ud.ID,
		Auth0ID: ud.Auth0ID,
		Email: ud.Email,
		EmailVerified: ud.EmailVerified,
		FirstName: &ud.FirstName,
		LastName: &ud.LastName,
		Picture: &ud.Picture,
		Locale: &ud.Locale,
	}

	fmt.Printf("New user: %+v", nu)

	_, err = CreateUser(repo, nu, time.Now())
	if err !=nil {
		log.Fatal(err)
	}

	m.Ack()

	bytes, err := u.Marshal()
	if err !=nil {
		log.Fatal(err)
	}

	c.Publish(string(events.EventsUserAdded), bytes)
}