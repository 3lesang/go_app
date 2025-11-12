-- name: CountReviewsByProduct :one
SELECT
  COUNT(*)
FROM
  reviews
WHERE product_id = $1;

-- name: GetAverageRatingByProduct :one
SELECT
    AVG(rating) AS average_rating,
    COUNT(*) AS total_reviews
FROM reviews
WHERE product_id = $1;

-- name: GetReviewsByProductID :many
SELECT
  r.id,
  r.rating,
  r.comment,
  (
    SELECT COALESCE(json_agg(
      rf.name
    ))
    FROM review_files rf
    WHERE rf.review_id = r.id
  ) AS files,
  (
    SELECT COALESCE(json_build_object(
      'id', c.id,
      'name' ,c.name,
      'avatar', c.avatar
    ), '{}'::json)
    FROM customers c
    WHERE c.id = r.customer_id
  ) AS customer
FROM reviews r
WHERE
  r.product_id = sqlc.arg(product_id)
  AND (
    sqlc.narg(rating)::int IS NULL
    OR r.rating = sqlc.narg(rating)
  )
  AND (
    sqlc.narg(has_file)::boolean IS NULL
    OR (r.has_file = sqlc.narg(has_file))
  )
ORDER BY
  CASE WHEN sqlc.narg(sort_flag)::int = 1 THEN r.created_at END DESC,
  CASE WHEN sqlc.narg(sort_flag)::int = 0 THEN r.created_at END ASC
LIMIT sqlc.arg(limit_count)
OFFSET sqlc.arg(offset_count);

-- name: CreateReview :one
INSERT INTO
  reviews (rating, comment, has_file, product_id, customer_id)
VALUES
  ($1, $2, $3, $4, $5)
RETURNING id;
