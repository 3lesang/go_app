-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS customers (
  id BIGSERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL,
  phone TEXT UNIQUE NOT NULL,
  avatar TEXT,
  email TEXT,
  password TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS customers CASCADE;
-- +goose StatementEnd
