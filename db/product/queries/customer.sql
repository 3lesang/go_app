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

-- name: GetCustomerByPhone :one
SELECT id, name, phone, password, phone_verified
FROM customers
WHERE phone = $1
LIMIT 1;

-- name: CreateCustomer :one
INSERT INTO
  customers (name, phone, password, zns_otp)
VALUES
  ($1, $2, $3, $4)
RETURNING
  id;

-- name: VerifyPhone :exec
UPDATE customers
SET phone_verified = true,
    zns_otp = NULL
WHERE phone = $1 AND zns_otp = $2;

-- name: UpdateCustomer :one
UPDATE customers
SET
    name     = COALESCE($2, name),
    email    = COALESCE($3, email),
    password = CASE WHEN $4 = '' THEN password ELSE $4 END,
    phone_verified = $5
WHERE id = $1
RETURNING id;

-- name: BulkDeleteCustomers :exec
DELETE FROM customers
WHERE
  id = ANY ($1::bigint[]);