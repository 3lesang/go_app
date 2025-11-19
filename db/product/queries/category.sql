-- name: CountCategories :one
SELECT
  COUNT(*)
FROM
  categories;

-- name: GetCategories :many
SELECT
  id,
  name,
  slug
FROM
  categories
ORDER BY
  id
LIMIT
  $1
OFFSET
  $2;

-- name: GetCategory :one
SELECT
  id,
  name,
  slug
FROM
  categories
WHERE
  id = $1
LIMIT
  1;

-- name: CreateCategory :exec
INSERT INTO
  categories (name, slug)
VALUES
  ($1, $2);

-- name: UpdateCategory :exec
UPDATE categories
SET
  name = $2,
  slug = $3
WHERE
  id = $1;

-- name: BulkDeleteCategories :exec
DELETE FROM categories
WHERE
  id = ANY ($1::bigint[]);
