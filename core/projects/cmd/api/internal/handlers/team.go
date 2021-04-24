package handlers

import (
	"github.com/go-chi/chi"
	"github.com/ivorscott/devpie-client-backend-go/internal/platform/mauth"
	"github.com/pkg/errors"
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
	repo    *database.Repository
	log     *log.Logger
	auth0   *mid.Auth0
	origins string
}

func (t *Team) Create(w http.ResponseWriter, r *http.Request) error {
	var nt team.NewTeam

	pid := chi.URLParam(r, "pid")
	uid := t.auth0.GetUserBySubject(r)

	if err := web.Decode(r, &nt); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	ts, err := team.Create(r.Context(), t.repo, nt, pid, uid, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, ts, http.StatusCreated)
}

func (t *Team) Retrieve(w http.ResponseWriter, r *http.Request) error {
	pid := chi.URLParam(r, "pid")

	ts, err := team.Retrieve(r.Context(), t.repo, pid)
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

	return web.Respond(r.Context(), w, ts, http.StatusOK)
}

func (t *Team) Invite(w http.ResponseWriter, r *http.Request) error {
	//tid := chi.URLParam(r, "tid")
	var token *mauth.Token
	var invite mauth.NewInvite

	if err := web.Decode(r, &invite); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	token, err := mauth.Retrieve(r.Context(), t.repo)
	if err == mauth.ErrNotFound || mauth.IsExpired(token, t.auth0.GetPemCert) {
		// create new management api token
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

	resultUrl := strings.Split(t.origins, ",")[0]

	for _, email := range invite.Emails {
		user, err := mauth.CreateUser(token, t.auth0.Domain, email)
		if err != nil {
			return err
		}

		ticket, err := mauth.ChangePasswordTicket(token, t.auth0.Domain, user, resultUrl)
		if err != nil {
			return err
		}

		t.SendInvite(ticket)
	}

	// TODO: Create Team member in database

	// Create Team Member in database
	//in, err := team.CreateMember(r.Context(), t.repo, ni, tid, time.Now())
	//if err != nil {
	//	return err
	//}

	return web.Respond(r.Context(), w, nil, http.StatusCreated)
}

func (t *Team) SendInvite(ticket string) {
	//TODO: send email with send grid
	log.Printf("Sending invite.... %s", ticket)
}
