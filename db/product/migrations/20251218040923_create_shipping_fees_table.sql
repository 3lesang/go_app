-- +goose Up
-- +goose StatementBegin
CREATE TABLE shipping_fees (
  id BIGSERIAL PRIMARY KEY,
  min_weight INTEGER NOT NULL DEFAULT 0,
  max_weight INTEGER NOT NULL DEFAULT 0,
  fee_amount INTEGER NOT NULL DEFAULT 0,
  min_order_value INTEGER,
  free_shipping BOOLEAN DEFAULT FALSE,
  shipping_method TEXT,
  effective_from DATE,
  effective_to DATE,
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shipping_fees;
-- +goose StatementEnd
