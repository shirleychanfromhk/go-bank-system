package db

import (
	"context"
	"database/sql"
	"simplebank/db/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)

	checkAccount := CreateAccountParams{
		Username: user.Username,
		Balance:  util.RandomBalance(),
		Currency: util.RandomCurrency(),
		Location: util.RandomLocation(),
	}

	testAccount, err := testQueries.CreateAccount(context.Background(), checkAccount)

	require.NoError(t, err)
	require.NotEmpty(t, testAccount)

	require.Equal(t, checkAccount.Username, testAccount.Username)
	require.Equal(t, checkAccount.Balance, testAccount.Balance)
	require.Equal(t, checkAccount.Currency, testAccount.Currency)

	require.NotZero(t, testAccount.ID)
	require.NotZero(t, testAccount.CreatedAt)

	return testAccount
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Username, account2.Username)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Location, account2.Location)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomBalance(),
	}

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Username, account2.Username)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Location, account2.Location)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	// TODO - EN2360: Make the current createRandomAccount to Run in Parallel
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}

func TestGetAccountForUpdate(t *testing.T) {
	mockAccount := createRandomAccount(t)
	account, err := testQueries.GetAccountForUpdate(context.Background(), mockAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.ID, mockAccount.ID)
	require.Equal(t, account.Username, mockAccount.Username)
	require.Equal(t, account.Balance, mockAccount.Balance)
	require.Equal(t, account.Currency, mockAccount.Currency)
	require.Equal(t, account.Location, mockAccount.Location)
	require.WithinDuration(t, account.CreatedAt, mockAccount.CreatedAt, time.Second)
}
