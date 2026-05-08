-- name: CreateAuthor :one
INSERT INTO authors (name, bio)
VALUES ($1, $2)
RETURNING id, name, bio, created_at, updated_at;

-- name: GetAuthor :one
SELECT id, name, bio, created_at, updated_at
FROM authors
WHERE id = $1;

-- name: ListAuthors :many
SELECT id, name, bio, created_at, updated_at
FROM authors
ORDER BY name;

-- name: UpdateAuthor :one
UPDATE authors
SET name = COALESCE($2, name),
    bio = COALESCE($3, bio),
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, bio, created_at, updated_at;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1;
