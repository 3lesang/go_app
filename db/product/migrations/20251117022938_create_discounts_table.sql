-- +goose Up
-- +goose StatementBegin
CREATE TABLE discounts (
  id BIGSERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  code TEXT UNIQUE, -- nullable for automatic discounts
  discount_type TEXT NOT NULL, -- 'code' | 'automatic'
  status TEXT NOT NULL, -- 'draft' | 'active' | 'scheduled' | 'expired'
  usage_limit INTEGER, -- global usage limit
  usage_count INTEGER DEFAULT 0,
  per_customer_limit INTEGER, -- optional per-customer limit
  starts_at TIMESTAMPTZ NOT NULL,
  ends_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS discounts CASCADE;
-- +goose StatementEnd
