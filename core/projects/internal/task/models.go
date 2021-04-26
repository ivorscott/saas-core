package task

import (
	"time"
)

type Task struct {
	ID        string    `db:"task_id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	ProjectID string    `db:"project_id" json:"projectId"`
	Created   time.Time `db:"created" json:"created"`
}

type NewTask struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content"`
}

type UpdateTask struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

type MoveTask struct {
	To      string   `json:"to"`
	From    string   `json:"from"`
	TaskIds []string `json:"taskIds"`
}
