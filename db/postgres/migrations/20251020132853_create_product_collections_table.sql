-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product_collections (
    product_id INT REFERENCES products(id) ON DELETE CASCADE,
    collection_id INT REFERENCES collections(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, collection_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE product_collections;
-- +goose StatementEnd
