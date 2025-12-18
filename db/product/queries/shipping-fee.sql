-- name: CountShippingFees :one
SELECT
  COUNT(*)
FROM
  shipping_fees;

-- name: GetShippingFees :many
SELECT
  id,
  min_weight,
  max_weight,
  fee_amount,
  min_order_value,
  free_shipping
FROM
  shipping_fees
ORDER BY
  id
LIMIT
  $1
OFFSET
  $2;

-- name: GetShippingFee :one
SELECT
  id,
  min_weight,
  max_weight,
  fee_amount,
  min_order_value,
  free_shipping
FROM
  shipping_fees
WHERE
  id = $1
LIMIT
  1;

-- name: GetShippingFeeByWeight :one
SELECT
  id,
  min_weight,
  max_weight,
  fee_amount,
  min_order_value,
  free_shipping
FROM
  shipping_fees
WHERE
  $1 >= min_weight
  AND $1 < max_weight
ORDER BY
  max_weight ASC
LIMIT 1;

-- name: CreateShippingFee :one
INSERT INTO
  shipping_fees (
  min_weight,
  max_weight,
  fee_amount,
  min_order_value,
  free_shipping
  )
VALUES
  ($1, $2, $3, $4, $5)
RETURNING
  id;

-- name: UpdateShippingFee :exec
UPDATE shipping_fees
SET
  min_weight = $2,
  max_weight = $3,
  fee_amount = $4,
  min_order_value = $5,
  free_shipping = $6
WHERE
  id = $1;

-- name: BulkDeleteShippingFees :exec
DELETE FROM shipping_fees
WHERE
  id = ANY ($1::bigint[]);
