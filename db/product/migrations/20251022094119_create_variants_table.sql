-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS variants (
  id BIGSERIAL PRIMARY KEY,
  origin_price INT NOT NULL DEFAULT 0,
  sale_price INT NOT NULL DEFAULT 0,
  file TEXT,
  stock INT NOT NULL DEFAULT 0,
  sku TEXT NOT NULL DEFAULT '',
  product_id BIGINT NOT NULL REFERENCES products (id) ON DELETE CASCADE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS variants CASCADE;

-- +goose StatementEnd
