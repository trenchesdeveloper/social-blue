// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user_invite.sql

package db

import (
	"context"
	"time"
)

const createUserInvitation = `-- name: CreateUserInvitation :one
INSERT INTO user_invitations (token, user_id, expiry)
VALUES ($1, $2, $3)
RETURNING token, user_id, expiry
`

type CreateUserInvitationParams struct {
	Token  []byte    `json:"token"`
	UserID int64     `json:"user_id"`
	Expiry time.Time `json:"expiry"`
}

func (q *Queries) CreateUserInvitation(ctx context.Context, arg CreateUserInvitationParams) (UserInvitation, error) {
	row := q.db.QueryRowContext(ctx, createUserInvitation, arg.Token, arg.UserID, arg.Expiry)
	var i UserInvitation
	err := row.Scan(&i.Token, &i.UserID, &i.Expiry)
	return i, err
}

const deleteUserInvitation = `-- name: DeleteUserInvitation :exec
DELETE FROM user_invitations
WHERE user_id = $1
`

func (q *Queries) DeleteUserInvitation(ctx context.Context, userID int64) error {
	_, err := q.db.ExecContext(ctx, deleteUserInvitation, userID)
	return err
}

const getUserInvitationByToken = `-- name: GetUserInvitationByToken :one
SELECT token, user_id, expiry
FROM user_invitations
WHERE token = $1
`

func (q *Queries) GetUserInvitationByToken(ctx context.Context, token []byte) (UserInvitation, error) {
	row := q.db.QueryRowContext(ctx, getUserInvitationByToken, token)
	var i UserInvitation
	err := row.Scan(&i.Token, &i.UserID, &i.Expiry)
	return i, err
}

const listUserInvitations = `-- name: ListUserInvitations :many
SELECT token, user_id, expiry
FROM user_invitations
WHERE user_id = $1
`

func (q *Queries) ListUserInvitations(ctx context.Context, userID int64) ([]UserInvitation, error) {
	rows, err := q.db.QueryContext(ctx, listUserInvitations, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []UserInvitation{}
	for rows.Next() {
		var i UserInvitation
		if err := rows.Scan(&i.Token, &i.UserID, &i.Expiry); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUserInvitationUserID = `-- name: UpdateUserInvitationUserID :exec
UPDATE user_invitations
SET user_id = $2
WHERE token = $1
`

type UpdateUserInvitationUserIDParams struct {
	Token  []byte `json:"token"`
	UserID int64  `json:"user_id"`
}

func (q *Queries) UpdateUserInvitationUserID(ctx context.Context, arg UpdateUserInvitationUserIDParams) error {
	_, err := q.db.ExecContext(ctx, updateUserInvitationUserID, arg.Token, arg.UserID)
	return err
}
