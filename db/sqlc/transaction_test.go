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

func TestCreateTransaction(t *testing.T) {
	mockAccount1 := createRandomAccount(t)
	mockAccount2 := createRandomAccount(t)
	createRandomTransfer(t, mockAccount1, mockAccount2)
}

func TestGetTransaction(t *testing.T) {
	mockAccount1 := createRandomAccount(t)
	mockAccount2 := createRandomAccount(t)
	mockTransaction := createRandomTransfer(t, mockAccount1, mockAccount2)

	transaction, err := testQueries.GetTransaction(context.Background(), mockTransaction.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transaction)

	require.Equal(t, mockAccount1.ID, transaction.FromAccountID)
	require.Equal(t, mockAccount2.ID, transaction.ToAccountID)
	require.Equal(t, transaction.Amount, transaction.Amount)

	require.NotZero(t, transaction.ID)
	require.NotZero(t, transaction.CreatedAt)
}

func TestListTransactions(t *testing.T) {
	// TODO - BK2734 Write the unit test method for [get] bank transaction's' operation
	mockAccount1 := createRandomAccount(t)
	mockAccount2 := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, mockAccount1, mockAccount2)
		createRandomTransfer(t, mockAccount2, mockAccount1)
	}

	arg := ListTransactionParams{
		FromAccountID: mockAccount1.ID,
		ToAccountID:   mockAccount1.ID,
		Limit:         5,
		Offset:        5,
	}

	transactions, err := testQueries.ListTransaction(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transactions, 5)

	for _, transaction := range transactions {
		require.NotEmpty(t, transaction)
		require.True(t, transaction.FromAccountID == mockAccount1.ID || transaction.ToAccountID == mockAccount1.ID)
	}
}
