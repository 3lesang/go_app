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

-- name: GetHotspotsByProductId :many
SELECT 
  h.id,
  h.file
FROM hotspots h
INNER JOIN product_hotspots ph ON h.id = ph.hotspot_id
WHERE ph.product_id = $1;

-- name: GetProductsByHotspotId :many
SELECT 
  p.id,
  p.name,
  f.name as file,
  p.origin_price,
  p.sale_price,
  p.slug,
  ph.x,
  ph.y
FROM product_hotspots ph
INNER JOIN products p ON p.id = ph.product_id
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
WHERE ph.hotspot_id = $1;