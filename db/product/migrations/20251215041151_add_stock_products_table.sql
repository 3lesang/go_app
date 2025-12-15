-- +goose Up
-- +goose StatementBegin
ALTER TABLE products
ADD COLUMN stock INT DEFAULT 0;

ALTER TABLE products
ADD COLUMN sku TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
