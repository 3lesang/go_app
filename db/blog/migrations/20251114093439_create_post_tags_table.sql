-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS post_tags (
  name TEXT NOT NULL,
  post_id BIGINT NOT NULL REFERENCES posts (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS post_tags CASCADE;
-- +goose StatementEnd
