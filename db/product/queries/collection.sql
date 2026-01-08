-- name: CountCollections :one
SELECT
  COUNT(*)
FROM
  collections;

-- name: GetCollections :many
SELECT
  id,
  name,
  file,
  slug,
  created_at
FROM
  collections
ORDER BY
  created_at DESC
LIMIT $1
OFFSET $2;

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

-- name: GetCollectionsByLayout :many
SELECT
  id,
  file,
  slug
FROM
  collections
WHERE
  layout = $1;

-- name: GetCollectionsBySlug :one
SELECT
  id,
  name,
  file,
  slug,
  meta_title,
  meta_description
FROM
  collections
WHERE
  slug = $1
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
RETURNING
  id;

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
