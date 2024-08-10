-- name: ListFiles :many
SELECT * FROM files;

-- name: CreateFile :one
INSERT INTO files (key, name, size, last_modified)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteProjectByKey :exec
DELETE FROM files
WHERE key = $1;