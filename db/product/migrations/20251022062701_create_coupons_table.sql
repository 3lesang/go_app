-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS coupons (
  id BIGSERIAL PRIMARY KEY,
  code TEXT UNIQUE NOT NULL,
  description TEXT,
  discount_percent DECIMAL(5, 2) CHECK (discount_percent BETWEEN 0 AND 100),
  valid_from TIMESTAMP,
  valid_until TIMESTAMP,
  is_active BOOLEAN DEFAULT TRUE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS coupons CASCADE;

-- +goose StatementEnd
