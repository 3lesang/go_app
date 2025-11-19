-- name: UpsertCustomerUsage :exec
INSERT INTO discount_customer_usages (discount_id, customer_id, used_count)
VALUES ($1, $2, 1)
ON CONFLICT (discount_id, customer_id) DO UPDATE
SET used_count = discount_customer_usages.used_count + 1;

-- name: GetCustomerUsage :one
SELECT used_count FROM discount_customer_usages WHERE discount_id = $1 AND customer_id = $2;
