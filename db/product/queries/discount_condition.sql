-- name: GetDiscountConditions :many
SELECT * FROM discount_conditions WHERE discount_id = $1;

-- name: CreateDiscountCondition :one
INSERT INTO discount_conditions (
  discount_id, condition_type, operator, value
)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateDiscountCondition :exec
UPDATE discount_conditions
SET
  condition_type = $2,
  value = $3
WHERE id = $1;

-- name: BulkDeleteDiscountConditions :exec
DELETE FROM discount_conditions
WHERE id = ANY($1::bigint[]);
