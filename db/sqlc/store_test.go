package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransactionTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// Run in goroutines
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransactionTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransactionTx(context.Background(), TransactionTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check Transaction
		transaction := result.Transaction
		require.NotEmpty(t, transaction)
		require.Equal(t, account1.ID, transaction.FromAccountID)
		require.Equal(t, account2.ID, transaction.ToAccountID)
		require.Equal(t, amount, transaction.Amount)
		require.NotZero(t, transaction.ID)
		require.NotZero(t, transaction.CreatedAt)
		_, err = store.GetTransaction(context.Background(), transaction.ID)
		require.NoError(t, err)

		// check From Record
		fromRecord := result.FromRecord
		require.NotEmpty(t, fromRecord)
		require.Equal(t, account1.ID, fromRecord.AccountID)
		require.Equal(t, -amount, fromRecord.Amount)
		require.NotZero(t, fromRecord.ID)
		require.NotZero(t, fromRecord.CreatedAt)
		_, err = store.GetRecord(context.Background(), fromRecord.ID)
		require.NoError(t, err)

		// check To Record
		toRecord := result.ToRecord
		require.NotEmpty(t, toRecord)
		require.Equal(t, account2.ID, toRecord.AccountID)
		require.Equal(t, amount, toRecord.Amount)
		require.NotZero(t, toRecord.ID)
		require.NotZero(t, toRecord.CreatedAt)
		_, err = store.GetRecord(context.Background(), toRecord.ID)
		require.NoError(t, err)
	}
}
