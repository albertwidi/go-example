// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

const createAccount = `-- name: CreateAccount :exec
INSERT INTO accounts(
	account_id,
	parent_account_id,
	account_status,
	currency_id,
	created_at
) VALUES($1,$2,$3,$4,$5)
`

type CreateAccountParams struct {
	AccountID       string
	ParentAccountID string
	AccountStatus   AccountStatus
	CurrencyID      int32
	CreatedAt       time.Time
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) error {
	_, err := q.db.Exec(ctx, createAccount,
		arg.AccountID,
		arg.ParentAccountID,
		arg.AccountStatus,
		arg.CurrencyID,
		arg.CreatedAt,
	)
	return err
}

const createAccountBalance = `-- name: CreateAccountBalance :exec
INSERT INTO accounts_balance(
	account_id,
	allow_negative,
	balance,
	last_ledger_id,
	currency_id,
	created_at
) VALUES($1,$2,$3,$4,$5,$6)
`

type CreateAccountBalanceParams struct {
	AccountID     string
	AllowNegative bool
	Balance       decimal.Decimal
	LastLedgerID  string
	CurrencyID    int32
	CreatedAt     time.Time
}

func (q *Queries) CreateAccountBalance(ctx context.Context, arg CreateAccountBalanceParams) error {
	_, err := q.db.Exec(ctx, createAccountBalance,
		arg.AccountID,
		arg.AllowNegative,
		arg.Balance,
		arg.LastLedgerID,
		arg.CurrencyID,
		arg.CreatedAt,
	)
	return err
}

const createMovements = `-- name: CreateMovements :exec
INSERT INTO movements(
	movement_id,
	idempotency_key,
	created_at,
	updated_at
) VALUES($1,$2,$3,$4)
`

type CreateMovementsParams struct {
	MovementID     string
	IdempotencyKey string
	CreatedAt      time.Time
	UpdatedAt      sql.NullTime
}

func (q *Queries) CreateMovements(ctx context.Context, arg CreateMovementsParams) error {
	_, err := q.db.Exec(ctx, createMovements,
		arg.MovementID,
		arg.IdempotencyKey,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const getAccounts = `-- name: GetAccounts :many
SELECT account_id, parent_account_id, account_status, currency_id, created_at, updated_at
FROM accounts
WHERE account_id = ANY($1::varchar[])
ORDER BY created_at
`

func (q *Queries) GetAccounts(ctx context.Context, dollar_1 []string) ([]Account, error) {
	rows, err := q.db.Query(ctx, getAccounts, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.AccountID,
			&i.ParentAccountID,
			&i.AccountStatus,
			&i.CurrencyID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAccountsBalance = `-- name: GetAccountsBalance :many
SELECT ab.account_id,
	ab.allow_negative,
	ab.balance,
	ab.currency_id,
	ab.last_ledger_id,
	ab.created_at,
	ab.updated_at,
	ac.account_status
FROM accounts_balance ab,
	accounts ac
WHERE ab.account_id = ANY($1::varchar[])
	AND ab.account_id = ac.account_id
`

type GetAccountsBalanceRow struct {
	AccountID     string
	AllowNegative bool
	Balance       decimal.Decimal
	CurrencyID    int32
	LastLedgerID  string
	CreatedAt     time.Time
	UpdatedAt     sql.NullTime
	AccountStatus AccountStatus
}

func (q *Queries) GetAccountsBalance(ctx context.Context, dollar_1 []string) ([]GetAccountsBalanceRow, error) {
	rows, err := q.db.Query(ctx, getAccountsBalance, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAccountsBalanceRow
	for rows.Next() {
		var i GetAccountsBalanceRow
		if err := rows.Scan(
			&i.AccountID,
			&i.AllowNegative,
			&i.Balance,
			&i.CurrencyID,
			&i.LastLedgerID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.AccountStatus,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAccountsLedgerByMovementID = `-- name: GetAccountsLedgerByMovementID :many
SELECT ledger_id,
	movement_id,
	movement_sequence,
	amount,
	previous_ledger_id,
	created_at,
	timestamp,
	client_id
FROM accounts_ledger
WHERE movement_id = $1
ORDER BY timestamp
`

type GetAccountsLedgerByMovementIDRow struct {
	LedgerID         string
	MovementID       string
	MovementSequence int32
	Amount           decimal.Decimal
	PreviousLedgerID string
	CreatedAt        time.Time
	Timestamp        int64
	ClientID         sql.NullString
}

func (q *Queries) GetAccountsLedgerByMovementID(ctx context.Context, movementID string) ([]GetAccountsLedgerByMovementIDRow, error) {
	rows, err := q.db.Query(ctx, getAccountsLedgerByMovementID, movementID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAccountsLedgerByMovementIDRow
	for rows.Next() {
		var i GetAccountsLedgerByMovementIDRow
		if err := rows.Scan(
			&i.LedgerID,
			&i.MovementID,
			&i.MovementSequence,
			&i.Amount,
			&i.PreviousLedgerID,
			&i.CreatedAt,
			&i.Timestamp,
			&i.ClientID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMovementByIdempotencyKey = `-- name: GetMovementByIdempotencyKey :one
SELECT movement_id,
    idempotency_key,
    created_at,
    updated_at
FROM movements
WHERE idempotency_key = $1
`

func (q *Queries) GetMovementByIdempotencyKey(ctx context.Context, idempotencyKey string) (Movement, error) {
	row := q.db.QueryRow(ctx, getMovementByIdempotencyKey, idempotencyKey)
	var i Movement
	err := row.Scan(
		&i.MovementID,
		&i.IdempotencyKey,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
