-- name: CreateAddress :one
INSERT INTO
  addresses (full_name, phone, address_line)
VALUES
  ($1, $2, $3)
RETURNING
  id;
