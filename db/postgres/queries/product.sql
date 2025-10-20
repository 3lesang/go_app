-- name: ListProducts :many
SELECT id, name, slug
FROM products
ORDER BY id;

-- name: GetProduct :one
SELECT id, name
FROM products
WHERE id = $1
LIMIT 1;

-- name: CreateProduct :exec
INSERT INTO products (
    name, slug, description, origin_price, sale_price, specs
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: UpdateProduct :exec
UPDATE products
    set name = $2,
    slug = $3
WHERE id = $1;

-- name: DeleteProducts :exec
DELETE FROM products
WHERE id IN (sqlc.slice(ids));
