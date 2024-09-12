-- name: NewAcc :one
INSERT INTO accounts (
    owner, balance, currency
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetAcc :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: GetAllAccs :many
SELECT * FROM accounts
ORDER BY id;

-- name: UpdateAcc :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: UpdateAccBalance :one
UPDATE accounts
SET balance =  balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAcc :exec
DELETE FROM accounts
WHERE id = $1;