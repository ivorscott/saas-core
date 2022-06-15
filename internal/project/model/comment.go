package model

import "time"

type Comment struct {
	ID        string    `db:"comment_id" json:"commentId"`
	Content   string    `db:"content" json:"content"`
	UserID    string    `db:"user_id" json:"userId"`
	Likes     int       `db:"likes" json:"likes"`
	Edited    bool      `db:"edited" json:"edited"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

type NewComment struct {
	Content string `json:"content"`
}

type UpdateComment struct {
	Content *string `json:"content"`
	Likes   *int    `json:"likes"`
}
