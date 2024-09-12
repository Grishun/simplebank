package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	if err := fn(q); err != nil {
		if txErr := tx.Rollback(); err != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, txErr)
		}
		return err
	}

	return tx.Commit()
}

type TxTransParams struct {
	FromAccId int64 `json:"from_acc_id"`
	ToAccId   int64 `json:"to_acc_id"`
	Amount    int64 `json:"amount"`
}

type TxTransResult struct {
	FromAcc   Account  `json:"from_acc"`
	ToAcc     Account  `json:"to_acc"`
	FromEntry Entry    `json:"from_entry"`
	ToEntry   Entry    `json:"to_entry"`
	Transfer  Transfer `json:"transfer"`
}

var txKey = struct{}{}

func (s *Store) MakeTransfer(ctx context.Context, params TxTransParams) (res TxTransResult, err error) {

	err = s.execTx(ctx, func(queries *Queries) error {

		txName := ctx.Value(txKey)

		log.Println(txName, "create transfer")
		res.Transfer, err = queries.NewTransfer(ctx, NewTransferParams{
			FromAccID: sql.NullInt64{params.FromAccId, true},
			ToAccID:   sql.NullInt64{params.ToAccId, true},
			Amount:    params.Amount,
		})
		if err != nil {
			return err
		}

		log.Println(txName, "create fromEntry")
		res.FromEntry, err = queries.NewEntry(ctx, NewEntryParams{
			AccID:  sql.NullInt64{params.FromAccId, true},
			Amount: -params.Amount,
		})
		if err != nil {
			return err
		}

		log.Println(txName, "create toEntry")
		res.ToEntry, err = queries.NewEntry(ctx, NewEntryParams{
			AccID:  sql.NullInt64{params.ToAccId, true},
			Amount: params.Amount,
		})
		if err != nil {
			return err
		}

		// get acc -> update its balance

		if params.FromAccId < params.ToAccId {
			res.FromAcc, res.ToAcc, err = updBalance(ctx, queries, params.FromAccId, params.ToAccId, -params.Amount, params.Amount)
		} else {
			res.FromAcc, res.ToAcc, err = updBalance(ctx, queries, params.ToAccId, params.FromAccId, params.Amount, -params.Amount)
		}

		return err
	})

	return res, err
}

func updBalance(ctx context.Context, q *Queries, AccId1, AccId2, amount1, amount2 int64) (Acc1, Acc2 Account, err error) {
	Acc1, err = q.UpdateAccBalance(ctx, UpdateAccBalanceParams{
		Amount: amount1,
		ID:     AccId1,
	})
	if err != nil {
		return
	}
	Acc2, err = q.UpdateAccBalance(ctx, UpdateAccBalanceParams{
		Amount: amount2,
		ID:     AccId2,
	})
	if err != nil {
		return
	}

	return
}
