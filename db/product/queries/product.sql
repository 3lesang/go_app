-- name: CountProducts :one
SELECT
  COUNT(*)
FROM
  products;

-- name: GetProducts :many
SELECT
  p.id,
  p.name,
  p.is_active,
  f.name AS file
FROM
  products p
  LEFT JOIN LATERAL (
    SELECT
      name
    FROM
      product_files
    WHERE
      product_id = p.id
      AND is_primary = true
    LIMIT
      1
  ) f ON true
ORDER BY
  created_at DESC
LIMIT
  $1
OFFSET
  $2;

-- name: GetProduct :one
SELECT
  id,
  name,
  slug,
  origin_price,
  sale_price,
  meta_title,
  meta_description,
  category_id,
  is_active
FROM
  products
WHERE
  id = $1
LIMIT
  1;

-- name: CreateProduct :one
INSERT INTO
  products (
    name,
    slug,
    origin_price,
    sale_price,
    meta_title,
    meta_description,
    category_id
  )
VALUES
  ($1, $2, $3, $4, $5, $6, $7)
RETURNING
  id;

-- name: UpdateProduct :exec
UPDATE products
SET
  name = $2,
  slug = $3,
  origin_price = $4,
  sale_price = $5,
  meta_title = $6,
  meta_description = $7,
  category_id = $8
WHERE
  id = $1;

-- name: BulkDeleteProducts :exec
DELETE FROM products
WHERE
  id = ANY ($1::bigint[]);
