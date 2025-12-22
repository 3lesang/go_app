-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
ADD COLUMN code TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
