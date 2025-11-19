-- name: BulkInsertProductFiles :exec
INSERT INTO
  product_files (name, is_primary, no, product_id)
SELECT
  unnest(@names::text[]),
  unnest(@is_primaries::bool[]),
  unnest(@nos::int[]),
  unnest(@product_ids::bigint[]);

-- name: GetFilesByProductID :many
SELECT
  name
FROM
  product_files
WHERE
  product_id = $1;

-- name: DeleteProductFiles :exec
DELETE FROM product_files
WHERE
  product_id = $1;
