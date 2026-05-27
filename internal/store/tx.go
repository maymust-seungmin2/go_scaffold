// Package store provides database storage operations for the application.
package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/oops"
)

// WithinTx runs fn inside a transaction, committing on success and rolling back
// on error. A rollback failure is joined with the original error so neither is
// lost; pgx.ErrTxClosed is ignored since it only means the tx already ended.
func WithinTx(ctx context.Context, db *pgxpool.Pool, fn func(pgx.Tx) error) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return oops.In("store").Tags("database").Wrapf(err, "begin transaction")
	}

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
			return errors.Join(err, oops.In("store").Tags("database").Wrapf(rbErr, "rollback transaction"))
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return oops.In("store").Tags("database").Wrapf(err, "commit transaction")
	}
	return nil
}
