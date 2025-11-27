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

-- name: CreateCustomer :one
INSERT INTO
  customers (name, phone, password)
VALUES
  ($1, $2, $3)
RETURNING
  id;

-- name: BulkDeleteCustomers :exec
DELETE FROM customers
WHERE
  id = ANY ($1::bigint[]);

-- name: GetCustomerByPhone :one
SELECT id, name, phone, password
FROM customers
WHERE phone = $1;