package model

import (
	"github.com/go-playground/validator/v10"
	"time"
)

var columnValidator *validator.Validate

func init() {
	v := NewValidator()
	columnValidator = v
}

// Column represents a Project Column.
type Column struct {
	ID         string    `db:"column_id" json:"id"`
	TenantID   string    `db:"tenant_id" json:"tenantID"`
	Title      string    `db:"title" json:"title"`
	ColumnName string    `db:"column_name" json:"columnName"`
	TaskIDS    []string  `db:"task_ids" json:"taskIds"`
	ProjectID  string    `db:"project_id" json:"projectId"`
	UpdatedAt  time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
}

// NewColumn represents a new Column.
type NewColumn struct {
	Title      string `json:"title" validate:"required,max=24"`
	ColumnName string `json:"columnName" validate:"required,max=10"`
	ProjectID  string `json:"projectId" validate:"required,max=36"`
}

func (nc *NewColumn) Validate() error {
	return columnValidator.Struct(nc)
}

// UpdateColumn represents a Column update.
type UpdateColumn struct {
	Title     *string   `json:"title" validate:"omitempty,max=24"`
	TaskIDS   *[]string `json:"taskIds"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (uc *UpdateColumn) Validate() error {
	return columnValidator.Struct(uc)
}
