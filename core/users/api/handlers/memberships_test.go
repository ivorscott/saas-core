package handlers

import (
	"time"

	"github.com/devpies/devpie-client-core/users/domain/memberships"
	"github.com/devpies/devpie-client-core/users/domain/teams"
)

func newMembership(t teams.Team) memberships.NewMembership {
	return memberships.NewMembership{
		TeamID: t.ID,
		UserID: t.UserID,
		Role:   "administrator",
	}
}

func membership(nm memberships.NewMembership) memberships.Membership {
	return memberships.Membership{
		ID:        "085cb8a0-b221-4a6d-95be-592eb68d5670",
		TeamID:    nm.TeamID,
		UserID:    nm.UserID,
		Role:      nm.Role,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
}
