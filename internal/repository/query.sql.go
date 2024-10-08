// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package repository

import (
	"context"
	"time"
)

const createFile = `-- name: CreateFile :one
INSERT INTO files (key, name, size, last_modified)
VALUES ($1, $2, $3, $4)
RETURNING key, name, size, last_modified
`

type CreateFileParams struct {
	Key          string
	Name         string
	Size         int32
	LastModified time.Time
}

func (q *Queries) CreateFile(ctx context.Context, arg CreateFileParams) (File, error) {
	row := q.db.QueryRowContext(ctx, createFile,
		arg.Key,
		arg.Name,
		arg.Size,
		arg.LastModified,
	)
	var i File
	err := row.Scan(
		&i.Key,
		&i.Name,
		&i.Size,
		&i.LastModified,
	)
	return i, err
}

const deleteProjectByKey = `-- name: DeleteProjectByKey :exec
DELETE FROM files
WHERE key = $1
`

func (q *Queries) DeleteProjectByKey(ctx context.Context, key string) error {
	_, err := q.db.ExecContext(ctx, deleteProjectByKey, key)
	return err
}

const listFiles = `-- name: ListFiles :many
SELECT key, name, size, last_modified FROM files
`

func (q *Queries) ListFiles(ctx context.Context) ([]File, error) {
	rows, err := q.db.QueryContext(ctx, listFiles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []File
	for rows.Next() {
		var i File
		if err := rows.Scan(
			&i.Key,
			&i.Name,
			&i.Size,
			&i.LastModified,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
