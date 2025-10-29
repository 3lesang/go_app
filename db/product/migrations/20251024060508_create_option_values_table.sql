-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS option_values (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL DEFAULT '',
  no INT,
  option_id BIGINT NOT NULL REFERENCES options (id) ON DELETE CASCADE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS option_values CASCADE;

-- +goose StatementEnd
