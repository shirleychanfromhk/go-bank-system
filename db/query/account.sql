-- name: CreateAccount :one
INSERT INTO accounts (
    username,
    balance,
    currency,
    location
) VALUES (
             $1, $2, $3, $4
         )
    RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts
WHERE username = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateAccount :one
UPDATE accounts
SET username = $2, balance = $3, currency = $4, location = $5
WHERE id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
    RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;