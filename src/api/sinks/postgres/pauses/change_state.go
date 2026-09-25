package pauses

import (
	"context"
	"fmt"
	"time"

	"github.com/HackUCF/Quincy/src/common/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ChangePauseState pauses or unpauses the competition.
func ChangePauseState(
	ctx context.Context,
	db *pgxpool.Pool,
	desiredState types.PauseState,
	desiredType types.PauseType,
) (changed bool, err error) {

	// start a transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin db transaction: %w", err)
		return false, err
	}
	defer tx.Rollback(ctx)

	// lock the pauses table
	err = blockLock(ctx, tx, pauseLock)
	if err != nil {
		return false, err
	}

	// check current state
	currentState, err := IsPausedTx(ctx, tx, desiredType)
	if err != nil {
		err = fmt.Errorf("failed to check current pause state: %w", err)
		return false, err
	}

	if currentState == desiredState {
		// exit early if no changes required
		return false, nil
	}

	ts := time.Now().UnixMicro()

	// insert new state
	tag, err := tx.Exec(
		ctx,
		"INSERT INTO pauses (timestamp, state, type) VALUES ($1, $2, $3);",
		ts,
		desiredState,
		desiredType,
	)
	if err != nil {
		err = fmt.Errorf("failed to update pause state: %w", err)
		return false, err
	}
	if tag.RowsAffected() == 0 {
		err = fmt.Errorf("failed to update pause state: no rows affected")
		return false, err
	}

	// commit transaction
	if err = tx.Commit(ctx); err != nil {
		err = fmt.Errorf("failed to commit transaction: %w", err)
		return false, err
	}

	return true, nil
}
