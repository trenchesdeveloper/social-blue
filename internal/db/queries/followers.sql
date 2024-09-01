-- name: FollowUser :exec
INSERT INTO followers (user_id, follower_id)
VALUES ($1, $2);

-- name: UnFollowUser :exec
DELETE FROM followers
WHERE user_id = $1 AND follower_id = $2;