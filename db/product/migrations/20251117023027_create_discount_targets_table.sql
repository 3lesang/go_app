-- +goose Up
-- +goose StatementBegin
CREATE TABLE discount_targets (
  id BIGSERIAL PRIMARY KEY,
  discount_id BIGINT NOT NULL,
  target_type TEXT NOT NULL, -- 'product' | 'collection'
  target_id INTEGER NOT NULL,
  FOREIGN KEY (discount_id) REFERENCES discounts (id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS discount_targets CASCADE;
-- +goose StatementEnd
