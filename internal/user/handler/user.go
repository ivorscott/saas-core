package handler

import (
	"fmt"
	"github.com/devpies/saas-core/internal/user/fail"
	"github.com/devpies/saas-core/internal/user/model"
	"github.com/devpies/saas-core/pkg/web"
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

// Create adds a new seat to the tenant account. The tenant admin is
// stored automatically through listening to the TENANTS.registered event and a separate concern.
func (uh *UserHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var (
		nu  model.NewUser
		err error
	)

	if err = web.Decode(r, &nu); err != nil {
		return err
	}

	user, err := uh.userService.RetrieveMe(r.Context())
	if err != nil {
		err = uh.userService.AddSeat(r.Context(), nu, time.Now())
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		return web.Respond(r.Context(), w, user, http.StatusCreated)
	}

	return web.Respond(r.Context(), w, nil, http.StatusAccepted)
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
			return fmt.Errorf("failed to retrieve authenticated user: %w", err)
		}
	}

	return web.Respond(r.Context(), w, us, http.StatusOK)
}
