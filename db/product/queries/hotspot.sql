-- name: CountHotspots :one
SELECT
  COUNT(*)
FROM
  hotspots;

-- name: GetHotspots :many
SELECT
  id,
  file
FROM
  hotspots
ORDER BY
  id
LIMIT
  $1
OFFSET
  $2;

-- name: GetHotspot :one
SELECT
  h.id,
  h.file,
  COALESCE(
    json_agg(
      json_build_object(
        'id', ph.id,
        'product_id', ph.product_id,
        'x', ph.x,
        'y', ph.y,
        'product', json_build_object(
          'id', p.id,
          'name', p.name
        )
      )
    ) FILTER (WHERE ph.id IS NOT NULL),
    '[]'
  ) AS spots
FROM hotspots h
LEFT JOIN product_hotspots ph ON ph.hotspot_id = h.id
LEFT JOIN products p ON p.id = ph.product_id
WHERE h.id = $1
GROUP BY h.id;

-- name: CreateHotspot :one
INSERT INTO hotspots (file)
VALUES ($1)
RETURNING id;

-- name: BulkDeleteHotspots :exec
DELETE FROM hotspots
WHERE id = ANY($1::bigint[]);
