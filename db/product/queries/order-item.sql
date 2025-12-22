-- name: BulkInsertOrderItems :exec
INSERT INTO
  order_items (
    quantity,
    sale_price,
    order_id,
    product_id,
    variant_id
  )
SELECT
  UNNEST(@quantities::int[]),
  UNNEST(@sale_prices::int[]),
  UNNEST(@order_ids::bigint[]),
  UNNEST(@product_ids::bigint[]),
  UNNEST(@variant_ids::bigint[]);
