-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product_tags (
  name TEXT NOT NULL,
  product_id BIGINT NOT NULL REFERENCES products (id) ON DELETE CASCADE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_tags;

-- +goose StatementEnd
