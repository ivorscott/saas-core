package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/devpies/saas-core/internal/user/fail"
	"github.com/devpies/saas-core/internal/user/model"
	"github.com/devpies/saas-core/pkg/msg"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/go-chi/chi/v5"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
)

type teamService interface {
	Create(ctx context.Context, nt model.NewTeam, now time.Time) (model.Team, error)
	Retrieve(ctx context.Context, tid string) (model.Team, error)
	List(ctx context.Context) ([]model.Team, error)
}

type projectService interface {
	Create(ctx context.Context, p model.ProjectCopy) error
	Retrieve(ctx context.Context, pid string) (model.ProjectCopy, error)
	Update(ctx context.Context, pid string, update model.UpdateProjectCopy) error
	Delete(ctx context.Context, pid string) error
}

type inviteService interface {
	Create(ctx context.Context, ni model.NewInvite, now time.Time) (model.Invite, error)
	RetrieveInvite(ctx context.Context, iid string) (model.Invite, error)
	RetrieveInvites(ctx context.Context) ([]model.Invite, error)
	Update(ctx context.Context, update model.UpdateInvite, iid string, now time.Time) (model.Invite, error)
}

type publisher interface {
	Publish(subject string, message []byte)
}

// TeamHandler handles the team requests.
type TeamHandler struct {
	logger            *zap.Logger
	js                publisher
	sendgridAPIKey    string
	teamService       teamService
	projectService    projectService
	membershipService membershipService
	inviteService     inviteService
	userService       userService
}

// NewTeamHandler returns a new team handler.
func NewTeamHandler(
	logger *zap.Logger,
	js publisher,
	sendgridAPIKey string,
	teamService teamService,
	projectService projectService,
	membershipService membershipService,
	inviteService inviteService,
	userService userService,
) *TeamHandler {
	return &TeamHandler{
		logger:            logger,
		js:                js,
		sendgridAPIKey:    sendgridAPIKey,
		teamService:       teamService,
		projectService:    projectService,
		membershipService: membershipService,
		inviteService:     inviteService,
		userService:       userService,
	}
}

// CreateTeamForProject creates a new team for a project.
func (th *TeamHandler) CreateTeamForProject(w http.ResponseWriter, r *http.Request) error {
	values, ok := web.FromContext(r.Context())
	if !ok {
		return web.CtxErr()
	}

	var nt model.NewTeam
	var role model.Role = model.Administrator

	if err := web.Decode(r, &nt); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	if _, err := th.projectService.Retrieve(r.Context(), nt.ProjectID); err != nil {
		switch err {
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("failed to retrieve project: %w", err)
		}
	}

	tm, err := th.teamService.Create(r.Context(), nt, time.Now())
	if err != nil {
		return err
	}

	nm := model.NewMembership{
		UserID: values.UserID,
		TeamID: tm.ID,
		Role:   role.String(),
	}

	m, err := th.membershipService.Create(r.Context(), nm, time.Now())
	if err != nil {
		return err
	}

	up := model.UpdateProjectCopy{
		TeamID: &tm.ID,
	}

	if err = th.projectService.Update(r.Context(), nt.ProjectID, up); err != nil {
		return err
	}

	e := msg.TeamAssignedEvent{
		Type: msg.TypeTeamAssignedEventType,
		Data: msg.TeamAssignedEventData{
			MembershipID: m.ID,
			TenantID:     values.TenantID,
			TeamID:       m.TeamID,
			Role:         m.Role,
			UserID:       m.UserID,
			ProjectID:    nt.ProjectID,
			UpdatedAt:    m.UpdatedAt.String(),
			CreatedAt:    m.CreatedAt.String(),
		},
		Metadata: msg.Metadata{
			TenantID: values.TenantID,
			TraceID:  values.TraceID,
			UserID:   values.UserID,
		},
	}

	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}

	th.js.Publish(msg.SubjectProjectTeamAssigned, bytes)

	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, tm, http.StatusCreated)
}

// AssignExistingTeam assigns an existing team to a project.
func (th *TeamHandler) AssignExistingTeam(w http.ResponseWriter, r *http.Request) error {
	values, ok := web.FromContext(r.Context())
	if !ok {
		return web.CtxErr()
	}

	tid := chi.URLParam(r, "tid")
	pid := chi.URLParam(r, "pid")
	uid := values.UserID

	tm, err := th.teamService.Retrieve(r.Context(), tid)
	if err != nil {
		switch err {
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("failed to retrieve team: %w", err)
		}
	}

	var up = model.UpdateProjectCopy{
		TeamID:    &tm.ID,
		UpdatedAt: time.Now().UTC(),
	}

	err = th.projectService.Update(r.Context(), pid, up)
	if err != nil {
		switch err {
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return fmt.Errorf("failed to update project: %w", err)
		}
	}
	ue := msg.ProjectUpdatedEvent{
		Type: msg.TypeProjectUpdated,
		Data: msg.ProjectUpdatedEventData{
			TeamID:    &tid,
			ProjectID: pid,
			UpdatedAt: time.Now().UTC().String(),
		},
		Metadata: msg.Metadata{
			TenantID: values.TenantID,
			TraceID:  values.TraceID,
			UserID:   uid,
		},
	}

	bytes, err := json.Marshal(ue)
	if err != nil {
		return err
	}

	th.js.Publish(msg.SubjectProjectUpdated, bytes)

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

// LeaveTeam destroys a team membership.
func (th *TeamHandler) LeaveTeam(w http.ResponseWriter, r *http.Request) error {
	values, ok := web.FromContext(r.Context())
	if !ok {
		return web.CtxErr()
	}

	tid := chi.URLParam(r, "tid")
	uid := values.UserID

	mid, err := th.membershipService.Delete(r.Context(), tid)
	if err != nil {
		switch err {
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("failed to delete membership: %w", err)
		}
	}
	me := msg.MembershipDeletedEvent{
		Type: msg.TypeMembershipDeleted,
		Data: msg.MembershipDeletedEventData{
			MembershipID: mid,
		},
		Metadata: msg.Metadata{
			TenantID: values.TenantID,
			TraceID:  values.TraceID,
			UserID:   uid,
		},
	}

	bytes, err := json.Marshal(me)
	if err != nil {
		return err
	}

	th.js.Publish(msg.SubjectMembershipDeleted, bytes)
	return web.Respond(r.Context(), w, nil, http.StatusOK)
}

// List returns all teams associated with the authenticated user.
func (th *TeamHandler) List(w http.ResponseWriter, r *http.Request) error {
	teams, err := th.teamService.List(r.Context())
	if err != nil {
		switch err {
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("failed to retrieve teams: %w", err)
		}
	}

	return web.Respond(r.Context(), w, teams, http.StatusOK)
}

func (th *TeamHandler) Retrieve(w http.ResponseWriter, r *http.Request) error {
	tid := chi.URLParam(r, "tid")

	team, err := th.teamService.Retrieve(r.Context(), tid)
	if err != nil {
		switch err {
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("failed to retrieve team: %w", err)
		}
	}

	return web.Respond(r.Context(), w, team, http.StatusOK)
}

// Sender describes the function dependency required by SendMail
type Sender func(email *mail.SGMailV3) (*rest.Response, error)

// SendMail sends mail via the sendgrid client
func SendMail(email *mail.SGMailV3, send Sender) (*rest.Response, error) {
	var resp *rest.Response

	resp, err := send(email)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// CreateInvite sends new team invitations.
func (th *TeamHandler) CreateInvite(w http.ResponseWriter, r *http.Request) error {
	var (
		list model.NewList
		err  error
	)

	tid := chi.URLParam(r, "tid")
	link := "http://localhost/activation"

	if err = web.Decode(r, &list); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	for _, email := range list.Emails {
		ni := model.NewInvite{
			TeamID: tid,
		}
		// when user exists
		var user model.User
		user, err = th.userService.RetrieveByEmail(r.Context(), email)
		if err != nil {
			name := strings.Split(email, "@")
			nu := model.NewUser{
				Email:     email,
				FirstName: name[0],
			}

			err = th.userService.AddUser(r.Context(), nu, time.Now())
			if err != nil {
				return err
			}
			return web.Respond(r.Context(), w, nil, http.StatusAccepted)
		}
		ni.UserID = user.ID

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

		_, err = SendMail(message, sendgrid.NewSendClient(th.sendgridAPIKey).Send)
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}

		_, err = th.inviteService.Create(r.Context(), ni, time.Now())
		if err != nil {
			return err
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusCreated)
}

// RetrieveInvites returns invitations for the authenticated user.
func (th *TeamHandler) RetrieveInvites(w http.ResponseWriter, r *http.Request) error {
	var (
		is  []model.Invite
		err error
	)

	is, err = th.inviteService.RetrieveInvites(r.Context())
	if err != nil {
		switch err {
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("failed to retrieve invites: %w", err)
		}
	}

	var result = make([]model.InviteEnhanced, 0)

	for _, invite := range is {
		var team model.Team

		team, err = th.teamService.Retrieve(r.Context(), invite.TeamID)
		if err != nil {
			switch err {
			case fail.ErrInvalidID:
				return web.NewRequestError(err, http.StatusBadRequest)
			case fail.ErrNotFound:
				return web.NewRequestError(err, http.StatusNotFound)
			default:
				return fmt.Errorf("failed to retrieve team: %w", err)
			}
		}

		ie := model.InviteEnhanced{
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

		result = append(result, ie)
	}

	return web.Respond(r.Context(), w, result, http.StatusOK)
}

// UpdateInvite updates an existing invitation.
func (th *TeamHandler) UpdateInvite(w http.ResponseWriter, r *http.Request) error {
	var (
		update model.UpdateInvite
		role   model.Role = model.Editor
		err    error
	)

	if err = web.Decode(r, &update); err != nil {
		return err
	}

	values, ok := web.FromContext(r.Context())
	if !ok {
		return web.CtxErr()
	}

	uid := values.UserID
	tid := chi.URLParam(r, "tid")
	iid := chi.URLParam(r, "iid")

	iv, err := th.inviteService.Update(r.Context(), update, iid, time.Now())
	if err != nil {
		switch err {
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return fmt.Errorf("failed to update invite: %w", err)
		}
	}

	if update.Accepted {
		var m model.Membership

		nm := model.NewMembership{
			UserID: uid,
			TeamID: tid,
			Role:   role.String(),
		}

		m, err = th.membershipService.Create(r.Context(), nm, time.Now())
		if err != nil {
			return fmt.Errorf("failed to insert membership: %w", err)
		}

		e := msg.MembershipCreatedEvent{
			Type: msg.TypeMembershipCreated,
			Data: msg.MembershipCreatedEventData{
				MembershipID: m.ID,
				TenantID:     values.TenantID,
				TeamID:       m.TeamID,
				Role:         m.Role,
				UserID:       m.UserID,
				UpdatedAt:    m.UpdatedAt.String(),
				CreatedAt:    m.CreatedAt.String(),
			},
			Metadata: msg.Metadata{
				TenantID: values.TenantID,
				TraceID:  values.TraceID,
				UserID:   uid,
			},
		}

		var bytes []byte

		bytes, err = json.Marshal(e)
		if err != nil {
			return err
		}

		th.js.Publish(msg.SubjectMembershipCreated, bytes)
	}

	return web.Respond(r.Context(), w, iv, http.StatusOK)
}
