-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_items (
  id BIGSERIAL PRIMARY KEY,
  quantity INT NOT NULL DEFAULT 0,
  sale_price INT NOT NULL DEFAULT 0,
  order_id BIGINT NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
  product_id BIGINT REFERENCES products (id) ON DELETE SET NULL,
  variant_id BIGINT REFERENCES variants (id) ON DELETE SET NULL
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_items CASCADE;

-- +goose StatementEnd
