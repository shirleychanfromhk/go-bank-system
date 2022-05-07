package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

// Create a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// Executes a function within a database transaction
func (store *Store) execTX(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransactionTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransactionTxResult struct {
	Transaction Transaction `json:"transaction"`
	FromAccount Account     `json:"from_account"`
	ToAccount   Account     `json:"to_account"`
	FromRecord  Record      `json:"from_record"`
	ToRecord    Record      `json:"to_record"`
}

func (store *Store) TransactionTx(ctx context.Context, arg TransactionTxParams) (TransactionTxResult, error) {
	var result TransactionTxResult

	err := store.execTX(ctx, func(q *Queries) error {
		var err error

		result.Transaction, err = q.CreateTransaction(ctx, CreateTransactionParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromRecord, err = q.CreateRecord(ctx, CreateRecordParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToRecord, err = q.CreateRecord(ctx, CreateRecordParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}
