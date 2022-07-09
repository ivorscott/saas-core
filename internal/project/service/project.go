package service

import (
	"context"
	"fmt"
	"github.com/devpies/saas-core/pkg/msg"
	"github.com/jmoiron/sqlx"
	"time"

	"github.com/devpies/saas-core/internal/project/model"

	"go.uber.org/zap"
)

type projectRepository interface {
	RunTx(ctx context.Context, fn func(*sqlx.Tx) error) error
	RetrieveTeamID(ctx context.Context, pid string) (string, error)
	Retrieve(ctx context.Context, pid string) (model.Project, error)
	List(ctx context.Context) ([]model.Project, error)
	Create(ctx context.Context, np model.NewProject, now time.Time) (model.Project, error)
	Update(ctx context.Context, pid string, update model.UpdateProject, now time.Time) (model.Project, error)
	UpdateTx(ctx context.Context, tx *sqlx.Tx, pid string, update model.UpdateProject, now time.Time) (model.Project, error)
	Delete(ctx context.Context, pid string) error
}

// ProjectService is responsible for managing project business logic.
type ProjectService struct {
	logger         *zap.Logger
	projectRepo    projectRepository
	membershipRepo membershipRepository
}

// NewProjectService returns a new ProjectService.
func NewProjectService(logger *zap.Logger, projectRepo projectRepository, membershipRepo membershipRepository) *ProjectService {
	return &ProjectService{
		logger:         logger,
		projectRepo:    projectRepo,
		membershipRepo: membershipRepo,
	}
}

// List lists projects.
func (ps *ProjectService) List(ctx context.Context, all bool) ([]model.Project, error) {
	if all {
		return forEachT(ctx, ps.projectRepo.List)
	}
	return ps.projectRepo.List(ctx)
}

// Retrieve retrieves an owned project.
func (ps *ProjectService) Retrieve(ctx context.Context, projectID string) (model.Project, error) {
	return ps.projectRepo.Retrieve(ctx, projectID)
}

// Create creates a project.
func (ps *ProjectService) Create(ctx context.Context, project model.NewProject, now time.Time) (model.Project, error) {
	return ps.projectRepo.Create(ctx, project, now)
}

// Update updates a project.
func (ps *ProjectService) Update(ctx context.Context, projectID string, update model.UpdateProject, now time.Time) (model.Project, error) {
	return ps.projectRepo.Update(ctx, projectID, update, now)
}

func (ps *ProjectService) AssignProjectTeamFromEvent(ctx context.Context, message interface{}) error {
	m, err := msg.Bytes(message)
	if err != nil {
		return err
	}
	event, err := msg.UnmarshalTeamAssignedEvent(m)
	if err != nil {
		return err
	}

	projectUpdate, membership := newTeamAssignment(event.Data)

	return ps.projectRepo.RunTx(ctx, func(tx *sqlx.Tx) error {
		tid := event.Metadata.TenantID

		_, err = tx.ExecContext(ctx, fmt.Sprintf("select set_config('app.current_tenant', '%s', false);", tid))
		if err != nil {
			ps.logger.Error("setting session variable failed", zap.Error(err))
			return err
		}

		err = ps.membershipRepo.CreateTx(ctx, tx, membership)
		if err != nil {
			return err
		}

		_, err = ps.projectRepo.UpdateTx(ctx, tx, event.Data.ProjectID, projectUpdate, msg.ParseTime(event.Data.UpdatedAt))
		if err != nil {
			return err
		}
		return nil
	})
}

func newTeamAssignment(data msg.TeamAssignedEventData) (model.UpdateProject, model.MembershipCopy) {
	return model.UpdateProject{
			TeamID: &data.TeamID,
		}, model.MembershipCopy{
			ID:        data.MembershipID,
			TenantID:  data.TenantID,
			UserID:    data.UserID,
			TeamID:    data.TeamID,
			Role:      data.Role,
			UpdatedAt: msg.ParseTime(data.UpdatedAt),
			CreatedAt: msg.ParseTime(data.CreatedAt),
		}
}

// Delete deletes a project.
func (ps *ProjectService) Delete(ctx context.Context, projectID string) error {
	return ps.projectRepo.Delete(ctx, projectID)
}
