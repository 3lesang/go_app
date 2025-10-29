-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product_files (
  name TEXT,
  no INT NOT NULL,
  is_primary BOOLEAN NOT NULL DEFAULT FALSE,
  product_id BIGINT NOT NULL REFERENCES products (id) ON DELETE CASCADE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_files CASCADE;

-- +goose StatementEnd
