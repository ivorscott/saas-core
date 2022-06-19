package service

import (
	"context"
	"github.com/devpies/saas-core/internal/user/model"
	"github.com/devpies/saas-core/pkg/msg"
	"go.uber.org/zap"
)

type projectRepository interface {
	Create(ctx context.Context, p model.ProjectCopy) error
	Retrieve(ctx context.Context, pid string) (model.ProjectCopy, error)
	Update(ctx context.Context, pid string, update model.UpdateProjectCopy) error
	Delete(ctx context.Context, pid string) error
}

// ProjectService manages the redundant copy of project data.
type ProjectService struct {
	logger      *zap.Logger
	projectRepo projectRepository
}

// NewProjectService returns a new project service.
func NewProjectService(
	logger *zap.Logger,
	projectRepo projectRepository,
) *ProjectService {
	return &ProjectService{
		logger:      logger,
		projectRepo: projectRepo,
	}
}

func (ps *ProjectService) Create(ctx context.Context, p model.ProjectCopy) error {
	return ps.projectRepo.Create(ctx, p)
}

func (ps *ProjectService) Retrieve(ctx context.Context, pid string) (model.ProjectCopy, error) {
	return ps.projectRepo.Retrieve(ctx, pid)
}

func (ps *ProjectService) Update(ctx context.Context, pid string, update model.UpdateProjectCopy) error {
	return ps.projectRepo.Update(ctx, pid, update)
}

func (ps *ProjectService) Delete(ctx context.Context, pid string) error {
	return ps.projectRepo.Delete(ctx, pid)
}

func (ps *ProjectService) CreateProjectCopyFromEvent(ctx context.Context, message interface{}) error {
	m, err := msg.Bytes(message)
	if err != nil {
		return err
	}
	event, err := msg.UnmarshalProjectCreatedEvent(m)
	if err != nil {
		return err
	}

	data := newProject(event.Data)

	return ps.projectRepo.Create(ctx, data)
}

func (ps *ProjectService) UpdateProjectCopyFromEvent(ctx context.Context, message interface{}) error {
	m, err := msg.Bytes(message)
	if err != nil {
		return err
	}
	event, err := msg.UnmarshalProjectUpdatedEvent(m)
	if err != nil {
		return err
	}

	data := newUpdateProject(event.Data)

	return ps.projectRepo.Update(ctx, event.Data.ProjectID, data)
}

func (ps *ProjectService) DeleteProjectCopyFromEvent(ctx context.Context, message interface{}) error {
	m, err := msg.Bytes(message)
	if err != nil {
		return err
	}
	event, err := msg.UnmarshalProjectDeletedEvent(m)
	if err != nil {
		return err
	}

	return ps.projectRepo.Delete(ctx, event.Data.ProjectID)
}

func newProject(data msg.ProjectCreatedEventData) model.ProjectCopy {
	return model.ProjectCopy{
		ID:          data.ProjectID,
		TenantID:    data.TenantID,
		Name:        data.Name,
		Prefix:      data.Prefix,
		Description: data.Description,
		TeamID:      data.TeamID,
		UserID:      data.UserID,
		Active:      data.Active,
		Public:      data.Public,
		ColumnOrder: data.ColumnOrder,
		UpdatedAt:   msg.ParseTime(data.UpdatedAt),
		CreatedAt:   msg.ParseTime(data.CreatedAt),
	}
}

func newUpdateProject(data msg.ProjectUpdatedEventData) model.UpdateProjectCopy {
	return model.UpdateProjectCopy{
		Name:        data.Name,
		Active:      data.Active,
		Public:      data.Public,
		TeamID:      data.TeamID,
		ColumnOrder: data.ColumnOrder,
		Description: data.Description,
		UpdatedAt:   msg.ParseTime(data.UpdatedAt),
	}
}
