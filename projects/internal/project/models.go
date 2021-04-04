package project

import (
	"time"
)

type Project struct {
	ID          string    `db:"project_id" json:"id"`
	UserID      string    `db:"user_id" json:"userId"`
	Name        string    `db:"name" json:"name"`
	Open        bool      `db:"open" json:"open"`
	ColumnOrder []string  `db:"column_order" json:"columnOrder"`
	Created     time.Time `db:"created" json:"created"`
}

type NewProject struct {
	Name string `db:"name" json:"name"`
}

type UpdateProject struct {
	Name        string   `db:"name" json:"name"`
	Open        bool     `db:"open" json:"open"`
	ColumnOrder []string `db:"column_order" json:"columnOrder"`
}
