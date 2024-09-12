package db

import (
	"context"
	"fmt"
	"github.com/Grishun/simplebank/db/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandAcc(t *testing.T) Account {
	arg := NewAccParams{
		Owner:    util.RandString(6),
		Balance:  util.RandInt(0, 1000),
		Currency: util.RandCurrency(),
	}

	acc, err := testQueries.NewAcc(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, acc)

	require.Equal(t, arg.Owner, acc.Owner)
	require.Equal(t, arg.Balance, acc.Balance)
	require.Equal(t, arg.Currency, acc.Currency)

	require.NotZero(t, acc.ID)
	require.NotZero(t, acc.CreatedAt)

	return acc
}

func TestQueries_NewAcc(t *testing.T) {
	createRandAcc(t)
}

func TestQueries_GetAcc(t *testing.T) {
	testAcc := createRandAcc(t)

	received, err := testQueries.GetAcc(context.Background(), testAcc.ID)

	require.NoError(t, err)
	require.NotEmpty(t, received)

	require.Equal(t, testAcc, received)
}

func TestQueries_UpdateAcc(t *testing.T) {
	testAcc := createRandAcc(t)

	arg := UpdateAccParams{
		Balance: util.RandInt(0, 100),
		ID:      testAcc.ID,
	}

	acc, err := testQueries.UpdateAcc(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, acc)

	require.Equal(t, testAcc.Owner, acc.Owner)
	require.Equal(t, arg.Balance, acc.Balance)
	require.Equal(t, testAcc.Currency, acc.Currency)

	require.NotZero(t, acc.ID)
	require.NotZero(t, acc.CreatedAt)

}

func TestQueries_DeleteAcc(t *testing.T) {
	testAcc := createRandAcc(t)
	id := testAcc.ID

	err := testQueries.DeleteAcc(context.Background(), id)

	require.NoError(t, err)

	if _, err = testQueries.GetAcc(context.Background(), id); err != nil {
		fmt.Println("account successfully deleted")
	}
}

func TestQueries_GetAllAccs(t *testing.T) {

	accounts, err := testQueries.GetAllAccs(context.Background())
	require.NoError(t, err)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.NotEmpty(t, account.ID)
		require.NotEmpty(t, account.Owner)
		require.NotEmpty(t, account.Currency)
		require.NotEmpty(t, account.CreatedAt)
	}
}
