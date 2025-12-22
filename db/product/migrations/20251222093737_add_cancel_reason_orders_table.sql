-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
ADD COLUMN cancel_reason TEXT DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
