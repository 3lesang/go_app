-- name: CountPosts :one
SELECT
  COUNT(*)
FROM
  posts;

-- name: GetPosts :many
SELECT id, title, slug
FROM posts
LIMIT
  $1
OFFSET
  $2;

-- name: GetPost :one
SELECT id, title, slug
FROM posts
WHERE id = $1;

-- name: GetPostBySlug :one
SELECT id, title, slug
FROM posts
WHERE slug = $1;

-- name: CreatePost :one
INSERT INTO
  posts (title, slug)
VALUES
  ($1, $2)
RETURNING id;

-- name: UpdatePost :exec
UPDATE posts
SET
  title = $2,
  slug = $3
WHERE
  id = $1;

-- name: BulkDeletePosts :exec
DELETE FROM posts
WHERE
  id = ANY ($1::bigint[]);
