-- name: NewEntry :one
INSERT INTO entries (
    acc_id, amount
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = $1 LIMIT 1;

-- name: GetAllEntries :many
SELECT * FROM entries
WHERE acc_id = $1;

