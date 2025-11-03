-- name: GetItemsByOrderID :many
SELECT
  oi.id AS order_item_id,
  oi.quantity,
  oi.sale_price,
  p.id AS product_id,
  p.name,
  p.slug,
  jsonb_object_agg(o.name, ov.name) AS options
FROM order_items oi
JOIN products p ON oi.product_id = p.id
JOIN variants v ON oi.variant_id = v.id
JOIN variant_options vo ON vo.variant_id = v.id
JOIN options o ON o.id = vo.option_id
JOIN option_values ov ON ov.id = vo.option_value_id
WHERE oi.order_id = $1
GROUP BY oi.id, p.id, p.name, p.slug, oi.quantity, oi.sale_price
ORDER BY oi.id;

-- name: BulkInsertOrderItems :exec
INSERT INTO order_items (quantity, sale_price, order_id, product_id, variant_id)
SELECT UNNEST(@quantities::int[]),
			UNNEST(@sale_prices::int[]),
			UNNEST(@order_ids::bigint[]),
			UNNEST(@product_ids::bigint[]),
			UNNEST(@variant_ids::bigint[]);