-- name: GetAuthor :one
-- http: GET /authors/{author_id}
SELECT * FROM authors
WHERE author_id = $1;

-- name: GetBook :one
-- http: GET /books/{book_id}
SELECT * FROM books
WHERE book_id = $1;

-- name: DeleteBook :exec
-- http: DELETE /books/{book_id}
DELETE FROM books
WHERE book_id = $1;

-- name: BooksByTitleYear :many
-- http: GET /books
SELECT * FROM books
WHERE title = $1 AND year = $2;

-- name: BooksByTags :many
-- skip: true
SELECT 
  book_id,
  title,
  name,
  isbn,
  tags
FROM books
LEFT JOIN authors ON books.author_id = authors.author_id
WHERE tags && $1::varchar[];

-- name: CreateAuthor :one
-- http: POST /authors
INSERT INTO authors (name) VALUES ($1)
RETURNING *;

-- name: CreateBook :one
-- http: POST /books
INSERT INTO books (
    author_id,
    isbn,
    book_type,
    title,
    year,
    available,
    tags
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: UpdateBook :exec
-- http: PUT /books/{book_id}
UPDATE books
SET title = $1, tags = $2, book_type = $3
WHERE book_id = $4;

-- name: UpdateBookISBN :exec
-- http: PATCH /books/{book_id}/isbn
UPDATE books
SET isbn = $1
WHERE book_id = $2;