-- name: CountCustomers :one
SELECT
  COUNT(*)
FROM
  customers;

-- name: GetCustomers :many
SELECT
  id,
  name,
  phone
FROM
  customers
LIMIT
  $1
OFFSET
  $2;

-- name: GetCustomer :one
SELECT id, name, phone, email
FROM customers
WHERE id = $1
LIMIT 1;

-- name: CreateCustomer :one
INSERT INTO
  customers (name, phone, password)
VALUES
  ($1, $2, $3)
RETURNING
  id;

-- name: UpdateCustomer :one
UPDATE customers
SET
    name     = COALESCE($2, name),
    email    = COALESCE($3, email),
    password = CASE WHEN $4 = '' THEN password ELSE $4 END
WHERE id = $1
RETURNING id;

-- name: BulkDeleteCustomers :exec
DELETE FROM customers
WHERE
  id = ANY ($1::bigint[]);

-- name: GetCustomerByPhone :one
SELECT id, name, phone, password
FROM customers
WHERE phone = $1;