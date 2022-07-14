package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var projectValidator *validator.Validate

func init() {
	v := NewValidator()
	projectValidator = v
}

// Project represents a tenant Project.
type Project struct {
	ID          string    `db:"project_id" json:"id"`
	TenantID    string    `db:"tenant_id" json:"tenantId"`
	Name        string    `db:"name" json:"name"`
	Prefix      string    `db:"prefix" json:"prefix"`
	Description string    `db:"description" json:"description"`
	UserID      string    `db:"user_id" json:"userId"`
	Active      bool      `db:"active" json:"active"`
	Public      bool      `db:"public" json:"public"`
	ColumnOrder []string  `db:"column_order" json:"columnOrder"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}

// NewProject represents a new Project.
type NewProject struct {
	Name string `json:"name" validate:"required,max=22"`
}

func (np *NewProject) Validate() error {
	return projectValidator.Struct(np)
}

// UpdateProject represents a Project update.
type UpdateProject struct {
	Name        *string  `json:"name" validate:"omitempty,min=3,max=22"`
	Active      *bool    `json:"active"`
	Public      *bool    `json:"public"`
	Description *string  `json:"description" validate:"omitempty,max=72"`
	ColumnOrder []string `json:"columnOrder"`
}

func (up *UpdateProject) Validate() error {
	return projectValidator.Struct(up)
}
