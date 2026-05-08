-- name: CreateTag :one
INSERT INTO tags (name)
VALUES ($1)
RETURNING id, name;

-- name: GetTag :one
SELECT id, name
FROM tags
WHERE id = $1;

-- name: ListTags :many
SELECT id, name
FROM tags
ORDER BY name
LIMIT sqlc.narg('limit')
OFFSET sqlc.narg('offset');

-- name: AddTagToBook :exec
INSERT INTO book_tags (book_id, tag_id)
VALUES ($1, $2);

-- name: RemoveTagFromBook :exec
DELETE FROM book_tags
WHERE book_id = $1 AND tag_id = $2;

-- name: GetBookTags :many
SELECT t.id, t.name
FROM book_tags bt
JOIN tags t ON t.id = bt.tag_id
WHERE bt.book_id = $1
ORDER BY t.name;
