-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS review_files (
  name TEXT NOT NULL,
  review_id BIGINT NOT NULL REFERENCES reviews (id) on DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS review_files CASCADE;
-- +goose StatementEnd
