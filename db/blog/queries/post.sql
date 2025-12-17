-- name: CountPosts :one
SELECT
  COUNT(*)
FROM
  posts;

-- name: CountPublicPosts :one
SELECT
  COUNT(*)
FROM
  posts
WHERE is_active = true;

-- name: GetPosts :many
SELECT id, title, slug, file
FROM posts
LIMIT
  $1
OFFSET
  $2;

-- name: GetPublicPosts :many
SELECT id, title, slug, file, meta_title, meta_description, created_at
FROM posts
WHERE is_active = true
LIMIT
  $1
OFFSET
  $2;

-- name: GetPost :one
SELECT id, title, slug, file
FROM posts
WHERE id = $1;

-- name: GetPostBySlug :one
SELECT id, title, slug, file, meta_title, meta_description, created_at
FROM posts
WHERE slug = $1;

-- name: CreatePost :one
INSERT INTO
  posts (title, slug, file)
VALUES
  ($1, $2, $3)
RETURNING id;

-- name: UpdatePost :exec
UPDATE posts
SET
  title = $2,
  slug = $3,
  file = $4
WHERE
  id = $1;

-- name: BulkDeletePosts :exec
DELETE FROM posts
WHERE
  id = ANY ($1::bigint[]);
