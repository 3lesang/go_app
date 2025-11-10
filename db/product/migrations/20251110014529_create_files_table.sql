-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS files (
  id BIGSERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS files CASCADE;
-- +goose StatementEnd
