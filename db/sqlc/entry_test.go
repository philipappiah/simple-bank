package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entries {
	account := createRandomAccount(t)
	amount := int64(10)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    amount,
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	// assertion testing
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	// unit tests should be independent from each other
	account1 := createRandomEntry(t)
	account2, err := testQueries.GetEntry(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account1.AccountID, account2.AccountID)
	require.Equal(t, account1.Amount, account2.Amount)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestUpdateEntry(t *testing.T) {
	account1 := createRandomEntry(t)

	arg := UpdateEntryParams{
		ID:     account1.ID,
		Amount: 5,
	}
	account2, err := testQueries.UpdateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	account1.Amount = arg.Amount
	require.Equal(t, account1.AccountID, account2.AccountID)
	require.Equal(t, account1.Amount, account2.Amount)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestDeleteEntry(t *testing.T) {
	account1 := createRandomEntry(t)
	err := testQueries.DeleteEntry(context.Background(), account1.ID)
	require.NoError(t, err)
	account2, err := testQueries.GetEntry(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListEntries(t *testing.T) {
	for i := 1; i <= 10; i++ {
		createRandomEntry(t)
	}

	arg := ListEntriesParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
