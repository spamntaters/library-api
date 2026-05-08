-- name: CreateSeries :one
INSERT INTO series (name, description)
VALUES ($1, $2)
RETURNING id, name, description, created_at, updated_at;

-- name: GetSeries :one
SELECT id, name, description, created_at, updated_at
FROM series
WHERE id = $1;

-- name: ListSeries :many
SELECT id, name, description, created_at, updated_at
FROM series
ORDER BY name
LIMIT sqlc.narg('limit')
OFFSET sqlc.narg('offset');

-- name: UpdateSeries :one
UPDATE series
SET name = COALESCE($2, name),
    description = COALESCE($3, description),
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, description, created_at, updated_at;

-- name: DeleteSeries :exec
DELETE FROM series
WHERE id = $1;

-- name: AddBookToSeries :one
INSERT INTO series_books (series_id, book_id, position)
VALUES ($1, $2, $3)
RETURNING series_id, book_id, position;

-- name: RemoveBookFromSeries :exec
DELETE FROM series_books
WHERE series_id = $1 AND book_id = $2;

-- name: GetSeriesBooks :many
SELECT sb.series_id, sb.book_id, sb.position,
       b.id as book_id, b.title, b.author_id, b.isbn, b.isbn13,
       b.published_date, b.page_count, b.description, b.cover_url, b.owned,
       b.created_at as book_created_at, b.updated_at as book_updated_at
FROM series_books sb
JOIN books b ON b.id = sb.book_id
WHERE sb.series_id = $1
ORDER BY sb.position;

-- name: GetMissingBooks :many
SELECT b.id, b.title, b.author_id, b.isbn, b.isbn13, b.published_date,
       b.page_count, b.description, b.cover_url, b.owned, b.created_at, b.updated_at
FROM series_books sb
JOIN books b ON b.id = sb.book_id
WHERE sb.series_id = $1 AND b.owned = false
ORDER BY sb.position;

-- name: GetSeriesBooksByBookID :many
SELECT series_id, book_id, position
FROM series_books
WHERE book_id = $1
ORDER BY position;
