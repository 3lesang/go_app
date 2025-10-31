-- name: BulkInsertProductCollection :exec
INSERT INTO product_collections (product_id, collection_id)
SELECT p.id, c.id
FROM unnest(@product_ids::bigint[]) AS p(id),
     unnest(@collection_ids::bigint[]) AS c(id)
ON CONFLICT (product_id, collection_id) DO NOTHING;

-- name: GetCollectionsByProductID :many
SELECT
  c.id,
  c.name
FROM
  product_collections pc
  LEFT JOIN collections c ON pc.collection_id = c.id
WHERE
  product_id = $1;

-- name: GetProductsByCollectionID :many
SELECT 
  p.id,
  p.name
FROM
  product_collections pc
  LEFT JOIN products p ON pc.product_id = p.id
WHERE collection_id = $1;

-- name: DeleteCollectionsByProductID :exec
DELETE FROM product_collections WHERE product_id = $1;

-- name: DeleteProductsByCollectionID :exec
DELETE FROM product_collections WHERE collection_id = $1;

-- name: DeleteCollectionProductsNotInIDsByCollection :exec
DELETE FROM product_collections
WHERE collection_id IN (SELECT UNNEST(@collection_ids::bigint[])) 
  AND product_id NOT IN (SELECT UNNEST(@product_ids::bigint[]));
