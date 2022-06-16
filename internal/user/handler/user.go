package handler

import (
	"context"
	"github.com/devpies/saas-core/internal/user/model"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type userService interface {
	Create(ctx context.Context, nu model.NewUser, now time.Time) (model.User, error)
	RetrieveByEmail(email string) (model.User, error)
	RetrieveMe(ctx context.Context, uid string) (model.User, error)
	RetrieveMeByAuthID(ctx context.Context, aid string) (model.User, error)
}

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

func (uh *UserHandler) Create(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (uh *UserHandler) RetrieveMe(w http.ResponseWriter, r *http.Request) error {
	return nil
}
