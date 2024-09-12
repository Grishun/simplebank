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

func createEntry(t *testing.T, AccId, amount int64) (Entry, error) {

	input := NewEntryParams{
		AccID: sql.NullInt64{
			Int64: AccId,
			Valid: true,
		},
		Amount: amount,
	}
	entry, err := testQueries.NewEntry(context.Background(), input)
	if err != nil {
		return Entry{}, err
	}
	require.NotEmpty(t, entry)

	require.Equal(t, AccId, entry.AccID.Int64)
	require.Equal(t, amount, entry.Amount)

	require.NotEmpty(t, entry.ID)
	require.NotEmpty(t, entry.CreatedAt)

	return entry, nil
}

func TestQueries_NewEntry(t *testing.T) {
	testCases := []struct {
		id       int64
		amount   int64
		expected *pq.Error
	}{
		{
			id:       1,
			amount:   util.RandInt(0, 1000),
			expected: nil,
		},
		{
			id:       2,
			amount:   -util.RandInt(0, 1000),
			expected: nil,
		},
		{
			id:       -1,
			amount:   util.RandInt(0, 1000),
			expected: &pq.Error{},
		},
		{
			id:       0,
			amount:   -util.RandInt(0, 1000),
			expected: &pq.Error{},
		},
	}

	for i, testCase := range testCases {
		t.Run("case_"+strconv.Itoa(i), func(t *testing.T) {
			_, err := createEntry(t, testCase.id, testCase.amount)
			if testCase.expected != nil {
				require.NotEmpty(t, err)
			}
		})
	}

}

func TestQueries_GetEntry(t *testing.T) {

	accounts, err := testQueries.GetAllAccs(context.Background())

	for _, acc := range accounts {
		testEntry, _ := createEntry(t, acc.ID, util.RandInt(1, 100))
		require.NoError(t, err)

		receivedEntry, err := testQueries.GetEntry(context.Background(), testEntry.ID)
		require.NoError(t, err)

		require.Equal(t, testEntry, receivedEntry)
	}
}

func TestQueries_GetAllEntries(t *testing.T) {

	testAcc := createRandAcc(t)

	entries := make([]Entry, 0, 10)

	for i := 0; int64(i) < util.RandInt(0, 10); i++ {
		entry, err := createEntry(t, testAcc.ID, util.RandInt(1, 100))
		require.NoError(t, err)
		entries = append(entries, entry)
	}

	allEntries, err := testQueries.GetAllEntries(context.Background(), sql.NullInt64{testAcc.ID, true})
	require.NoError(t, err)
	require.Equal(t, entries, allEntries)

}
