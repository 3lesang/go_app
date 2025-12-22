-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
ADD COLUMN status VARCHAR(50) DEFAULT 'pending';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
