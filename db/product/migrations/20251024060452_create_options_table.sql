-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS options (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL DEFAULT '',
  no INT NOT NULL,
  product_id BIGINT NOT NULL REFERENCES products (id) ON DELETE CASCADE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS options CASCADE;

-- +goose StatementEnd
