-- name: CreateUserInvitation :one
INSERT INTO user_invitations (token, user_id, expiry)
VALUES ($1, $2, $3)
RETURNING token, user_id, expiry;


-- name: GetUserInvitationByToken :one
SELECT token, user_id, expiry
FROM user_invitations
WHERE token = $1;


-- name: ListUserInvitations :many
SELECT token, user_id, expiry
FROM user_invitations
WHERE user_id = $1;


-- name: DeleteUserInvitation :exec
DELETE FROM user_invitations
WHERE token = $1;


-- name: UpdateUserInvitationUserID :exec
UPDATE user_invitations
SET user_id = $2
WHERE token = $1;
