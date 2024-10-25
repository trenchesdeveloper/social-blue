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
	ErrInvalidToken      = errors.New("invalid token")
)

type Store interface {
	Querier
	CreateAndInviteUser(ctx context.Context, token string, exp time.Duration, arg CreateUserParams) (CreateUserRow, error)
	ActivateUser(ctx context.Context, token string) (GetUserFromInvitationRow, error)
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

func (s *SQLStore) ActivateUser(ctx context.Context, token string) (GetUserFromInvitationRow, error) {
	tx, err := s.connPool.BeginTx(ctx, nil)
	if err != nil {
		return GetUserFromInvitationRow{}, err
	}

	defer tx.Rollback()

	user, err := s.GetUserFromInvitation(ctx, GetUserFromInvitationParams{
		Token:  []byte(token),
		Expiry: time.Now(),
	})

	if err != nil {
		return GetUserFromInvitationRow{}, err
	}

	// update the user to be active
	err = s.UpdateUserActivation(ctx, UpdateUserActivationParams{
		ID:       user.ID,
		IsActive: true,
	})

	if err != nil {
		return GetUserFromInvitationRow{}, err
	}

	// delete the invitation
	err = s.DeleteUserInvitation(ctx, user.ID)

	if err != nil {
		return GetUserFromInvitationRow{}, err
	}

	if err := tx.Commit(); err != nil {
		return GetUserFromInvitationRow{}, err
	}

	return user, nil
}
