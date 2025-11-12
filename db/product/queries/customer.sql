-- name: CountCustomers :one
SELECT
  COUNT(*)
FROM
  customers;

-- name: GetCustomers :many
SELECT id, name
FROM customers
LIMIT
  $1
OFFSET
  $2;
