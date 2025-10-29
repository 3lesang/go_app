-- name: BulkInsertProductCollection :exec
INSERT INTO
  product_collections (product_id, collection_id)
SELECT
  unnest(@product_ids::bigint[]),
  unnest(@collection_ids::bigint[]);

-- name: GetCollectionsByProductID :many
SELECT
  c.id,
  c.name
FROM
  product_collections pc
  LEFT JOIN collections c ON pc.collection_id = c.id
WHERE
  product_id = $1;

-- name: DeleteProductCollections :exec
DELETE FROM product_collections WHERE product_id = $1;