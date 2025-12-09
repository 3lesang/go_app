-- name: BulkInsertProductHotspots :exec
INSERT INTO
  product_hotspots (x, y, product_id, hotspot_id)
SELECT
  unnest(@xs::real[]) as x,
  unnest(@ys::real[]) as y,
  unnest(@product_ids::bigint[]) as product_id,
  unnest(@hotspot_ids::bigint[]) as hotspot_id;

-- name: BulkUpdateProductHotspots :exec
UPDATE product_hotspots ph
SET product_id = u.product_id
FROM (
    SELECT
        unnest(@hotspot_ids::bigint[]) AS hotspot_id,
        unnest(@product_ids::bigint[]) AS product_id
) AS u
WHERE ph.hotspot_id = u.hotspot_id;

-- name: BulkDeleteProductHotspots :exec
DELETE FROM product_hotspots
WHERE id = ANY($1::bigint[]);

-- name: GetProductHotspotsByHotspot :many
SELECT id
FROM product_hotspots
WHERE hotspot_id = $1;

-- name: GetHotspotByProduct :many
SELECT h.id, h.file,
  COALESCE(
    json_agg(
      json_build_object(
        'id', ph.id,
        'product_id', ph.product_id,
        'x', ph.x,
        'y', ph.y,
        'product', json_build_object(
          'id', p.id,
          'name', p.name,
          'slug', p.slug,
          'origin_price', p.origin_price,
          'sale_price', p.sale_price,
          'file', f.name
        )
      )
    ) FILTER (WHERE ph.id IS NOT NULL),
    '[]'
  ) AS spots
FROM hotspots h
JOIN product_hotspots ph ON ph.hotspot_id = h.id
LEFT JOIN products p ON p.id = ph.product_id
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
WHERE ph.product_id = $1
GROUP BY h.id, h.file;