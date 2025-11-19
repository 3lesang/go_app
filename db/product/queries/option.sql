-- name: BulkInsertOptions :many
INSERT INTO
  options (name, no, product_id)
SELECT
  unnest(@names::text[]),
  unnest(@nos::int[]),
  unnest(@product_ids::bigint[])
RETURNING
  id,
  name;

-- name: GetOptionsByProductID :many
SELECT
  id,
  name,
  no
FROM
  options
WHERE
  product_id = $1
ORDER BY
  no ASC;

-- name: DeleteOptionsByProductID :exec
DELETE FROM options
WHERE
  product_id = $1;

-- name: BulkUpdateOptions :exec
UPDATE options AS o
SET
  name = data.name
FROM
  (
    SELECT
      UNNEST(@ids::bigint[]) AS id,
      UNNEST(@names::text[]) AS name
  ) AS data
WHERE
  o.id = data.id
  AND (o.name IS DISTINCT FROM data.name);

-- name: DeleteOptionsNotInIDs :exec
DELETE FROM options
WHERE
  product_id = @product_id
  AND id NOT IN (
    SELECT
      UNNEST(@ids::bigint[])
  );
