-- name: GetUrl :one
SELECT * FROM shurls WHERE id = $1 LIMIT 1;

-- name: GetUrlFromHash :one
SELECT * FROM shurls WHERE hash = $1 LIMIT 1;

-- name: CreateUrl :one
INSERT INTO shurls (hash, url, owner) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateUrl :exec
UPDATE shurls SET hits = $2 WHERE id = $1;