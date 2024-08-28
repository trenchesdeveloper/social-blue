package store

import (
	"context"
	"database/sql"
)

type PostsStore struct {
	db *sql.DB
}

func NewPostStore(db *sql.DB) *PostsStore {
	return &PostsStore{
		db: db,
	}
}

func (p *PostsStore) Create(ctx context.Context) error {
	return nil

}
