-- name: GetCollections :many
SELECT
  id,
  name,
  file,
  slug
FROM
  collections
ORDER BY
  id;

-- name: GetCollection :one
SELECT
  id,
  name,
  slug,
  file,
  meta_title,
  meta_description,
  layout
FROM
  collections
WHERE
  id = $1
LIMIT
  1;

-- name: CreateCollection :one
INSERT INTO
  collections (
    name,
    slug,
    meta_title,
    meta_description,
    file,
    layout
  )
VALUES
  ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: UpdateCollection :exec
UPDATE collections
SET
  name = $2,
  slug = $3,
  meta_title = $4,
  meta_description = $5,
  file = $6,
  layout = $7
WHERE
  id = $1;

-- name: BulkDeleteCollections :exec
DELETE FROM collections
WHERE
  id = ANY ($1::bigint[]);
