-- name: CreateAccount :exec
INSERT INTO accounts(
	account_id,
	parent_account_id,
	account_status,
	currency_id,
	created_at
) VALUES($1,$2,$3,$4,$5);

-- name: CreateAccountBalance :exec
INSERT INTO accounts_balance(
	account_id,
	allow_negative,
	balance,
	last_ledger_id,
	last_movement_id,
	currency_id,
	created_at
) VALUES($1,$2,$3,$4,$5,$6,$7);

-- name: GetAccounts :many
SELECT *
FROM accounts
WHERE account_id = ANY($1::varchar[])
ORDER BY created_at;

-- name: GetAccountsBalance :many
SELECT ab.account_id,
	ab.allow_negative,
	ab.balance,
	ab.currency_id,
	ab.last_ledger_id,
	ab.last_movement_id,
	ab.created_at,
	ab.updated_at,
	ac.account_status
FROM accounts_balance ab,
	accounts ac
WHERE ab.account_id = ANY($1::varchar[])
	AND ab.account_id = ac.account_id;

-- name: GetMovementByIdempotencyKey :one
SELECT movement_id,
    idempotency_key,
    created_at,
    updated_at
FROM movements
WHERE idempotency_key = $1;

-- name: CreateMovement :exec
INSERT INTO movements(
	movement_id,
	idempotency_key,
	movement_status,
	created_at,
	updated_at
) VALUES($1,$2,$3,$4,$5);

-- name: GetMovement :one
SELECT * FROM movements
WHERE movement_id = $1;

-- name: GetAccountsLedgerByMovementID :many
SELECT ledger_id,
	movement_id,
	movement_sequence,
	account_id,
	amount,
	previous_ledger_id,
	client_id,
	created_at,
	client_id
FROM accounts_ledger
WHERE movement_id = $1
ORDER BY created_at;
