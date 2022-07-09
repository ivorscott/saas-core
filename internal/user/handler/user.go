package handler

import (
	"github.com/devpies/saas-core/internal/user/fail"
	"github.com/devpies/saas-core/internal/user/model"
	"github.com/devpies/saas-core/pkg/web"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// UserHandler handles the user requests.
type UserHandler struct {
	logger      *zap.Logger
	userService userService
}

// NewUserHandler returns a new user handler.
func NewUserHandler(
	logger *zap.Logger,
	userService userService,
) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
	}
}

// Create adds a new seat to the tenant account. Not to be confused with the tenant admin user.
// The tenant admin user is defined automatically during initial registration of the tenant account.
func (uh *UserHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var (
		nu  model.NewUser
		err error
	)

	if err = web.Decode(r, &nu); err != nil {
		return err
	}

	err = uh.userService.AddUser(r.Context(), nu, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, nil, http.StatusCreated)
}

// List retrieves all users on the tenant account.
func (uh *UserHandler) List(w http.ResponseWriter, r *http.Request) error {
	var err error

	users, err := uh.userService.List(r.Context())
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, users, http.StatusOK)
}

func (uh *UserHandler) SeatsAvailable(w http.ResponseWriter, r *http.Request) error {
	result, err := uh.userService.SeatsAvailable(r.Context())
	if err != nil {
		return err
	}
	return web.Respond(r.Context(), w, result, http.StatusOK)
}

// RetrieveMe retrieves the authenticated user.
func (uh *UserHandler) RetrieveMe(w http.ResponseWriter, r *http.Request) error {
	var us model.User

	values, ok := web.FromContext(r.Context())
	if !ok {
		return web.CtxErr()
	}

	if values.UserID == "" {
		return web.NewRequestError(fail.ErrNotFound, http.StatusNotFound)
	}

	us, err := uh.userService.RetrieveMe(r.Context())
	if err != nil {
		switch err {
		case fail.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case fail.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return err
		}
	}

	return web.Respond(r.Context(), w, us, http.StatusOK)
}

func (uh *UserHandler) RemoveUser(w http.ResponseWriter, r *http.Request) error {

	uid := chi.URLParam(r, "uid")

	err := uh.userService.RemoveUser(r.Context(), uid)
	if err != nil {
		return err
	}

	return web.Respond(r.Context(), w, nil, http.StatusOK)
}
