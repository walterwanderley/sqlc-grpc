-- name: GetAuthor :one
-- http: GET /authors/{id}
SELECT * FROM authors
WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
-- http: GET /authors
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :one
-- http: POST /authors
INSERT INTO authors (
  name, bio, created_at
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteAuthor :exec
-- http: DELETE /authors/{id}
DELETE FROM authors
WHERE id = $1;

-- name: UpdateAuthorBio :exec
-- http: PATCH /authors/{id}/bio
UPDATE authors
SET bio = $1
WHERE id = $2;