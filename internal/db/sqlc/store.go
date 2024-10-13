package db

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrorUniqueViolation = errors.New("unique_violation")
	ErrConflict          = errors.New("resource already exists")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicateUsername = errors.New("username already exists")
)

type Store interface {
	Querier
	CreateAndInviteUser(ctx context.Context, token string, exp time.Duration, arg CreateUserParams) (CreateUserRow, error)
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

func (s *SQLStore) CreateAndInviteUser(ctx context.Context, token string, exp time.Duration, arg CreateUserParams) (CreateUserRow, error) {
	tx, err := s.connPool.BeginTx(ctx, nil)
	if err != nil {
		return CreateUserRow{}, err
	}

	defer tx.Rollback()

	user, err := s.CreateUser(ctx, arg)
	if err != nil {
		return CreateUserRow{}, err
	}

	_, err = s.CreateUserInvitation(ctx, CreateUserInvitationParams{
		Token:  []byte(token),
		UserID: user.ID,
		Expiry: time.Now().Add(exp),
	})
	if err != nil {
		return CreateUserRow{}, err
	}

	if err := tx.Commit(); err != nil {
		return CreateUserRow{}, err
	}

	return user, nil
}
