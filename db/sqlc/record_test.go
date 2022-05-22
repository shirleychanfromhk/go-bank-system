package db

import (
	"context"
	"simplebank/db/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomRecord(t *testing.T, account Account) Record {
	arg := CreateRecordParams{
		AccountID: account.ID,
		Amount:    util.RandomBalance(),
	}

	record, err := testQueries.CreateRecord(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, record)

	require.Equal(t, arg.AccountID, record.AccountID)
	require.Equal(t, arg.Amount, record.Amount)

	require.NotZero(t, record.AccountID)
	require.NotZero(t, record.CreatedAt)

	return record
}

func TestGetRecord(t *testing.T) {
	account := createRandomAccount(t)
	record1 := createRandomRecord(t, account)
	record2, err := testQueries.GetRecord(context.Background(), record1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, record2)

	require.Equal(t, record1.ID, record2.ID)
	require.Equal(t, record1.AccountID, record2.AccountID)
	require.Equal(t, record1.Amount, record2.Amount)
	require.WithinDuration(t, record1.CreatedAt, record2.CreatedAt, time.Second)
}

func TestCreateRecord(t *testing.T) {
	createRandomRecord(t, createRandomAccount(t))
}

func TestListRecords(t *testing.T) {
	account := createRandomAccount(t)
	for i := 0; i < 10; i++ {
		createRandomRecord(t, account)
	}

	arg := ListRecordsParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	records, err := testQueries.ListRecords(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, records, 5)

	for _, records := range records {
		require.NotEmpty(t, records)
	}

}
