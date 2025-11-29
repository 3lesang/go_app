-- +goose Up
-- +goose StatementBegin
CREATE TABLE hotspots (
  id BIGSERIAL PRIMARY KEY,
  file TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS hotspots CASCADE;
-- +goose StatementEnd
