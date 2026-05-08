-- name: CreateBook :one
INSERT INTO books (title, author_id, isbn, isbn13, published_date, page_count, description, cover_url, owned)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, title, author_id, isbn, isbn13, published_date, page_count, description, cover_url, owned, created_at, updated_at;

-- name: GetBook :one
SELECT id, title, author_id, isbn, isbn13, published_date, page_count, description, cover_url, owned, created_at, updated_at
FROM books
WHERE id = $1;

-- name: ListBooks :many
SELECT b.id, b.title, b.author_id, b.isbn, b.isbn13, b.published_date, b.page_count, b.description, b.cover_url, b.owned, b.created_at, b.updated_at
FROM books b
WHERE (sqlc.narg(owned)::boolean IS NULL OR b.owned = sqlc.narg(owned))
  AND (sqlc.narg(author_id)::int IS NULL OR b.author_id = sqlc.narg(author_id))
  AND (sqlc.narg(tag_id)::int IS NULL OR EXISTS (
    SELECT 1 FROM book_tags bt WHERE bt.book_id = b.id AND bt.tag_id = sqlc.narg(tag_id)
  ))
ORDER BY b.title
LIMIT sqlc.narg('limit')
OFFSET sqlc.narg('offset');

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
