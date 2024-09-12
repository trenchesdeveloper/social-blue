package db

import (
	"database/sql"
	"errors"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrorUniqueViolation = errors.New("unique_violation")
	ErrConflict          = errors.New("resource already exists")
)

type Store interface {
	Querier
}

type SQLStore struct {
	connPool *sql.DB
	*Queries
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		connPool: db,
		Queries:  New(db),
	}
}
