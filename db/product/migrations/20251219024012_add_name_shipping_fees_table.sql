-- +goose Up
-- +goose StatementBegin
ALTER TABLE shipping_fees
ADD COLUMN name TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
