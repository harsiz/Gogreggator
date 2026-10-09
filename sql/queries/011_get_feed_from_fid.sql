-- name: GetFeedFromFeedId :one
SELECT
    *
FROM
    feeds
WHERE
    id = $1
LIMIT 1;