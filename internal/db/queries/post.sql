-- name: CreatePost :one
INSERT INTO posts (content, title, user_id, tags)
VALUES ($1, $2, $3, $4)
RETURNING id, content, title, user_id, tags, created_at, updated_at;

-- name: GetPostByID :one
SELECT id, content, title, user_id, tags, created_at, updated_at
FROM posts
WHERE id = $1;

-- name: ListPosts :many
SELECT id, content, title, user_id, tags, created_at, updated_at
FROM posts
ORDER BY created_at DESC;

-- name: UpdatePost :one
UPDATE posts
SET content = $2, title = $3, tags = $4, updated_at = now()
WHERE id = $1
RETURNING id, content, title, user_id, tags, created_at, updated_at;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;
