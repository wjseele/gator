-- name: GetFeedByURL :one
SELECT id
FROM feeds
WHERE url = $1;
