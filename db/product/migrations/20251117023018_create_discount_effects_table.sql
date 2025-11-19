-- +goose Up
-- +goose StatementBegin
CREATE TABLE discount_effects (
  id BIGSERIAL PRIMARY KEY,
  discount_id BIGINT NOT NULL,
  effect_type TEXT NOT NULL, -- 'percent' | 'fixed' | 'free_shipping' | 'bogo'
  value TEXT, -- percent: "20", fixed: "5", bogo: JSON {"buy":1, "get":1}
  applies_to TEXT NOT NULL, -- 'entire_order' | 'specific_products' | 'specific_collections'
  FOREIGN KEY (discount_id) REFERENCES discounts (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS discount_effects CASCADE;
-- +goose StatementEnd
