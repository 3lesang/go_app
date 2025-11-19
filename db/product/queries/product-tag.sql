-- name: BulkInsertProductTags :exec
INSERT INTO
  product_tags (name, product_id)
SELECT
  unnest(@names::text[]),
  unnest(@product_ids::bigint[]);

-- name: GetTagsByProductID :many
SELECT
  name
FROM
  product_tags
WHERE
  product_id = $1;

-- name: DeleteProductTags :exec
DELETE FROM product_tags
WHERE
  product_id = $1;
