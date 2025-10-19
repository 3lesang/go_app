-- name: CreateUser :exec
INSERT INTO users (name, phone, email, password, username) VALUES (?, ?, ?, ?, ?);

-- name: GetUser :one
SELECT id, name, phone, email, username FROM users WHERE id = ?;

-- name: ListUsers :many
SELECT id, name, phone, email, username FROM users ORDER BY id;

-- name: DeleteUsers :exec
DELETE FROM users WHERE id IN (sqlc.slice(ids));

-- name: GetUserByIdentify :one
SELECT id, name, phone, email, username, password FROM users WHERE username = ? OR email = ? OR phone = ? LIMIT 1;
