// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package authors

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const CreateAuthor = `-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio, created_at
) VALUES (
  $1, $2, $3
)
RETURNING id, name, bio, created_at
`

type CreateAuthorParams struct {
	Name      string           `json:"name"`
	Bio       pgtype.Text      `json:"bio"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

// http: POST /authors
func (q *Queries) CreateAuthor(ctx context.Context, db DBTX, arg *CreateAuthorParams) (*Authors, error) {
	row := db.QueryRow(ctx, CreateAuthor, arg.Name, arg.Bio, arg.CreatedAt)
	var i Authors
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Bio,
		&i.CreatedAt,
	)
	return &i, err
}

const DeleteAuthor = `-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1
`

// http: DELETE /authors/{id}
func (q *Queries) DeleteAuthor(ctx context.Context, db DBTX, id int64) error {
	_, err := db.Exec(ctx, DeleteAuthor, id)
	return err
}

const GetAuthor = `-- name: GetAuthor :one
SELECT id, name, bio, created_at FROM authors
WHERE id = $1 LIMIT 1
`

// http: GET /authors/{id}
func (q *Queries) GetAuthor(ctx context.Context, db DBTX, id int64) (*Authors, error) {
	row := db.QueryRow(ctx, GetAuthor, id)
	var i Authors
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Bio,
		&i.CreatedAt,
	)
	return &i, err
}

const ListAuthors = `-- name: ListAuthors :many
SELECT id, name, bio, created_at FROM authors
ORDER BY name
`

// http: GET /authors
func (q *Queries) ListAuthors(ctx context.Context, db DBTX) ([]*Authors, error) {
	rows, err := db.Query(ctx, ListAuthors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Authors{}
	for rows.Next() {
		var i Authors
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Bio,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const UpdateAuthorBio = `-- name: UpdateAuthorBio :exec
UPDATE authors
SET bio = $1
WHERE id = $2
`

type UpdateAuthorBioParams struct {
	Bio pgtype.Text `json:"bio"`
	ID  int64       `json:"id"`
}

// http: PATCH /authors/{id}/bio
func (q *Queries) UpdateAuthorBio(ctx context.Context, db DBTX, arg *UpdateAuthorBioParams) error {
	_, err := db.Exec(ctx, UpdateAuthorBio, arg.Bio, arg.ID)
	return err
}
