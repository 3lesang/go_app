-- name: CreateAddress :one
INSERT INTO
  addresses (full_name, phone, email, address_line)
VALUES
  ($1, $2, $3, $4)
RETURNING
  id;
