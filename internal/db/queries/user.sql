-- name: CreateUser :one
INSERT INTO users (first_name, last_name, username, password, email, role_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, first_name, last_name, username, email,role_id, created_at, updated_at;

-- name: GetUserByID :one
SELECT *
FROM users JOIN roles ON users.role_id = roles.id
WHERE users.id = $1;


-- name: GetUserByUsername :one
SELECT id, username, password, email, created_at, updated_at, is_active, role_id
FROM users
WHERE username = $1;

-- name: ListUsers :many
SELECT id, username, email, created_at, updated_at, is_active
FROM users
ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users
SET username = $2, email = $3, updated_at = now()
WHERE id = $1
RETURNING id, username, email, created_at, updated_at, is_active;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetUserFromInvitation :one
SELECT u.id, u.username, u.email, u.created_at, u.updated_at, u.is_active, u.role_id
FROM users u
JOIN user_invitations ui ON u.id = ui.user_id
WHERE ui.token = $1 AND ui.expiry > $2;


-- name: UpdateUserActivation :exec
UPDATE users
SET is_active = $2
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, username, email, password, created_at, updated_at, is_active, role_id
FROM users
WHERE email = $1;

-- name: GetActiveUserByEmail :one
SELECT id, username, email, password, created_at, updated_at, is_active, role_id
FROM users
WHERE email = $1 AND is_active = true;