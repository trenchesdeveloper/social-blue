-- name: CreatePost :one
INSERT INTO posts (content, title, user_id, tags)
VALUES ($1, $2, $3, $4)
RETURNING id, content, title, user_id, tags, version, created_at, updated_at;

-- name: GetPostByID :one
SELECT id, content, title, user_id, tags, version, created_at, updated_at
FROM posts
WHERE id = $1;

-- name: ListPosts :many
SELECT id, content, title, user_id, tags,version, created_at, updated_at
FROM posts
ORDER BY created_at DESC;

-- name: UpdatePost :one
UPDATE posts
SET
     content = COALESCE(NULLIF($2, ''), content),
    title = COALESCE(NULLIF($3, ''), title),
     tags = COALESCE(NULLIF($4::text[], '{}'), tags),
        version = version + 1,
    updated_at = now()
 WHERE id = $1 AND version = $5
    RETURNING id, content, title, user_id, tags, version, created_at, updated_at;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;
