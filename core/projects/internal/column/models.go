package column

import (
	"time"
)

type Column struct {
	ID         string    `db:"column_id" json:"id"`
	Title      string    `db:"title" json:"title"`
	ColumnName string    `db:"column_name" json:"columnName"`
	TaskIDS    []string  `db:"task_ids" json:"taskIds"`
	ProjectID  string    `db:"project_id" json:"projectId"`
	Created    time.Time `db:"created" json:"created"`
}

type NewColumn struct {
	Title      string `json:"title"`
	ColumnName string `json:"columnName"`
	ProjectID  string `json:"projectId"`
}

type UpdateColumn struct {
	Title   *string   `json:"title"`
	TaskIDS *[]string `json:"taskIds"`
}
