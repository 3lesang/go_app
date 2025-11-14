-- name: CountMenus :one
SELECT
  COUNT(*)
FROM
  menus;

-- name: GetMenus :many
SELECT id, name, position
FROM menus
LIMIT
  $1
OFFSET
  $2;

-- name: GetMenu :one
SELECT id, name, position
FROM menus
WHERE id = $1;

-- name: GetMenuByPosition :one
SELECT id, name, position
FROM menus
WHERE position = $1;

-- name: CreateMenu :one
INSERT INTO
  menus (name, position)
VALUES
  ($1, $2)
RETURNING id;

-- name: UpdateMenu :exec
UPDATE menus
SET
  name = $2,
  position = $3
WHERE
  id = $1;

-- name: BulkDeleteMenus :exec
DELETE FROM menus
WHERE
  id = ANY ($1::bigint[]);
