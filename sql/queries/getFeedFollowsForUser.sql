-- name: GetFeedFollowsForUser :many
SELECT feeds.name, users.name
FROM feed_follows
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
INNER JOIN users ON users.id = feed_follows.user_id;
