package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/devpies/devpie-client-core/users/domain/invites"
	"github.com/devpies/devpie-client-core/users/domain/memberships"
	"github.com/devpies/devpie-client-core/users/domain/projects"
	"github.com/devpies/devpie-client-core/users/domain/teams"
	"github.com/devpies/devpie-client-core/users/domain/users"
	"github.com/devpies/devpie-client-core/users/platform/auth0"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/devpies/devpie-client-core/users/platform/web"
	"github.com/devpies/devpie-client-events/go/events"
)

type Team struct {
	repo        *database.Repository
	log         *log.Logger
	auth0       *auth0.Auth0
	nats        *events.Client
	origins     string
	sendgridKey string
}

func (t *Team) Create(w http.ResponseWriter, r *http.Request) error {
	var nt teams.NewTeam
	var role memberships.Role = memberships.Administrator

	uid := t.auth0.GetUser(r, users.RetrieveMeByAuthID)

	if err := web.Decode(r, &nt); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	if _, err := projects.Retrieve(r.Context(), t.repo, nt.ProjectID); err != nil {
		switch err {
		case projects.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case projects.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "creating team for project %q", nt.ProjectID)
		}
	}


	tm, err := teams.Create(r.Context(), t.repo, nt, uid, time.Now())
	if err != nil {
		return err
	}

	nm := memberships.NewMembership{
		UserID: uid,
		TeamID: tm.ID,
		Role:   role.String(),
	}

	m, err := memberships.Create(r.Context(), t.repo, nm, time.Now())
	if err != nil {
		return err
	}

	if nt.ProjectID == "" {
		// TODO: remove, no longer creating team outside of project context
		e := events.MembershipCreatedEvent{
			ID:   uuid.New().String(),
			Type: events.TypeMembershipCreated,
			Data: events.MembershipCreatedEventData{
				MembershipID: m.ID,
				TeamID:       m.TeamID,
				Role:         m.Role,
				UserID:       m.UserID,
				UpdatedAt:    m.UpdatedAt.String(),
				CreatedAt:    m.CreatedAt.String(),
			},
			Metadata: events.Metadata{
				TraceID: uuid.New().String(),
				UserID:  uid,
			},
		}

		bytes, err := json.Marshal(e)
		if err != nil {
			return err
		}

		if t.nats != nil {
			t.nats.Publish(string(events.EventsMembershipCreated), bytes)
		}
	} else {
		e := events.MembershipCreatedForProjectEvent{
			ID:   uuid.New().String(),
			Type: events.TypeMembershipCreatedForProject,
			Data: events.MembershipCreatedForProjectEventData{
				MembershipID: m.ID,
				TeamID:       m.TeamID,
				Role:         m.Role,
				UserID:       m.UserID,
				ProjectID:    nt.ProjectID,
				UpdatedAt:    m.UpdatedAt.String(),
				CreatedAt:    m.CreatedAt.String(),
			},
			Metadata: events.Metadata{
				TraceID: uuid.New().String(),
				UserID:  uid,
			},
		}

		up := projects.UpdateProjectCopy{
			TeamID: &tm.ID,
		}

		if err := projects.Update(r.Context(), t.repo, nt.ProjectID, up); err != nil {
			return err
		}

		bytes, err := json.Marshal(e)
		if err != nil {
			return err
		}

		if t.nats != nil {
			t.nats.Publish(string(events.EventsMembershipCreatedForProject), bytes)
		}
	}

	return web.Respond(r.Context(), w, tm, http.StatusCreated)
}

func (t *Team) AssignExisting(w http.ResponseWriter, r *http.Request) error {
	tid := chi.URLParam(r, "tid")
	pid := chi.URLParam(r, "pid")
	uid := t.auth0.GetUser(r, users.RetrieveMeByAuthID)

	tm, err := teams.Retrieve(r.Context(), t.repo, tid)
	if err != nil {
		return web.NewRequestError(err, http.StatusNotFound)
	}

	var up = projects.UpdateProjectCopy{
		TeamID:    &tm.ID,
		UpdatedAt: time.Now().UTC(),
	}

	err = projects.Update(r.Context(), t.repo, pid, up)
	if err != nil {
		switch err {
		case projects.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case projects.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "updating project %q", pid)
		}
	}

	ue := events.ProjectUpdatedEvent{
		ID:   uuid.New().String(),
		Type: events.TypeProjectUpdated,
		Data: events.ProjectUpdatedEventData{
			TeamID:    &tm.ID,
			ProjectID: pid,
			UpdatedAt: time.Now().UTC().String(),
		},
		Metadata: events.Metadata{
			TraceID: uuid.New().String(),
			UserID:  uid,
		},
	}

	bytes, err := json.Marshal(ue)
	if err != nil {
		return err
	}

	t.nats.Publish(string(events.EventsProjectUpdated), bytes)

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

func (t *Team) LeaveTeam(w http.ResponseWriter, r *http.Request) error {
	tid := chi.URLParam(r, "tid")

	uid := t.auth0.GetUser(r, users.RetrieveMeByAuthID)

	// if user is the administrator
	// and is the last to leave
	// delete the team

	// if the user is the administrator
	// and is not the last to leave
	// ownership must be passed to another member

	mid, err := memberships.Delete(r.Context(), t.repo, tid, uid)
	if err != nil {
		switch err {
		case teams.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case teams.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for team %q", tid)
		}
	}

	me := events.MembershipDeletedEvent{
		ID:   uuid.New().String(),
		Type: events.TypeMembershipDeleted,
		Data: events.MembershipDeletedEventData{
			MembershipID: mid,
		},
		Metadata: events.Metadata{
			TraceID: uuid.New().String(),
			UserID:  uid,
		},
	}

	bytes, err := json.Marshal(me)
	if err != nil {
		return err
	}

	t.nats.Publish(string(events.EventsMembershipDeleted), bytes)

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

func (t *Team) Retrieve(w http.ResponseWriter, r *http.Request) error {
	tid := chi.URLParam(r, "tid")

	tm, err := teams.Retrieve(r.Context(), t.repo, tid)
	if err != nil {
		switch err {
		case teams.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case teams.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for team %q", tid)
		}
	}

	return web.Respond(r.Context(), w, tm, http.StatusOK)
}

func (t *Team) List(w http.ResponseWriter, r *http.Request) error {
	uid := t.auth0.GetUser(r, users.RetrieveMeByAuthID)

	tms, err := teams.List(r.Context(), t.repo, uid)
	if err != nil {
		switch err {
		case teams.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case teams.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for user's teams")
		}
	}

	return web.Respond(r.Context(), w, tms, http.StatusOK)
}

func (t *Team) CreateInvite(w http.ResponseWriter, r *http.Request) error {
	var token auth0.Token
	var list invites.NewList

	link := strings.Split(t.origins, ",")[0]

	if err := web.Decode(r, &list); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	// Get valid token
	token, err := t.auth0.RetrieveToken()
	if err == auth0.ErrNotFound || t.auth0.IsExpired(token) {
		nt, err := t.auth0.NewManagementToken()
		if err != nil {
			return err
		}
		// clean table before persisting
		if err := t.auth0.DeleteToken(); err != nil {
			return err
		}
		// persist management api token
		tk, err := t.auth0.PersistToken(nt, time.Now())
		if err != nil {
			return err
		}
		token = tk
	}

	tid := chi.URLParam(r, "tid")
	if err != nil {
		return err
	}

	for _, email := range list.Emails {

		ni := invites.NewInvite{
			TeamID: tid,
		}

		// existing user
		u, err := users.RetrieveByEmail(t.repo, email)
		if err != nil {
			// new user
			au, err := t.auth0.CreateUser(token, email)
			if err != nil {
				return err
			}

			nu := users.NewUser{
				Auth0ID:       au.Auth0ID,
				Email:         au.Email,
				EmailVerified: au.EmailVerified,
				FirstName:     au.FirstName,
				Picture:       au.Picture,
			}

			user, err := users.Create(r.Context(), t.repo, nu, au.Auth0ID, time.Now())
			if err != nil {
				return err
			}

			ni.UserID = user.ID

			if err := t.auth0.UpdateUserAppMetaData(token, au.Auth0ID, user.ID); err != nil {
				return err
			}

			link, err = t.auth0.ChangePasswordTicket(token, au, link)
			if err != nil {
				return err
			}

		} else {
			ni.UserID = u.ID
		}

		if err := t.SendMail(email, link); err != nil {
			return err
		}

		_, err = invites.Create(r.Context(), t.repo, ni, time.Now())
		if err != nil {
			return err
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusCreated)
}

func (t *Team) RetrieveInvites(w http.ResponseWriter, r *http.Request) error {
	uid := t.auth0.GetUser(r, users.RetrieveMeByAuthID)

	is, err := invites.RetrieveInvites(r.Context(), t.repo, uid)
	if err != nil {
		switch err {
		case teams.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case teams.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "searching team invites for %q", uid)
		}
	}

	var res []invites.InviteEnhanced
	for _, invite := range is {
		team, err := teams.Retrieve(r.Context(), t.repo, invite.TeamID)
		if err != nil {
			return err
		}
		ie := invites.InviteEnhanced{
			ID:         invite.ID,
			UserID:     invite.UserID,
			TeamID:     invite.TeamID,
			TeamName:   team.Name,
			Read:       invite.Read,
			Accepted:   invite.Accepted,
			Expiration: invite.Expiration,
			UpdatedAt:  invite.UpdatedAt,
			CreatedAt:  invite.CreatedAt,
		}
		res = append(res, ie)
	}

	return web.Respond(r.Context(), w, res, http.StatusOK)
}

func (t *Team) UpdateInvite(w http.ResponseWriter, r *http.Request) error {
	var update invites.UpdateInvite
	var role memberships.Role = memberships.Editor

	uid := t.auth0.GetUser(r, users.RetrieveMeByAuthID)
	tid := chi.URLParam(r, "tid")
	iid := chi.URLParam(r, "iid")

	if err := web.Decode(r, &update); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	iv, err := invites.Update(r.Context(), t.repo, update, uid, iid, time.Now())
	if err != nil {
		return err
	}

	if update.Accepted {
		nm := memberships.NewMembership{
			UserID: uid,
			TeamID: tid,
			Role:   role.String(),
		}
		m, err := memberships.Create(r.Context(), t.repo, nm, time.Now())
		if err != nil {
			return err
		}

		e := events.MembershipCreatedEvent{
			ID:   uuid.New().String(),
			Type: events.TypeMembershipCreated,
			Data: events.MembershipCreatedEventData{
				MembershipID: m.ID,
				TeamID:       m.TeamID,
				Role:         m.Role,
				UserID:       m.UserID,
				UpdatedAt:    m.UpdatedAt.String(),
				CreatedAt:    m.CreatedAt.String(),
			},
			Metadata: events.Metadata{
				TraceID: uuid.New().String(),
				UserID:  uid,
			},
		}

		bytes, err := json.Marshal(e)
		if err != nil {
			return err
		}

		t.nats.Publish(string(events.EventsMembershipCreated), bytes)
	}

	return web.Respond(r.Context(), w, iv, http.StatusOK)
}

func (t *Team) SendMail(email, link string) error {
	from := mail.NewEmail("DevPie", "people@devpie.io")
	subject := "You've been invited to a Team on DevPie!"
	to := mail.NewEmail("Invitee", email)

	html := ""
	html += "<strong>Join Devpie</strong>"
	html += "<br/>"
	html += "<p>To accept your invitation, <a href=\"%s\">create an account</a>.</p>"
	htmlContent := fmt.Sprintf(html, link)

	plainTextContent := fmt.Sprintf("You've been invited to a Team on DevPie! %s ", link)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(t.sendgridKey)

	response, err := client.Send(message)
	if err != nil {
		return err
	} else {
		t.log.Println(response.StatusCode)
		t.log.Println(response.Body)
		t.log.Println(response.Headers)
	}
	return nil
}
