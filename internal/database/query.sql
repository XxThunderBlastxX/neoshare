-- name: ListFiles :many
SELECT * FROM file;

-- name: CreateFile :one
INSERT INTO file (id, name, key, size, last_modified)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteProjectByKey :exec
DELETE FROM file
WHERE key = $1;