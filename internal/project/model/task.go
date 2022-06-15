package model

import (
	"time"
)

// Task represents a Project Task.
type Task struct {
	ID          string    `db:"task_id" json:"id"`
	Key         string    `db:"key" json:"key"`
	Seq         int       `db:"seq" json:"seq"`
	Title       string    `db:"title" json:"title"`
	Points      int       `db:"points" json:"points"`
	Content     string    `db:"content" json:"content"`
	ProjectID   string    `db:"project_id" json:"projectId"`
	AssignedTo  string    `db:"assigned_to" json:"assignedTo"`
	Attachments []string  `db:"attachments" json:"attachments"`
	Comments    []string  `db:"comments" json:"comments"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}

// NewTask represents a new Task.
type NewTask struct {
	Title string `json:"title" validate:"required"`
}

// UpdateTask represents a Task being updated.
type UpdateTask struct {
	Title       *string   `json:"title"`
	Key         *string   `json:"key"`
	Points      *int      `json:"points"`
	Content     *string   `json:"content"`
	AssignedTo  *string   `json:"assignedTo"`
	Attachments []string  `json:"attachments"`
	Comments    []string  `json:"comments"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// MoveTask represents a Task being moved between columns.
type MoveTask struct {
	To      string   `json:"to"`
	From    string   `json:"from"`
	TaskIds []string `json:"taskIds"`
}
