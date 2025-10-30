-- name: GetVariantsByProductID :many
SELECT
  id,
  origin_price,
  sale_price,
  file,
  stock,
  sku
FROM
  variants
WHERE
  product_id = $1
ORDER BY no ASC;

-- name: CreateVariant :exec
INSERT INTO
  variants (
    origin_price,
    sale_price,
    file,
    stock,
    sku,
    no,
    product_id
  )
VALUES
  ($1, $2, $3, $4, $5, $6, $7);

-- name: BulkInsertVariants :many
INSERT INTO
  variants (
    origin_price,
    sale_price,
    file,
    stock,
    sku,
    no,
    product_id
  )
SELECT
  unnest(@origin_prices::int[]),
  unnest(@sale_prices::int[]),
  unnest(@files::text[]),
  unnest(@stocks::int[]),
  unnest(@skus::text[]),
  unnest(@nos::int[]),
  unnest(@product_ids::bigint[])
RETURNING
  id;

-- name: UpdateVariant :exec
UPDATE variants
SET
  origin_price = $2,
  sale_price = $3,
  file = $4,
  stock = $5,
  sku = $6,
  product_id = $7
WHERE
  id = $1;

-- name: BulkUpdateVariants :exec
UPDATE variants AS v
SET
  origin_price = data.origin_price,
  sale_price   = data.sale_price,
  file         = data.file,
  stock        = data.stock,
  sku          = data.sku
FROM (
  SELECT
    UNNEST(@ids::bigint[])       AS id,
    UNNEST(@origin_prices::int[]) AS origin_price,
    UNNEST(@sale_prices::int[])   AS sale_price,
    UNNEST(@files::text[])        AS file,
    UNNEST(@stocks::int[])        AS stock,
    UNNEST(@skus::text[])         AS sku
) AS data
WHERE v.id = data.id
  AND (
    v.origin_price IS DISTINCT FROM data.origin_price OR
    v.sale_price   IS DISTINCT FROM data.sale_price OR
    v.file         IS DISTINCT FROM data.file OR
    v.stock        IS DISTINCT FROM data.stock OR
    v.sku          IS DISTINCT FROM data.sku
  );

-- name: DeleteVariantsByProductID :exec
DELETE FROM variants
WHERE
  product_id = $1;

-- name: DeleteVariantsNotInIDsByProductID :exec
DELETE FROM variants
WHERE product_id = @product_id
  AND id NOT IN (SELECT UNNEST(@ids::bigint[]));
