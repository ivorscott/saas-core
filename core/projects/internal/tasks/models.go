package tasks

import (
	"time"
)

type Task struct {
	ID        string    `db:"task_id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	ProjectID string    `db:"project_id" json:"projectId"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type NewTask struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content"`
}

type UpdateTask struct {
	Title     *string   `json:"title"`
	Content   *string   `json:"content"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type MoveTask struct {
	To      string   `json:"to"`
	From    string   `json:"from"`
	TaskIds []string `json:"taskIds"`
}
