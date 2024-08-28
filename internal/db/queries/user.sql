-- name: CreateUser :one
INSERT INTO users (first_name, last_name, username, password, email)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, first_name, last_name, username, email, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, username, email, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByUsername :one
SELECT id, username, password, email, created_at, updated_at
FROM users
WHERE username = $1;

-- name: ListUsers :many
SELECT id, username, email, created_at, updated_at
FROM users
ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users
SET username = $2, email = $3, updated_at = now()
WHERE id = $1
RETURNING id, username, email, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
