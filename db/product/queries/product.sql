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

-- name: GetProductBySlug :one
SELECT
  p.id,
  p.name,
  p.slug,
  p.origin_price,
  p.sale_price,
  p.meta_title,
  p.meta_description,
  p.category_id,
  p.is_active,
  (
    SELECT COALESCE(json_agg(pf.name), '[]'::json)
    FROM product_files pf
    WHERE pf.product_id = p.id
  ) as files,
  (
    SELECT COALESCE(json_agg(json_build_object(
      'id', o.id,
      'name', o.name,
      'values', (
        SELECT COALESCE(json_agg(json_build_object(
            'id', ov.name,
            'name', ov.name
          )), '[]'::json)
        FROM option_values ov
        WHERE ov.option_id = o.id
        )
      )), '[]'::json)
    FROM options o
    WHERE o.product_id = p.id
  ) as options,
  (
    SELECT COALESCE(
      json_agg(
        json_build_object(
          'id', v.id,
          'sku', v.sku,
          'origin_price', v.origin_price,
          'sale_price', v.sale_price,
          'file', v.file,
          'options', (
            SELECT COALESCE(jsonb_object_agg(o.name, ov.name), '{}'::jsonb)
            FROM variant_options vo
            JOIN options o ON o.id = vo.option_id
            JOIN option_values ov ON ov.id = vo.option_value_id
            WHERE vo.variant_id = v.id
          )
        )
        ORDER BY v.sale_price ASC
      ),
      '[]'::json
    )
    FROM variants v
    WHERE v.product_id = p.id
  ) AS variants
FROM
  products p
WHERE
  slug = $1
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
