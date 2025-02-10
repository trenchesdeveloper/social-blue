-- name: GetRoleByName :one
SELECT id, name, level, description
FROM roles
WHERE name = $1;

-- name: GetRoleByID :one
SELECT id, name, level, description
FROM roles
WHERE id = $1;