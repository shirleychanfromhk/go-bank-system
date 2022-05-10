package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransactionTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println("Before Transaction: ", account1.Balance, account2.Balance)

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

	existed := make(map[int]bool)

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

		// check account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		fmt.Println("Transaction: ", fromAccount.Balance, toAccount.Balance)

		// check accounts' balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println("After Transaction: ", updateAccount1.Balance, updateAccount2.Balance)
	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)
}
