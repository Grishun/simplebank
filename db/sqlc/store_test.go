package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore_MakeTransfer(t *testing.T) {
	store := NewStore(testDB)

	fromAcc, toAcc := createRandAcc(t), createRandAcc(t)
	fmt.Println(">>> start:", fromAcc.Balance, toAcc.Balance)
	errs, results := make(chan error), make(chan TxTransResult)
	unique1 := make(map[int64]struct{})

	n, amount := int64(10), int64(100)

	for i := 0; int64(i) < n; i++ {

		go func() {
			res, err := store.MakeTransfer(context.Background(), TxTransParams{
				FromAccId: fromAcc.ID,
				ToAccId:   toAcc.ID,
				Amount:    amount,
			})
			errs <- err
			results <- res
		}()
	}

	for i := 0; int64(i) < n; i++ {
		err := <-errs
		require.NoError(t, err)
		res := <-results

		//check the accounts
		require.NotEmpty(t, res)
		require.Equal(t, fromAcc.ID, res.FromAcc.ID)
		require.Equal(t, toAcc.ID, res.ToAcc.ID)

		// check the transfer
		require.NotEmpty(t, res.Transfer.ID)
		require.NotEmpty(t, res.Transfer.CreatedAt)
		require.Equal(t, amount, res.Transfer.Amount)
		require.Equal(t, fromAcc.ID, res.Transfer.FromAccID.Int64)
		require.Equal(t, toAcc.ID, res.Transfer.ToAccID.Int64)

		//check the entry of fromAcc entries
		require.NotEmpty(t, res.FromEntry.ID)
		require.Equal(t, fromAcc.ID, res.FromEntry.AccID.Int64)
		require.Equal(t, -amount, res.FromEntry.Amount)
		require.NotEmpty(t, res.FromEntry.ID)
		require.NotEmpty(t, res.FromEntry.CreatedAt)

		//check the entry of ToAcc entries
		require.NotEmpty(t, res.ToEntry.ID)
		require.Equal(t, toAcc.ID, res.ToEntry.AccID.Int64)
		require.Equal(t, amount, res.ToEntry.Amount)
		require.NotEmpty(t, res.ToEntry.ID)
		require.NotEmpty(t, res.ToEntry.CreatedAt)

		// check the balance updating
		fmt.Println(">>> tx:", res.FromAcc.Balance, res.ToAcc.Balance)
		diff1 := fromAcc.Balance - res.FromAcc.Balance
		require.True(t, diff1%amount == 0)

		k := diff1 / amount
		require.True(t, k >= 1 && k <= n)
		_, ok1 := unique1[k]
		if !ok1 {
			unique1[k] = struct{}{}
		} else {
			t.Errorf("problem in balance changing in fromAcc")
		}

		diff2 := res.ToAcc.Balance - toAcc.Balance
		require.Equal(t, diff1, diff2)

	}
	updFromAcc, _ := store.GetAcc(context.Background(), fromAcc.ID)
	updToAcc, _ := store.GetAcc(context.Background(), toAcc.ID)
	fmt.Println(">>> final:", updFromAcc.Balance, updToAcc.Balance)
	// check the final balance of the accounts
	require.True(t, fromAcc.Balance-n*amount == updFromAcc.Balance)
	require.True(t, toAcc.Balance+n*amount == updToAcc.Balance)
}

func TestStore_TransferDeadlock(t *testing.T) {
	store := NewStore(testDB)

	fromAcc, toAcc := createRandAcc(t), createRandAcc(t)
	fmt.Println(">>> start:", fromAcc.Balance, toAcc.Balance)
	errs := make(chan error)

	n, amount := int64(10), int64(100)

	for i := 0; int64(i) < n; i++ {
		fromAccId, toAccId := fromAcc.ID, toAcc.ID

		if i%2 == 0 {
			fromAccId, toAccId = toAcc.ID, fromAcc.ID
		}
		go func() {
			_, err := store.MakeTransfer(context.Background(), TxTransParams{
				FromAccId: fromAccId,
				ToAccId:   toAccId,
				Amount:    amount,
			})
			errs <- err

		}()
	}

	for i := 0; int64(i) < n; i++ {
		err := <-errs
		require.NoError(t, err)

		//check the accounts

	}
	updFromAcc, _ := store.GetAcc(context.Background(), fromAcc.ID)
	updToAcc, _ := store.GetAcc(context.Background(), toAcc.ID)
	fmt.Println(">>> final:", updFromAcc.Balance, updToAcc.Balance)
	// check the final balance of the accounts
	require.Equal(t, updToAcc.Balance, toAcc.Balance)
	require.Equal(t, updFromAcc.Balance, fromAcc.Balance)
}
