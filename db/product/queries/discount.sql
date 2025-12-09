-- name: ListDiscounts :many
SELECT id, title, description, code, discount_type, status, usage_limit, per_customer_limit, starts_at, ends_at
FROM discounts
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: CountDiscounts :one
SELECT COUNT(*)
FROM discounts;

-- name: CreateDiscount :one
INSERT INTO discounts (
  title, description, code, discount_type, status, usage_limit, per_customer_limit, starts_at, ends_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id;

-- name: GetDiscountByID :one
SELECT id, title, description, code, discount_type, status, usage_limit, per_customer_limit, starts_at, ends_at
FROM discounts
WHERE id = $1;

-- name: GetDiscountByCode :one
SELECT id, title, description, code, discount_type, status, usage_limit, per_customer_limit, starts_at, ends_at
FROM discounts
WHERE code = $1
  AND starts_at <= NOW()
  AND (ends_at IS NULL OR ends_at >= NOW())
  AND status = 'active'
  AND (usage_limit IS NULL OR usage_count < usage_limit);

-- name: UpdateDiscount :one
UPDATE discounts
SET
  title = $2,
  description = $3,
  status = $4,
  usage_limit = $5,
  per_customer_limit = $6,
  starts_at = $7,
  ends_at = $8,
  updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id;

-- name: UpsertDiscountUsageUsage :exec
INSERT INTO discounts (id, usage_count)
VALUES ($1, 1)
ON CONFLICT (id) DO UPDATE
SET usage_count = discounts.usage_count + 1;

-- name: BulkDeleteDiscounts :exec
DELETE FROM discounts
WHERE id = ANY($1::bigint[]);

-- name: GetDiscountWithRelations :one
SELECT
  d.id, d.title, d.description, d.code, d.discount_type, d.status, d.usage_limit, d.per_customer_limit, d.starts_at, d.ends_at,
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

-- name: GetValidDiscounts :many
SELECT
    d.id, d.title, d.description, d.code, d.discount_type, d.status, d.usage_limit, d.per_customer_limit, d.starts_at, d.ends_at,
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
    ) AS effects
FROM
    discounts d
WHERE
    d.status = 'active'
    AND d.starts_at <= NOW()
    AND (d.ends_at IS NULL OR d.ends_at > NOW())
    AND (d.usage_limit IS NULL OR d.usage_count < d.usage_limit);
