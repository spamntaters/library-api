-- name: CreateBook :one
INSERT INTO books (title, author_id, isbn, isbn13, published_date, page_count, description, cover_url, owned)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, title, author_id, isbn, isbn13, published_date, page_count, description, cover_url, owned, created_at, updated_at;

-- name: GetBook :one
SELECT id, title, author_id, isbn, isbn13, published_date, page_count, description, cover_url, owned, created_at, updated_at
FROM books
WHERE id = $1;

-- name: ListBooks :many
SELECT id, title, author_id, isbn, isbn13, published_date, page_count, description, cover_url, owned, created_at, updated_at
FROM books
WHERE (sqlc.narg(owned)::boolean IS NULL OR owned = sqlc.narg(owned))
  AND (sqlc.narg(author_id)::int IS NULL OR author_id = sqlc.narg(author_id))
ORDER BY title;

-- name: UpdateBook :one
UPDATE books
SET title = COALESCE($2, title),
    isbn = COALESCE($3, isbn),
    isbn13 = COALESCE($4, isbn13),
    published_date = COALESCE($5, published_date),
    page_count = COALESCE($6, page_count),
    description = COALESCE($7, description),
    updated_at = NOW()
WHERE id = $1
RETURNING id, title, author_id, isbn, isbn13, published_date, page_count, description, cover_url, owned, created_at, updated_at;

-- name: DeleteBook :exec
DELETE FROM books
WHERE id = $1;

-- name: ToggleOwned :one
UPDATE books
SET owned = NOT owned,
    updated_at = NOW()
WHERE id = $1
RETURNING id, title, author_id, isbn, isbn13, published_date, page_count, description, cover_url, owned, created_at, updated_at;

-- name: GetBookByISBN :one
SELECT id, title, author_id, isbn, isbn13, published_date, page_count, description, cover_url, owned, created_at, updated_at
FROM books
WHERE isbn = $1 OR isbn13 = $1;
