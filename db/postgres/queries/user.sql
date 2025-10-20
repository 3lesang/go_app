-- name: ListUsers :many
SELECT id, name, phone, email, username
FROM users
ORDER BY id;

-- name: GetUser :one
SELECT id, name, username
FROM users
WHERE id = $1 LIMIT 1;

-- name: CreateUser :exec
INSERT INTO users (
    name, phone, email, password, username
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: DeleteUsers :exec
DELETE FROM users
WHERE id IN (sqlc.slice(ids));

-- name: GetUserByIdentify :one
SELECT id, name, phone, email, username, password
FROM users
WHERE username = $1 OR email = $2 OR phone = $3
LIMIT 1;
