-- +goose Up
-- +goose StatementBegin
CREATE TABLE discount_conditions (
  id BIGSERIAL PRIMARY KEY,
  discount_id BIGINT NOT NULL,
  condition_type TEXT NOT NULL, -- 'product' | 'collection' | 'order_amount' | 'quantity' | 'customer'
  operator TEXT NOT NULL, -- 'eq' | 'gt' | 'gte' | 'lt' | 'lte' | 'in' | 'not_in'
  value TEXT NOT NULL, -- JSON (product IDs, min_subtotal, min_qty, etc)
  FOREIGN KEY (discount_id) REFERENCES discounts (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS discount_conditions CASCADE;
-- +goose StatementEnd
