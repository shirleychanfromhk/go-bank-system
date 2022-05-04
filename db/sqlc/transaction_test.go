package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"simplebank/db/util"
	"testing"
)

func createRandomTransfer(t *testing.T, account1, account2 Account) Transaction {
	arg := CreateTransactionParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomBalance(),
	}

	transaction, err := testQueries.CreateTransaction(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transaction)

	require.Equal(t, arg.FromAccountID, transaction.FromAccountID)
	require.Equal(t, arg.ToAccountID, transaction.ToAccountID)
	require.Equal(t, arg.Amount, transaction.Amount)

	require.NotZero(t, transaction.ID)
	require.NotZero(t, transaction.CreatedAt)

	return transaction
}

func TestGetTransaction(t *testing.T) {
	// TODO - BK2734 Write the unit test method for [get] bank transaction operation
}

func TestCreateTransaction(t *testing.T) {
	// TODO - BK2734 Write the unit test method for [create] bank transaction operation
}

func TestListTransactions(t *testing.T) {
	// TODO - BK2734 Write the unit test method for [get] bank transaction's' operation
}

func TestUpdateTransaction(t *testing.T) {
	// TODO - BK2734 Write the unit test method for [update] bank transaction operation
}

func TestDeleteTransaction(t *testing.T) {
	// TODO - BK2734 Write the unit test method for [delete] bank transaction operation
}
