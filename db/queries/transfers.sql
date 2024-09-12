-- name: NewTransfer :one
INSERT INTO transfers (
    from_acc_id, to_acc_id, amount
) VALUES (
     $1, $2, $3
) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: GetAllTransfers :many
SELECT * FROM transfers
WHERE
    from_acc_id=$1 OR
    to_acc_id=$2;
