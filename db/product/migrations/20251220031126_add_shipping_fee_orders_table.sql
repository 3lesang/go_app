-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
ADD COLUMN shipping_fee_amount INT DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
