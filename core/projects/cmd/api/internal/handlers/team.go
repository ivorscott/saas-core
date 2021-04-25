package handlers

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/mauth"
	"github.com/ivorscott/devpie-client-backend-go/internal/user"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ivorscott/devpie-client-backend-go/internal/mid"
	"github.com/ivorscott/devpie-client-backend-go/internal/team"

	"github.com/ivorscott/devpie-client-backend-go/internal/platform/database"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/web"
)

type Team struct {
	repo           *database.Repository
	log            *log.Logger
	auth0          *mid.Auth0
	origins        string
	sendgridAPIKey string
}

func (t *Team) Create(w http.ResponseWriter, r *http.Request) error {
	// TODO: Create Team member for team leader in database

	var nt team.NewTeam

	pid := chi.URLParam(r, "pid")
	uid := t.auth0.GetUserBySubject(r)

	if err := web.Decode(r, &nt); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	tm, err := team.Create(r.Context(), t.repo, nt, pid, uid, time.Now())
	if err != nil {
		return err
	}

	nm := team.NewMember{
		UserID:         uid,
		TeamID:         tm.ID,
		IsLeader:       true,
		InviteAccepted: true,
	}

	m, err := team.CreateMember(r.Context(), t.repo, nm, time.Now())
	if err != nil {
		return err
	}

	res := struct {
		Team   team.Team
		Member team.Member
	}{tm, m}

	return web.Respond(r.Context(), w, res, http.StatusCreated)
}

func (t *Team) Retrieve(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")

	tm, err := team.Retrieve(r.Context(), t.repo, pid)
	if err != nil {
		switch err {
		case team.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case team.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking for team connected to project %q", pid)
		}
	}

	return web.Respond(r.Context(), w, tm, http.StatusOK)
}

func (t *Team) Invite(w http.ResponseWriter, r *http.Request) error {
	var token *mauth.Token
	var invite mauth.NewInvite
	var link string

	if err := web.Decode(r, &invite); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	// Get valid management api token
	token, err := mauth.Retrieve(r.Context(), t.repo)
	if err == mauth.ErrNotFound || mauth.IsExpired(token, t.auth0.GetPemCert) {
		token, err = mauth.NewManagementToken(t.auth0.Domain, t.auth0.M2MClient, t.auth0.M2MSecret, t.auth0.MAPIAudience)
		if err != nil {
			return err
		}
		// clean table before persisting
		if err := mauth.Delete(r.Context(), t.repo); err != nil {
			return err
		}
		// persist management api token
		if err := mauth.Persist(r.Context(), t.repo, token, time.Now()); err != nil {
			return err
		}
	}

	tid := chi.URLParam(r, "tid")
	link = strings.Split(t.origins, ",")[0]

	for _, email := range invite.Emails {

		nm := team.NewMember{
			TeamID:         tid,
			IsLeader:       false,
			InviteSent:     true,
			InviteAccepted: false,
		}

		u, err := user.RetrieveByEmail(t.repo,email)
		if err != nil {
			 invitee, err := mauth.CreateUser(token, t.auth0.Domain, email)
			 if err != nil {
				 return err
			 }

			nm.UserID = invitee.UserId

			link, err = mauth.ChangePasswordTicket(token, t.auth0.Domain, invitee, link)
			if err != nil {
				 return err
			}
		 } else {
			nm.UserID = u.Auth0ID
		}

		t.SendInvite(email, link)

		// TODO: Create Team member in database
		_, err = team.CreateMember(r.Context(), t.repo, nm, time.Now())
		if err != nil {
			return err
		}
	}

	return web.Respond(r.Context(), w, nil, http.StatusCreated)
}

func (t *Team) SendInvite(email, link string) {
	html := ""
	html += "<strong>Join Devpie</strong>"
	html += "<br/>"
	html += "<p>To accept your invitation, <a href=\"%s\">create an account</a>.</p>"

	from := mail.NewEmail("DevPie", "people@devpie.io")
	subject := "You've been invited to a Project on DevPie!"
	to := mail.NewEmail("Invitee", email)
	plainTextContent := "What is this used for exactly.... subtitle???"
	htmlContent := fmt.Sprintf(html, link)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(t.sendgridAPIKey)
	response, err := client.Send(message)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
