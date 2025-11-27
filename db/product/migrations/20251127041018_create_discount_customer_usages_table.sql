-- +goose Up
-- +goose StatementBegin
CREATE TABLE discount_customer_usages (
  id BIGSERIAL PRIMARY KEY,
  discount_id BIGINT NOT NULL REFERENCES discounts(id) ON DELETE CASCADE,
  customer_id BIGINT NOT NULL,
  used_count INTEGER DEFAULT 0,
  UNIQUE (discount_id, customer_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS discount_customer_usages CASCADE;
-- +goose StatementEnd
