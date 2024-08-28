package db

import (
	"context"
	"database/sql"
)

type Post struct {
	ID        int64    `json:"id"`
	Content   string   `json:"content"`
	Title     string   `json:"title"`
	UserID    int64    `json:"user_id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type PostsStore struct {
	db *sql.DB
}

func NewPostStore(db *sql.DB) *PostsStore {
	return &PostsStore{
		db: db,
	}
}

func (p *PostsStore) Create(ctx context.Context, post *Post) error {
	return nil

}
