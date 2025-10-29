-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS variant_options (
  variant_id BIGINT NOT NULL REFERENCES variants (id) ON DELETE CASCADE,
  option_value_id BIGINT NOT NULL REFERENCES option_values (id) ON DELETE CASCADE,
  option_id BIGINT NOT NULL REFERENCES options (id) ON DELETE CASCADE,
  PRIMARY KEY (variant_id, option_id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS variant_options CASCADE;

-- +goose StatementEnd
