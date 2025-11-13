-- name: CountPages :one
SELECT
  COUNT(*)
FROM
  pages;

-- name: GetPages :many
SELECT id, name, slug
FROM pages
LIMIT
  $1
OFFSET
  $2;

-- name: GetPage :one
SELECT id, name, slug
FROM pages
WHERE id = $1;

-- name: CreatePage :one
INSERT INTO
  pages (name, slug)
VALUES
  ($1, $2)
RETURNING id;

-- name: UpdatePage :exec
UPDATE pages
SET
  name = $2,
  slug = $3
WHERE
  id = $1;

-- name: BulkDeletePages :exec
DELETE FROM pages
WHERE
  id = ANY ($1::bigint[]);
