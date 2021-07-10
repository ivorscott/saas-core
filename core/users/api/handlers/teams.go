package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/devpies/devpie-client-core/users/api/publishers"
	"github.com/devpies/devpie-client-core/users/domain/invites"
	"github.com/devpies/devpie-client-core/users/domain/memberships"
	"github.com/devpies/devpie-client-core/users/domain/projects"
	"github.com/devpies/devpie-client-core/users/domain/teams"
	"github.com/devpies/devpie-client-core/users/domain/users"
	"github.com/devpies/devpie-client-core/users/platform/auth0"
	"github.com/devpies/devpie-client-core/users/platform/database"
	"github.com/devpies/devpie-client-core/users/platform/sendgrid"
	"github.com/devpies/devpie-client-core/users/platform/web"
	"github.com/devpies/devpie-client-events/go/events"
	"github.com/go-chi/chi"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Team defines team handlers and their dependencies
type Team struct {
	repo    database.Storer
	log     *log.Logger
	auth0   auth0.Auther
	nats    *events.Client
	origins string
	sender  sendgrid.Sender
	query   TeamQueries
	publish publishers.Publisher
}

// TeamQueries defines queries required by team handlers
type TeamQueries struct {
	team       teams.TeamQuerier
	project    projects.ProjectQuerier
	membership memberships.MembershipQuerier
	user       users.UserQuerier
	invite     invites.InviteQuerier
}

// Create creates a new team for a project
func (t *Team) Create(w http.ResponseWriter, r *http.Request) error {
	var nt teams.NewTeam
	var role memberships.Role = memberships.Administrator

	if err := web.Decode(r, &nt); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	uid := t.auth0.UserByID(r.Context())

	if _, err := t.query.project.Retrieve(r.Context(), t.repo, nt.ProjectID); err != nil {
		switch err {
		case projects.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case projects.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("failed to retrieve project: %w", err)
		}
	}

	tm, err := t.query.team.Create(r.Context(), t.repo, nt, uid, time.Now())
	if err != nil {
		return err
	}

	nm := memberships.NewMembership{
		UserID: uid,
		TeamID: tm.ID,
		Role:   role.String(),
	}

	m, err := t.query.membership.Create(r.Context(), t.repo, nm, time.Now())
	if err != nil {
		return err
	}

	up := projects.UpdateProjectCopy{
		TeamID: &tm.ID,
	}

	if err = t.query.project.Update(r.Context(), t.repo, nt.ProjectID, up); err != nil {
		return err
	}

	if t.nats != nil {
		err = t.publish.MembershipCreatedForProject(t.nats, m, nt.ProjectID, uid)
		if err != nil {
			return err
		}
	}

	return web.Respond(r.Context(), w, tm, http.StatusCreated)
}

// AssignExistingTeam assigns an existing team to a project
func (t *Team) AssignExistingTeam(w http.ResponseWriter, r *http.Request) error {
	tid := chi.URLParam(r, "tid")
	pid := chi.URLParam(r, "pid")
	uid := t.auth0.UserByID(r.Context())

	tm, err := t.query.team.Retrieve(r.Context(), t.repo, tid)
	if err != nil {
		switch err {
		case teams.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case teams.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("failed to retrieve team: %w", err)
		}
	}

	var up = projects.UpdateProjectCopy{
		TeamID:    &tm.ID,
		UpdatedAt: time.Now().UTC(),
	}

	err = t.query.project.Update(r.Context(), t.repo, pid, up)
	if err != nil {
		switch err {
		case projects.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case projects.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("failed to update project: %w", err)
		}
	}

	if t.nats != nil {
		err = t.publish.ProjectUpdated(t.nats, &tm.ID, pid, uid)
		if err != nil {
			return err
		}
	}
	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

// LeaveTeam destroys a team membership
func (t *Team) LeaveTeam(w http.ResponseWriter, r *http.Request) error {
	tid := chi.URLParam(r, "tid")

	uid := t.auth0.UserByID(r.Context())

	mid, err := t.query.membership.Delete(r.Context(), t.repo, tid, uid)
	if err != nil {
		switch err {
		case memberships.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case memberships.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("failed to delete membership: %w", err)
		}
	}

	if t.nats != nil {
		err = t.publish.MembershipDeleted(t.nats, mid, uid)
		if err != nil {
			return err
		}
	}
	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

// Retrieve returns a team by id
func (t *Team) Retrieve(w http.ResponseWriter, r *http.Request) error {
	tid := chi.URLParam(r, "tid")

	tm, err := t.query.team.Retrieve(r.Context(), t.repo, tid)
	if err != nil {
		switch err {
		case teams.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case teams.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("failed to retrieve team: %w", err)
		}
	}

	return web.Respond(r.Context(), w, tm, http.StatusOK)
}

// List returns all teams associated with the authenticated user
func (t *Team) List(w http.ResponseWriter, r *http.Request) error {
	uid := t.auth0.UserByID(r.Context())

	tms, err := t.query.team.List(r.Context(), t.repo, uid)
	if err != nil {
		switch err {
		case teams.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case teams.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("failed to retrieve teams: %w", err)
		}
	}

	return web.Respond(r.Context(), w, tms, http.StatusOK)
}

// CreateInvite sends new team invitations
func (t *Team) CreateInvite(w http.ResponseWriter, r *http.Request) error {
	var list invites.NewList

	tid := chi.URLParam(r, "tid")
	link := strings.Split(t.origins, ",")[0]

	if err := web.Decode(r, &list); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	token, err := t.auth0.GenerateToken()
	if err != nil {
		return fmt.Errorf("failure during token generation: %w", err)
	}

	for _, email := range list.Emails {
		ni := invites.NewInvite{
			TeamID: tid,
		}
		// when user exists
		u, err := t.query.user.RetrieveByEmail(t.repo, email)
		if err != nil {
			var au auth0.AuthUser

			au, err = t.auth0.CreateUser(token, email)
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

			var us users.User

			us, err = t.query.user.Create(r.Context(), t.repo, nu, time.Now())
			if err != nil {
				return err
			}

			ni.UserID = us.ID

			if err = t.auth0.UpdateUserAppMetaData(token, au.Auth0ID, us.ID); err != nil {
				return err
			}

			link, err = t.auth0.ChangePasswordTicket(token, au, link)
			if err != nil {
				return err
			}

		} else {
			ni.UserID = u.ID
		}

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

		_, err = sendgrid.SendMail(message, t.sender)
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}

		_, err = t.query.invite.Create(r.Context(), t.repo, ni, time.Now())
		if err != nil {
			return err
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusCreated)
}

// RetrieveInvites returns invitations for the authenticated user
func (t *Team) RetrieveInvites(w http.ResponseWriter, r *http.Request) error {
	uid := t.auth0.UserByID(r.Context())

	is, err := t.query.invite.RetrieveInvites(r.Context(), t.repo, uid)
	if err != nil {
		switch err {
		case invites.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case invites.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("failed to retrieve invites: %w", err)
		}
	}

	var res []invites.InviteEnhanced
	for _, invite := range is {
		team, err := t.query.team.Retrieve(r.Context(), t.repo, invite.TeamID)
		if err != nil {
			switch err {
			case teams.ErrInvalidID:
				return web.NewRequestError(err, http.StatusBadRequest)
			case teams.ErrNotFound:
				return web.NewRequestError(err, http.StatusNotFound)
			default:
				return fmt.Errorf("failed to retrieve team: %w", err)
			}
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

// UpdateInvite updates an existing invitation
func (t *Team) UpdateInvite(w http.ResponseWriter, r *http.Request) error {
	var update invites.UpdateInvite
	var role memberships.Role = memberships.Editor

	if err := web.Decode(r, &update); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	uid := t.auth0.UserByID(r.Context())
	tid := chi.URLParam(r, "tid")
	iid := chi.URLParam(r, "iid")

	iv, err := t.query.invite.Update(r.Context(), t.repo, update, uid, iid, time.Now())
	if err != nil {
		switch err {
		case invites.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case invites.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("failed to update invite: %w", err)
		}
	}

	if update.Accepted {
		nm := memberships.NewMembership{
			UserID: uid,
			TeamID: tid,
			Role:   role.String(),
		}
		m, err := t.query.membership.Create(r.Context(), t.repo, nm, time.Now())
		if err != nil {
			return fmt.Errorf("failed to insert membership: %w", err)
		}

		if t.nats != nil {
			err = t.publish.MembershipCreated(t.nats, m, uid)
			if err != nil {
				return err
			}
		}
	}

	return web.Respond(r.Context(), w, iv, http.StatusOK)
}
