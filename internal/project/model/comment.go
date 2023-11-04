package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var commentValidator *validator.Validate

func init() {
	v := NewValidator()
	commentValidator = v
}

// Comment represents a comment on a Task.
type Comment struct {
	ID        string    `db:"comment_id" json:"commentId"`
	Content   string    `db:"content" json:"content"`
	UserID    string    `db:"user_id" json:"userId"`
	Liked     bool      `db:"liked" json:"liked"`
	Edited    bool      `db:"edited" json:"edited"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// NewComment represents a new comment.
type NewComment struct {
	Content string `json:"content" validate:"required,max=500"`
}

// Validate validates a NewComment.
func (nc *NewComment) Validate() error {
	return commentValidator.Struct(nc)
}

// UpdateComment represents a comment update.
type UpdateComment struct {
	Content *string `json:"content" validate:"omitempty,max=500"`
	Liked   *bool   `json:"liked"`
}

// Validate validates an UpdateComment.
func (uc *UpdateComment) Validate() error {
	return commentValidator.Struct(uc)
}
