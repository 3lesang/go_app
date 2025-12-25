-- name: CountOrders :one
SELECT
  COUNT(*)
FROM
  orders;

-- name: CountStatusOrder :one
SELECT
  SUM(CASE WHEN status = 'pending' THEN 1 ELSE 0 END) AS pending_count,
  SUM(CASE WHEN status = 'confirmed' THEN 1 ELSE 0 END) AS confirmed_count,
  SUM(CASE WHEN status = 'shipping' THEN 1 ELSE 0 END) AS shipping_count,
  SUM(CASE WHEN status = 'shipped' THEN 1 ELSE 0 END) AS shipped_count,
  SUM(CASE WHEN status = 'cancelled' THEN 1 ELSE 0 END) AS cancelled_count
FROM orders
LIMIT 1;

-- name: GetOrders :many
SELECT
  o.id,
  o.code,
  o.total_amount,
  o.discount_amount,
  o.shipping_fee_amount,
  o.status,
  o.created_at,
  a.full_name,
  a.phone,
  a.address_line,
  a.email,
  (
    SELECT COALESCE(
      json_agg(
        json_build_object(
          'id', oi.id,
          'product_id', p.id,
          'product_sku', p.sku,
          'variant_sku', v.sku,
          'name', p.name,
          'quantity', oi.quantity,
          'sale_price', oi.sale_price,
          'options', (
            SELECT COALESCE(
              json_agg(
                json_build_object(
                  'option', o.name,
                  'value', ov.name
                )
              ),
              '[]'::json
            )
            FROM variant_options vo
            LEFT JOIN options o ON o.id = vo.option_id
            LEFT JOIN option_values ov ON ov.id = vo.option_value_id 
            WHERE vo.variant_id = oi.variant_id        
          )
        )
      ) FILTER (WHERE oi.id IS NOT NULL),
      '[]'::json
    )
    FROM order_items oi
    LEFT JOIN products p ON p.id = oi.product_id
    LEFT JOIN variants v ON v.id = oi.variant_id
    WHERE oi.order_id = o.id
  ) AS items
FROM
  orders o
  LEFT JOIN addresses a ON o.shipping_address_id = a.id
ORDER BY o.id DESC
LIMIT
  $1
OFFSET
  $2;

-- name: CountOrdersByStatus :one
SELECT
  COUNT(*)
FROM
  orders
WHERE status = $1;

-- name: GetOrdersByStatus :many
SELECT
  o.id,
  o.code,
  o.total_amount,
  o.discount_amount,
  o.shipping_fee_amount,
  o.status,
  o.created_at,
  a.full_name,
  a.phone,
  a.address_line,
  a.email,
  (
    SELECT COALESCE(
      json_agg(
        json_build_object(
          'id', oi.id,
          'product_id', p.id,
          'product_sku', p.sku,
          'variant_sku', v.sku,
          'name', p.name,
          'quantity', oi.quantity,
          'sale_price', oi.sale_price,
          'options', (
            SELECT COALESCE(
              json_agg(
                json_build_object(
                  'option', o.name,
                  'value', ov.name
                )
              ),
              '[]'::json
            )
            FROM variant_options vo
            LEFT JOIN options o ON o.id = vo.option_id
            LEFT JOIN option_values ov ON ov.id = vo.option_value_id 
            WHERE vo.variant_id = oi.variant_id        
          )
        )
      ) FILTER (WHERE oi.id IS NOT NULL),
      '[]'::json
    )
    FROM order_items oi
    LEFT JOIN products p ON p.id = oi.product_id
    LEFT JOIN variants v ON v.id = oi.variant_id
    WHERE oi.order_id = o.id
  ) AS items
FROM
  orders o
  LEFT JOIN addresses a ON o.shipping_address_id = a.id
WHERE o.status = $3
ORDER BY o.id DESC
LIMIT
  $1
OFFSET
  $2;

-- name: GetOrder :one
SELECT
  o.id,
  o.code,
  o.total_amount,
  o.discount_amount,
  o.shipping_fee_amount,
  o.status,
  o.created_at,  
  a.full_name,
  a.phone,
  a.address_line,
  a.email,
  (
    SELECT COALESCE(
      json_agg(
        json_build_object(
          'id', oi.id,
          'product_id', p.id,
          'product_sku', p.sku,
          'variant_sku', v.sku,
          'name', p.name,
          'quantity', oi.quantity,
          'sale_price', oi.sale_price,
          'options', (
            SELECT COALESCE(
              json_agg(
                json_build_object(
                  'option', o.name,
                  'value', ov.name
                )
              ),
              '[]'::json
            )
            FROM variant_options vo
            LEFT JOIN options o ON o.id = vo.option_id
            LEFT JOIN option_values ov ON ov.id = vo.option_value_id 
            WHERE vo.variant_id = oi.variant_id        
          )
        )
      ) FILTER (WHERE oi.id IS NOT NULL),
      '[]'::json
    )
    FROM order_items oi
    LEFT JOIN products p ON p.id = oi.product_id
    LEFT JOIN variants v ON v.id = oi.variant_id
    WHERE oi.order_id = o.id
  ) AS items
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
    code,
    total_amount,
    discount_amount,
    shipping_fee_amount,
    shipping_address_id
  )
VALUES
  ($1, $2, $3, $4, $5)
RETURNING
  id;

-- name: UpdateOrder :exec
UPDATE orders
SET status = $2,
    cancel_reason = $3
WHERE id = $1;

-- name: BulkDeleteOrders :exec
DELETE FROM orders
WHERE
  id = ANY ($1::bigint[]);
