package db

import (
	"context"
	"database/sql"
	"github.com/Grishun/simplebank/db/util"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func createTransfer(t *testing.T, amount int64) Transfer {
	fromAcc, toAcc := createRandAcc(t), createRandAcc(t)

	args := NewTransferParams{
		FromAccID: sql.NullInt64{fromAcc.ID, true},
		ToAccID:   sql.NullInt64{toAcc.ID, true},
		Amount:    amount,
	}

	tran, err := testQueries.NewTransfer(context.Background(), args)
	require.NoError(t, err)

	return tran
}

func TestQueries_NewTransfer(t *testing.T) {

	fromAcc, toAcc := createRandAcc(t), createRandAcc(t)

	testCases := []struct {
		input       NewTransferParams
		expectedErr *pq.Error
	}{
		{
			input: NewTransferParams{
				FromAccID: sql.NullInt64{fromAcc.ID, true},
				ToAccID:   sql.NullInt64{toAcc.ID, true},
				Amount:    util.RandInt(1, 100),
			},
			expectedErr: nil,
		}, {
			input: NewTransferParams{
				FromAccID: sql.NullInt64{0, true},
				ToAccID:   sql.NullInt64{toAcc.ID, true},
				Amount:    util.RandInt(1, 100),
			},
			expectedErr: &pq.Error{},
		},
		{
			input: NewTransferParams{
				FromAccID: sql.NullInt64{0, true},
				ToAccID:   sql.NullInt64{0, true},
				Amount:    util.RandInt(1, 100),
			},
			expectedErr: &pq.Error{},
		},
		{
			input: NewTransferParams{
				FromAccID: sql.NullInt64{toAcc.ID, true},
				ToAccID:   sql.NullInt64{0, true},
				Amount:    util.RandInt(1, 100),
			},
			expectedErr: &pq.Error{},
		},
		{
			input: NewTransferParams{
				FromAccID: sql.NullInt64{fromAcc.ID, true},
				ToAccID:   sql.NullInt64{toAcc.ID, true},
				Amount:    util.RandInt(-100, -1),
			},
			expectedErr: &pq.Error{},
		},
	}

	for i, testCase := range testCases {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
			testTrans, err := testQueries.NewTransfer(context.Background(), testCase.input)
			if testCase.expectedErr != nil {
				require.NotZero(t, err)
				return
			}
			require.NotEmpty(t, testTrans)
			require.Equal(t, toAcc.ID, testTrans.ToAccID.Int64)
			require.Equal(t, fromAcc.ID, testTrans.FromAccID.Int64)
			require.NotEmpty(t, testTrans.ID)
			require.NotEmpty(t, testTrans.CreatedAt)
			require.NotEmpty(t, testTrans.Amount)
		})
	}
}

func TestQueries_GetTransfer(t *testing.T) {

	testTrans := createTransfer(t, util.RandInt(1, 100))

	received, err1 := testQueries.GetTransfer(context.Background(), testTrans.ID)
	require.NoError(t, err1)
	require.Equal(t, testTrans, received)

	_, err2 := testQueries.GetTransfer(context.Background(), 0)
	require.Error(t, err2)
}

func TestQueries_GetAllTransfers(t *testing.T) {
	trans := make([]Transfer, 0, 10)
	toAcc, fromAcc := createRandAcc(t), createRandAcc(t)
	arg := NewTransferParams{
		FromAccID: sql.NullInt64{fromAcc.ID, true},
		ToAccID:   sql.NullInt64{toAcc.ID, true},
		Amount:    util.RandInt(1, 100),
	}
	for i := 0; int64(i) < util.RandInt(0, 10); i++ {
		tr, _ := testQueries.NewTransfer(context.Background(), arg)
		trans = append(trans, tr)
	}

	params := GetAllTransfersParams{
		FromAccID: arg.FromAccID,
		ToAccID:   arg.ToAccID,
	}
	received, err := testQueries.GetAllTransfers(context.Background(), params)
	require.NoError(t, err)
	require.Equal(t, trans, received)
}
