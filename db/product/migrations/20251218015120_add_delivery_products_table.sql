-- +goose Up
-- +goose StatementBegin
ALTER TABLE products
ADD COLUMN weight INT DEFAULT 0,
ADD COLUMN long INT DEFAULT 0,
ADD COLUMN wide INT DEFAULT 0,
ADD COLUMN high INT DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
