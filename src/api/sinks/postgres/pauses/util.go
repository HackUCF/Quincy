package pauses

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// tryLock attempts to lock a postgres lock by name
// quickly exits
// lock is held until transaction commit or rollback
func tryLock(ctx context.Context, tx pgx.Tx, lockName string) (got bool, err error) {

	// select a lock by name, returns a sql boolean
	row := tx.QueryRow(ctx, `SELECT pg_try_advisory_xact_lock(hashtext($1))`, lockName)

	// check if we got the lock
	err = row.Scan(&got)

	if err != nil {
		// some db error
		return false, err
	}
	if !got {
		// someone else holds the lock
		return false, nil
	}

	return true, nil
}

// blockLock attempts to lock a postgres lock by name
// lock is held until transaction commit or rollback
// blocks until the lock is acquired
func blockLock(ctx context.Context, tx pgx.Tx, lockName string) error {

	_, err := tx.Exec(ctx, `SELECT pg_advisory_xact_lock(hashtext($1))`, lockName)
	return err
}
