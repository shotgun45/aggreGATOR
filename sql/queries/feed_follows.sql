-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT
    iff.id,
    iff.created_at,
    iff.updated_at,
    iff.user_id,
    iff.feed_id,
    f.name AS feed_name,
    u.name AS user_name
FROM inserted_feed_follow iff
INNER JOIN feeds f ON iff.feed_id = f.id
INNER JOIN users u ON iff.user_id = u.id;

-- name: GetFeedFollowWithNames :one
SELECT ff.id, ff.created_at, ff.updated_at, ff.user_id, ff.feed_id, u.name AS user_name, f.name AS feed_name
FROM feed_follows ff
JOIN users u ON ff.user_id = u.id
JOIN feeds f ON ff.feed_id = f.id
WHERE ff.id = $1;

-- name: GetFeedFollowsForUser :many
SELECT
    ff.id,
    ff.created_at,
    ff.updated_at,
    ff.user_id,
    ff.feed_id,
    f.name AS feed_name,
    u.name AS user_name
FROM feed_follows ff
INNER JOIN feeds f ON ff.feed_id = f.id
INNER JOIN users u ON ff.user_id = u.id
WHERE ff.user_id = $1
ORDER BY ff.created_at DESC;
