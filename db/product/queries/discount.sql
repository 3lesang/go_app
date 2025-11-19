-- name: ListDiscounts :many
SELECT *
FROM discounts
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: CountDiscounts :one
SELECT COUNT(*)
FROM discounts;

-- name: CreateDiscount :one
INSERT INTO discounts (
  title, code, discount_type, status, usage_limit, per_customer_limit, starts_at, ends_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;

-- name: GetDiscountByID :one
SELECT id, title, code, discount_type, status, usage_limit, per_customer_limit, starts_at, ends_at
FROM discounts
WHERE id = $1;

-- name: GetDiscountByCode :one
SELECT * FROM discounts
WHERE code = $1
  AND starts_at <= NOW()
  AND (ends_at IS NULL OR ends_at >= NOW())
  AND status = 'active'
  AND (usage_limit IS NULL OR usage_count < usage_limit);

-- name: UpdateDiscount :one
UPDATE discounts
SET
  title = $2,
  status = $3,
  usage_limit = $4,
  usage_count = $5,
  per_customer_limit = $6,
  starts_at = $7,
  ends_at = $8,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id;

-- name: IncrementDiscountUsage :exec
UPDATE discounts SET usage_count = usage_count + 1 WHERE id = $1;

-- name: BulkDeleteDiscounts :exec
DELETE FROM discounts
WHERE id = ANY($1::bigint[]);

-- name: GetDiscountWithRelations :one
SELECT
  d.*,
  COALESCE(
    (
      SELECT json_agg(dc)
      FROM discount_conditions dc
      WHERE dc.discount_id = d.id
    ),
    '[]'::json
  ) AS conditions,
  COALESCE(
    (
      SELECT json_agg(de)
      FROM discount_effects de
      WHERE de.discount_id = d.id
    ),
    '[]'::json
  ) AS effects,
  COALESCE(
    (
      SELECT json_agg(dt)
      FROM discount_targets dt
      WHERE dt.discount_id = d.id
    ),
    '[]'::json
  ) AS targets
FROM discounts d
WHERE d.id = $1
LIMIT 1;
