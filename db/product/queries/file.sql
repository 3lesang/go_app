-- name: CountFiles :one
SELECT
  COUNT(*)
FROM
  files;

-- name: GetFiles :many
SELECT
  id,
  name
FROM
  files
ORDER BY
  created_at DESC
LIMIT
  $1
OFFSET
  $2;

-- name: BulkInsertFiles :exec
INSERT INTO
  files (name)
SELECT
  unnest(@names::text[]) as name;

-- name: BulkDeleteFiles :exec
DELETE FROM files
WHERE
  id = ANY ($1::bigint[]);
