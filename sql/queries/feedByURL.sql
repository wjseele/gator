-- name: GetFeedByURL :one
SELECT id
FROM feeds
WHERE name = $1;
