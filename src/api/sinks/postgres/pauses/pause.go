package pauses

import (
	"context"
	"fmt"
	"time"

	"github.com/HackUCF/Quincy/src/common/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// IsPaused returns whether or not the competition is currently paused.
// Takes a tx instead of the entire db pool for use in reliable logic.
func IsPaused(ctx context.Context, db *pgxpool.Pool, pauseType types.PauseType) (types.PauseState, error) {

	var current types.PauseState

	row := db.QueryRow(ctx, "SELECT state FROM pauses WHERE type = $1 ORDER BY timestamp DESC LIMIT 1;", pauseType)

	err := row.Scan(&current)
	if err != nil {
		err = fmt.Errorf("failed to get state from db: %w", err)
		return types.PauseDefault, err
	}

	return current, nil
}

// IsPausedTx returns whether or not the competition is currently paused.
// Takes a tx instead of the entire db pool for use in reliable logic.
func IsPausedTx(ctx context.Context, tx pgx.Tx, pauseType types.PauseType) (types.PauseState, error) {

	var current types.PauseState

	row := tx.QueryRow(ctx, "SELECT state FROM pauses WHERE type = $1 ORDER BY timestamp DESC LIMIT 1;", pauseType)

	err := row.Scan(&current)
	if err != nil {
		err = fmt.Errorf("failed to get state from db: %w", err)
		return types.PauseDefault, err
	}

	return current, nil
}

// GetPauseRecord returns the current pause state, as well as the time that state went into effect.
func GetPauseRecord(ctx context.Context, db *pgxpool.Pool) (types.PauseRecord, error) {

	var record = types.PauseRecord{
		ScoringState: types.PauseDefault,
		ScoringSince: time.Time{},
		CheckState:   types.PauseDefault,
		CheckSince:   time.Time{},
	}
	var rawTime int64

	tx, err := db.Begin(ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin transaction: %w", err)
		return record, err
	}
	defer tx.Rollback(ctx)

	// scoring pause
	{
		// get pause record from db
		row := tx.QueryRow(
			ctx,
			"SELECT state, timestamp FROM pauses WHERE type = $1 ORDER BY timestamp DESC LIMIT 1;",
			types.ScoringPause,
		)

		// scan into variables
		err = row.Scan(&record.ScoringSince, &rawTime)
		if err != nil {
			err = fmt.Errorf("failed to get state from db: %w", err)
			return record, err
		}

		record.ScoringSince = time.UnixMicro(rawTime)
	}

	// checks pause
	{
		// get pause record from db
		row := tx.QueryRow(
			ctx,
			"SELECT state, timestamp FROM pauses WHERE type = $1 ORDER BY timestamp DESC LIMIT 1;",
			types.CheckPause,
		)

		// scan into variables
		err = row.Scan(&record.CheckState, &rawTime)
		if err != nil {
			err = fmt.Errorf("failed to get state from db: %w", err)
			return record, err
		}

		record.CheckSince = time.UnixMicro(rawTime)
	}

	// commit transaction
	if err = tx.Commit(ctx); err != nil {
		err = fmt.Errorf("failed to commit transaction: %w", err)
		return record, err
	}

	return record, nil
}
