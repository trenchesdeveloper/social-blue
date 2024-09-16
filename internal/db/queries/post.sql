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


-- name: GetUserFeed :many
SELECT p.id, p.content, p.title, p.user_id, p.tags, p.version, p.created_at, p.updated_at, u.username,
       COUNT(c.id) AS comments_count
FROM posts p
         LEFT JOIN comments c ON c.post_id = p.id
         LEFT JOIN users u ON p.user_id = u.id
         JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
WHERE f.user_id = $1 AND
    (p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') AND
    (p.tags @> $5 OR $5 = '{}')
GROUP BY p.id, u.username
ORDER BY p.created_at DESC
LIMIT $2
OFFSET $3;