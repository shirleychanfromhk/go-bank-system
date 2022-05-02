package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	checkAccount := CreateAccountParams{
		Username: "testUsername",
		Balance:  100,
		Currency: "USD",
		Location: "testLocation",
	}

	testAccount, err := testQueries.CreateAccount(context.Background(), checkAccount)

	require.NoError(t, err)
	require.NotEmpty(t, testAccount)

	require.Equal(t, checkAccount.Username, testAccount.Username)
	require.Equal(t, checkAccount.Balance, testAccount.Balance)
	require.Equal(t, checkAccount.Currency, testAccount.Currency)

	require.NotZero(t, testAccount.ID)
	require.NotZero(t, testAccount.CreatedAt)
}
