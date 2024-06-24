package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T) TransferTxResult {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	amount := int64(10)
	store := NewStore(testDB)
	result, err := store.TransferTx(context.Background(), TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        amount,
	})
	// assertion testing
	require.NoError(t, err)
	require.NotEmpty(t, result)
	diff1 := account1.Balance - result.FromAccount.Balance // 874 - 864 = 10
	diff2 := result.ToAccount.Balance - account2.Balance   // 555 - 545 = 10
	require.Equal(t, diff1, diff2)
	require.True(t, diff1 > 0)

	return result
}

func TestCreateTransfer(t *testing.T) {
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	// unit tests should be independent from each other
	result1 := createRandomTransfer(t)
	result2, err := testQueries.GetTransfer(context.Background(), result1.Transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result2)

}

func TestUpdateTransfer(t *testing.T) {
	result := createRandomTransfer(t)

	arg := UpdateTransferParams{
		ID:     result.Transfer.ID,
		Amount: 5,
	}
	result1, err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	result.Transfer.Amount = arg.Amount
	require.Equal(t, result.FromAccount.ID, result1.FromAccountID)

}

func TestDeleteTransfer(t *testing.T) {
	result := createRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(), result.Transfer.ID)
	require.NoError(t, err)
	account2, err := testQueries.GetTransfer(context.Background(), result.Transfer.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListTransfers(t *testing.T) {
	for i := 1; i <= 10; i++ {
		createRandomTransfer(t)
	}

	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
