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
ORDER BY name
LIMIT sqlc.narg('limit')
OFFSET sqlc.narg('offset');

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

-- name: GetBooksByAuthor :many
SELECT id, title, author_id, isbn, isbn13, published_date, page_count, description, cover_url, owned, created_at, updated_at
FROM books
WHERE author_id = $1
ORDER BY title;

-- name: GetSeriesByAuthor :many
SELECT DISTINCT s.id, s.name, s.description, s.created_at, s.updated_at
FROM series s
JOIN series_books sb ON sb.series_id = s.id
JOIN books b ON b.id = sb.book_id
WHERE b.author_id = $1
ORDER BY s.name;
