-- name: GetDiscountTargets :many
SELECT * FROM discount_targets WHERE discount_id = $1;

-- name: CreateDiscountTarget :one
INSERT INTO discount_targets (
  discount_id, target_type, target_id
)
VALUES ($1, $2, $3)
RETURNING *;

-- name: BulkInsertDiscountTargets :exec
INSERT INTO discount_targets (discount_id, target_type, target_id)
SELECT
unnest(@discount_ids::bigint[]),
unnest(@target_types::text[]),
unnest(@target_ids::bigint[]);

-- name: UpdateDiscountTarget :exec
UPDATE discount_targets
SET
  target_type = $2,
  target_id = $3
WHERE id = $1;

-- name: BulkDeleteDiscountTargets :exec
DELETE FROM discount_targets
WHERE id = ANY($1::bigint[]);
