package pauses

import (
	"context"
	"fmt"
	"time"

	"github.com/HackUCF/Quincy/src/common/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreatePause(ctx context.Context, db *pgxpool.Pool, start, end time.Time) error {

	// start a transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		err = fmt.Errorf("failed to begin db transaction: %w", err)
		return err
	}
	defer tx.Rollback(ctx)

	// insert pause start and end
	tag, err := tx.Exec(
		ctx,
		"INSERT INTO pause_states (timestamp, state) VALUES (?, ?), (?, ?);",
		start.UnixMicro(), types.Paused,
		end.UnixMicro(), types.Unpaused,
	)
	if err != nil {
		err = fmt.Errorf("failed to add new pause: %w", err)
		return err
	}
	if tag.RowsAffected() == 0 {
		err = fmt.Errorf("failed to add new pause: no rows affected")
		return err
	}

	// commit transaction
	if err = tx.Commit(ctx); err != nil {
		err = fmt.Errorf("failed to commit transaction: %w", err)
		return err
	}

	return nil
}
