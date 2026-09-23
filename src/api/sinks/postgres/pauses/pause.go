package pauses

import (
	"context"
	"fmt"
	"time"

	"github.com/HackUCF/Quincy/src/common/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// IsPausedTx returns whether or not the competition is currently paused.
// Takes a tx instead of the entire db pool for use in reliable logic.
func IsPaused(ctx context.Context, tx pgx.Tx) (types.PauseState, error) {

	var current types.PauseState

	row := tx.QueryRow(ctx, "SELECT state FROM pause_states ORDER BY timestamp DESC LIMIT 1;")

	err := row.Scan(&current)
	if err != nil {
		err = fmt.Errorf("failed to get state from db: %w", err)
		return types.PauseDefault, err
	}

	return current, nil
}

func GetPauseRecord(ctx context.Context, db *pgxpool.Pool) (types.PauseRecord, error) {

	var record = types.PauseRecord{
		State: types.PauseDefault,
		Since: time.Time{},
	}
	var rawTime int64

	// get pause record from db
	row := db.QueryRow(ctx, "SELECT state, timestamp FROM pause_states ORDER BY timestamp DESC LIMIT 1;")

	// scan into variable
	err := row.Scan(&record.State, &rawTime)
	if err != nil {
		err = fmt.Errorf("failed to get state from db: %w", err)
		return record, err
	}

	record.Since = time.UnixMicro(rawTime)

	return record, nil
}
