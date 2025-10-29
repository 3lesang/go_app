-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product_collections (
  product_id BIGINT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
  collection_id BIGINT NOT NULL REFERENCES collections (id) ON DELETE CASCADE,
  PRIMARY KEY (product_id, collection_id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_collections CASCADE;

-- +goose StatementEnd
