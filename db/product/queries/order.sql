-- name: CountOrders :one
SELECT
  COUNT(*)
FROM
  orders;

-- name: GetOrders :many
SELECT
  id,
  total_amount,
  discount_amount,
  created_at
FROM
  orders
LIMIT
  $1
OFFSET
  $2;

-- name: GetOrder :one
SELECT
  o.total_amount,
  o.discount_amount,
  a.full_name,
  a.phone,
  a.address_line
FROM
  orders o
  LEFT JOIN addresses a ON o.shipping_address_id = a.id
WHERE
  o.id = $1;

-- name: CheckOrderCreated :one
SELECT
  id
FROM
  orders
WHERE
  id = $1;

-- name: CreateOrder :one
INSERT INTO
  orders (
    total_amount,
    discount_amount,
    shipping_address_id
  )
VALUES
  ($1, $2, $3)
RETURNING
  id;

-- name: BulkDeleteOrders :exec
DELETE FROM orders
WHERE
  id = ANY ($1::bigint[]);
