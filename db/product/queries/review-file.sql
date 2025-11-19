-- name: BulkInsertReviewFiles :exec
INSERT INTO
  review_files (name, review_id)
SELECT
  unnest(@names::text[]),
  unnest(@review_ids::bigint[]);

-- name: CountReviewFilesByProduct :one
SELECT
  COUNT(*)
FROM
  review_files rf
  JOIN reviews r ON r.id = rf.review_id
WHERE
  r.product_id = $1;

-- name: GetReviewFilesByProduct :many
SELECT
  rf.name
FROM
  review_files rf
  JOIN reviews r ON r.id = rf.review_id
WHERE
  r.product_id = $1
LIMIT
  8
OFFSET
  0;
