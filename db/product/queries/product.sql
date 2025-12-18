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
  stock,
  sku,
  meta_title,
  meta_description,
  category_id,
  is_active,
  weight,
  long,
  wide,
  high
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
  p.stock,
  p.sku,
  p.weight,
  p.long,
  p.wide,
  p.high,
  p.meta_title,
  p.meta_description,
  p.category_id,
  p.is_active,
  (
    SELECT
      COALESCE(json_agg(pf.name), '[]'::json)
    FROM
      product_files pf
    WHERE
      pf.product_id = p.id
  ) as files,
  (
    SELECT
      COALESCE(
        json_agg(
          json_build_object(
            'id',
            o.id,
            'name',
            o.name,
            'values',
            (
              SELECT
                COALESCE(
                  json_agg(json_build_object('id', ov.name, 'name', ov.name)),
                  '[]'::json
                )
              FROM
                option_values ov
              WHERE
                ov.option_id = o.id
            )
          )
        ),
        '[]'::json
      )
    FROM
      options o
    WHERE
      o.product_id = p.id
  ) as options,
  (
    SELECT
      COALESCE(
        json_agg(
          json_build_object(
            'id',
            v.id,
            'sku',
            v.sku,
            'origin_price',
            v.origin_price,
            'sale_price',
            v.sale_price,
            'file',
            v.file,
            'options',
            (
              SELECT
                COALESCE(jsonb_object_agg(o.name, ov.name), '{}'::jsonb)
              FROM
                variant_options vo
                JOIN options o ON o.id = vo.option_id
                JOIN option_values ov ON ov.id = vo.option_value_id
              WHERE
                vo.variant_id = v.id
            )
          )
          ORDER BY
            v.sale_price ASC
        ),
        '[]'::json
      )
    FROM
      variants v
    WHERE
      v.product_id = p.id
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
    stock,
    sku,
    meta_title,
    meta_description,
    category_id,
    weight,
    long,
    wide,
    high
  )
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING
  id;

-- name: UpdateProduct :exec
UPDATE products
SET
  name = $2,
  slug = $3,
  origin_price = $4,
  sale_price = $5,
  stock = $6,
  sku = $7,
  meta_title = $8,
  meta_description = $9,
  category_id = $10,
  weight = $11,
  long = $12,
  wide = $13,
  high = $14
WHERE
  id = $1;

-- name: BulkDeleteProducts :exec
DELETE FROM products
WHERE
  id = ANY ($1::bigint[]);

-- name: CountProductsByCategory :one
SELECT
  COUNT(*)
FROM
  products
WHERE
  category_id = $1;

-- name: GetProductsByCategory :many
SELECT
  p.id,
  p.name,
  p.origin_price,
  p.slug,
  p.sale_price,
  (
    SELECT
      COALESCE(json_agg(pf.name))
    FROM
      product_files pf
    WHERE
      pf.product_id = p.id
  ) as files
FROM
  products p
WHERE
  category_id = $3
LIMIT
  $1
OFFSET
  $2;

-- name: SearchProducts :many
SELECT
    p.id,
    p.name,
    p.slug,
    p.origin_price,
    p.sale_price,
    (
        SELECT pf.name
        FROM product_files pf
        WHERE pf.product_id = p.id
          AND pf.is_primary = TRUE
        ORDER BY pf.no ASC
        LIMIT 1
    ) AS file
FROM products p
WHERE ($1 = '' OR p.name ILIKE '%' || $1 || '%')
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountSearchProducts :one
SELECT COUNT(*) AS total
FROM products
WHERE ($1 = '' OR name LIKE '%' || $1 || '%');
