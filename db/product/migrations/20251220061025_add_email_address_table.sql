-- +goose Up
-- +goose StatementBegin
ALTER TABLE addresses
ADD COLUMN email TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
