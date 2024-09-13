package dto

import (
	db "github.com/trenchesdeveloper/social-blue/internal/db/sqlc"
	"time"
)

type CreatPostDto struct {
	Content string   `json:"content" validate:"required,max=1000"`
	Title   string   `json:"title" validate:"required,max=100"`
	Tags    []string `json:"tags"`
}

type UpdatePostDto struct {
	Content string   `json:"content" validate:"omitempty,max=1000"`
	Title   string   `json:"title" validate:"omitempty,max=100"`
	Tags    []string `json:"tags" validate:"omitempty"`
}

type PostWithMetadata struct {
	GetPostWithCommentsDto
	CommentsCount int64 `json:"comments_count"`
}

type GetPostWithCommentsDto struct {
	ID        int64                       `json:"id"`
	Content   string                      `json:"content"`
	Title     string                      `json:"title"`
	UserID    int64                       `json:"user_id"`
	Tags      []string                    `json:"tags"`
	CreatedAt time.Time                   `json:"created_at"`
	UpdatedAt time.Time                   `json:"updated_at"`
	Comments  []db.GetCommentsByPostIDRow `json:"comments"`
	User      db.User                     `json:"user"`
}
