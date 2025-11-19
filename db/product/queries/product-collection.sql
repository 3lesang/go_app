-- name: BulkInsertProductCollection :exec
INSERT INTO
  product_collections (product_id, collection_id)
SELECT
  p.id,
  c.id
FROM
  unnest(@product_ids::bigint[]) AS p (id),
  unnest(@collection_ids::bigint[]) AS c (id)
ON CONFLICT (product_id, collection_id) DO NOTHING;

-- name: GetCollectionsByProductID :many
SELECT
  c.id,
  c.name
FROM
  product_collections pc
  LEFT JOIN collections c ON pc.collection_id = c.id
WHERE
  product_id = $1;

-- name: GetProductsByCollectionID :many
SELECT
  p.id,
  p.name
FROM
  product_collections pc
  LEFT JOIN products p ON pc.product_id = p.id
WHERE
  collection_id = $1;

-- name: GetProductsByCollection :many
SELECT
  p.id,
  p.name,
  p.slug,
  p.origin_price,
  p.sale_price,
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
            v.id,
            'sku',
            v.sku,
            'origin_price',
            v.origin_price,
            'sale_price',
            v.sale_price,
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
        ),
        '[]'::json
      )
    FROM
      variants v
    WHERE
      v.product_id = p.id
  ) as variants
FROM
  product_collections pc
  LEFT JOIN products p ON pc.product_id = p.id
WHERE
  collection_id = $1;

-- name: GetHomeCollectionsWithProductsAndVariants :many
SELECT
  c.id,
  c.name,
  c.file,
  c.slug,
  (
    SELECT
      COALESCE(
        json_agg(
          json_build_object(
            'id',
            p.id,
            'name',
            p.name,
            'slug',
            p.slug,
            'sale_price',
            p.sale_price,
            'origin_price',
            p.origin_price,
            'files',
            (
              SELECT
                COALESCE(json_agg(pf.name), '[]'::json)
              FROM
                product_files pf
              WHERE
                pf.product_id = p.id
            ),
            'variants',
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
                  ),
                  '[]'::json
                )
              FROM
                variants v
              WHERE
                v.product_id = p.id
            )
          )
        ),
        '[]'::json
      )
    FROM
      products p
      JOIN product_collections pc ON pc.product_id = p.id
    WHERE
      pc.collection_id = c.id
    LIMIT
      $1
    OFFSET
      $2
  ) AS products
FROM
  collections c
WHERE
  c.layout = 'home'
ORDER BY
  c.id;

-- name: DeleteCollectionsByProductID :exec
DELETE FROM product_collections
WHERE
  product_id = $1;

-- name: DeleteProductsByCollectionID :exec
DELETE FROM product_collections
WHERE
  collection_id = $1;

-- name: DeleteCollectionProductsNotInIDsByCollection :exec
DELETE FROM product_collections
WHERE
  collection_id IN (
    SELECT
      UNNEST(@collection_ids::bigint[])
  )
  AND product_id NOT IN (
    SELECT
      UNNEST(@product_ids::bigint[])
  );
