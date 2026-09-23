package pauses

import (
	"context"
	"fmt"
	"time"

	"github.com/HackUCF/Quincy/src/common/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ChangePauseState returns whether or not the competition is currently paused.
func ChangePauseState(ctx context.Context, db *pgxpool.Pool, desiredState types.PauseState) (changed bool, err error) {

	// start a transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin db transaction: %w", err)
		return false, err
	}
	defer tx.Rollback(ctx)

	// check current state
	currentState, err := IsPaused(ctx, tx)
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
	tag, err := tx.Exec(ctx, "INSERT INTO pause_states (timestamp, state) VALUES (?, ?);", ts, desiredState)
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
