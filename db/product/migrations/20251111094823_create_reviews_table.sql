-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS reviews (
  id BIGSERIAL PRIMARY KEY,
  rating INT CHECK (rating BETWEEN 1 AND 5),
  comment TEXT,
  has_file BOOLEAN NOT NULL DEFAULT FALSE,
  product_id BIGINT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
  customer_id BIGINT NOT NULL REFERENCES customers (id) ON DELETE CASCADE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reviews CASCADE;
-- +goose StatementEnd
