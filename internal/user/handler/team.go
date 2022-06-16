package handler

import (
	"context"
	"github.com/devpies/saas-core/internal/user/model"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type teamService interface {
	Create(ctx context.Context, nt model.NewTeam, uid string, now time.Time) (model.Team, error)
	Retrieve(ctx context.Context, tid string) (model.Team, error)
	List(ctx context.Context, uid string) ([]model.Team, error)
}

// TeamHandler handles the team requests.
type TeamHandler struct {
	logger      *zap.Logger
	teamService teamService
}

// NewTeamHandler returns a new team handler.
func NewTeamHandler(
	logger *zap.Logger,
	teamService teamService,
) *TeamHandler {
	return &TeamHandler{
		logger:      logger,
		teamService: teamService,
	}
}

func (th *TeamHandler) Create(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (th *TeamHandler) AssignExistingTeam(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (th *TeamHandler) LeaveTeam(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (th *TeamHandler) List(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (th *TeamHandler) Retrieve(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (th *TeamHandler) CreateInvite(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (th *TeamHandler) RetrieveInvites(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (th *TeamHandler) UpdateInvite(w http.ResponseWriter, r *http.Request) error {
	return nil
}
