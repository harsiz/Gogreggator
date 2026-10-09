-- name: GetUserFromUserId :one
SELECT
    name
FROM
    users
WHERE
    id = $1
LIMIT 1;