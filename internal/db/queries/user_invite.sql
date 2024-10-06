-- name: CreateUserInvitation :one
INSERT INTO user_invitations (token, user_id)
VALUES ($1, $2)
RETURNING token, user_id;


-- name: GetUserInvitationByToken :one
SELECT token, user_id
FROM user_invitations
WHERE token = $1;


-- name: ListUserInvitations :many
SELECT token, user_id
FROM user_invitations
WHERE user_id = $1;


-- name: DeleteUserInvitation :exec
DELETE FROM user_invitations
WHERE token = $1;


-- name: UpdateUserInvitationUserID :exec
UPDATE user_invitations
SET user_id = $2
WHERE token = $1;
