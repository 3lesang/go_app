-- +goose Up
-- +goose StatementBegin
CREATE TABLE product_hotspots (
  id BIGSERIAL PRIMARY KEY,
  product_id BIGINT NOT NULL,
  hotspot_id BIGINT NOT NULL,
  x REAL NOT NULL,
  y REAL NOT NULL,
  FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
  FOREIGN KEY (hotspot_id) REFERENCES hotspots(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product_hotspots CASCADE;
-- +goose StatementEnd
