-- name: GetCollections :many
SELECT
  id,
  name,
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
  image_url,
  meta_title,
  meta_description,
  layout
FROM
  collections
WHERE
  id = $1
LIMIT
  1;

-- name: CreateCollection :exec
INSERT INTO
  collections (
    name,
    slug,
    meta_title,
    meta_description,
    image_url,
    layout
  )
VALUES
  ($1, $2, $3, $4, $5, $6);

-- name: UpdateCollection :exec
UPDATE collections
SET
  name = $2,
  slug = $3,
  meta_title = $4,
  meta_description = $5,
  image_url = $6,
  layout = $7
WHERE
  id = $1;

-- name: BulkDeleteCollections :exec
DELETE FROM collections
WHERE
  id = ANY ($1::bigint[]);
