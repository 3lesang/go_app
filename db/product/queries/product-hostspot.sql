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