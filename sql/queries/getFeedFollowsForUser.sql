-- name: GetFeedFollowsForUser :many
SELECT feeds.name AS feed_name, users.name AS user_name
FROM feed_follows
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
INNER JOIN users ON $1 = feed_follows.user_id
WHERE users.id = $1;
