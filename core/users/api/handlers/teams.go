package handlers

import (
	"context"
	"encoding/json"
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
	repo        database.Storer
	log         *log.Logger
	auth0       auth0.Auther
	nats        *events.Client
	origins     string
	sendgridKey string
	query       TeamQueries
}

type TeamQueries struct {
	team       TeamQuerier
	project    ProjectQuerier
	membership MembershipQuerier
	user       UserQuerier
	invite     InviteQuerier
}

type TeamQuerier interface {
	Create(ctx context.Context, repo database.Storer, nt teams.NewTeam, uid string, now time.Time) (teams.Team, error)
	Retrieve(ctx context.Context, repo database.Storer, tid string) (teams.Team, error)
	List(ctx context.Context, repo database.Storer, uid string) ([]teams.Team, error)
}

type ProjectQuerier interface {
	Create(ctx context.Context, repo *database.Repository, p projects.ProjectCopy) error
	Retrieve(ctx context.Context, repo database.Storer, pid string) (projects.ProjectCopy, error)
	Update(ctx context.Context, repo database.Storer, pid string, update projects.UpdateProjectCopy) error
	Delete(ctx context.Context, repo database.Storer, pid string) error
}

type InviteQuerier interface {
	Create(ctx context.Context, repo database.Storer, ni invites.NewInvite, now time.Time) (invites.Invite, error)
	RetrieveInvite(ctx context.Context, repo database.Storer, uid string, iid string) (invites.Invite, error)
	RetrieveInvites(ctx context.Context, repo database.Storer, uid string) ([]invites.Invite, error)
	Update(ctx context.Context, repo database.Storer, update invites.UpdateInvite, uid, iid string, now time.Time) (invites.Invite, error)
}

func (t *Team) Create(w http.ResponseWriter, r *http.Request) error {
	var nt teams.NewTeam
	var role memberships.Role = memberships.Administrator

	uid := t.auth0.UserByID(r.Context()) // mock

	if err := web.Decode(r, &nt); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	if _, err := t.query.project.Retrieve(r.Context(), t.repo, nt.ProjectID); err != nil { //mock
		switch err {
		case projects.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case projects.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("failed to retrieve project: %w", err)
		}
	}

	tm, err := t.query.team.Create(r.Context(), t.repo, nt, uid, time.Now()) // mock
	if err != nil {
		return err
	}

	nm := memberships.NewMembership{
		UserID: uid,
		TeamID: tm.ID,
		Role:   role.String(),
	}

	m, err := t.query.membership.Create(r.Context(), t.repo, nm, time.Now()) // mock
	if err != nil {
		return err
	}

	up := projects.UpdateProjectCopy{
		TeamID: &tm.ID,
	}

	if err := t.query.project.Update(r.Context(), t.repo, nt.ProjectID, up); err != nil {
		return err
	} // mock

	err = t.PublishMembershipCreatedForProject(m, nt.ProjectID, uid)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, tm, http.StatusCreated)
}

func (t *Team)PublishMembershipCreatedForProject(m memberships.Membership, pid , uid string) error {
	e := events.MembershipCreatedForProjectEvent{
		ID:   uuid.New().String(),
		Type: events.TypeMembershipCreatedForProject,
		Data: events.MembershipCreatedForProjectEventData{
			MembershipID: m.ID,
			TeamID:       m.TeamID,
			Role:         m.Role,
			UserID:       m.UserID,
			ProjectID:    pid,
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
		t.nats.Publish(string(events.EventsMembershipCreatedForProject), bytes) // mock
	}
	return nil
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

	err = t.PublishProjectUpdateEvent(&tm.ID, pid,uid)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

func (t *Team) PublishProjectUpdateEvent(tid *string, pid, uid string) error {
	ue := events.ProjectUpdatedEvent{
		ID:   uuid.New().String(),
		Type: events.TypeProjectUpdated,
		Data: events.ProjectUpdatedEventData{
			TeamID:    tid,
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
	return nil
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

	err = t.PublishMembershipDeleted(mid, uid)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

func (t *Team) PublishMembershipDeleted(mid, uid string) error {
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

	return nil
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

	uid := t.auth0.UserByID(r.Context())
	tid := chi.URLParam(r, "tid")
	iid := chi.URLParam(r, "iid")

	if err := web.Decode(r, &update); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

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
			return err
		}

		err = t.PublishMembershipCreated(m, uid)
		if err != nil {
			return err
		}
	}

	return web.Respond(r.Context(), w, iv, http.StatusOK)
}

func(t *Team) PublishMembershipCreated(m memberships.Membership, uid string) error {
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

	return nil
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
	}

	t.log.Println(response.StatusCode)
	t.log.Println(response.Body)
	t.log.Println(response.Headers)

	return nil
}
