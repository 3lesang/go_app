-- name: GetDiscountEffects :many
SELECT * FROM discount_effects WHERE discount_id = $1;

-- name: CreateDiscountEffect :one
INSERT INTO discount_effects (
  discount_id, effect_type, value, applies_to
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateDiscountEffect :exec
UPDATE discount_effects
SET
  effect_type = $2,
  value = $3,
  applies_to = $4
WHERE id = $1;

-- name: BulkDeleteDiscountEffects :exec
DELETE FROM discount_effects
WHERE id = ANY($1::bigint[]);
