-- +goose Up
-- +goose StatementBegin
ALTER TABLE customers
ADD COLUMN phone_verified  BOOLEAN DEFAULT FALSE,
ADD COLUMN zns_otp TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
