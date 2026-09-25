package pauses

import (
	"context"
	"fmt"
	"time"

	"github.com/HackUCF/Quincy/src/api/config"
	"github.com/HackUCF/Quincy/src/common/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

const pauseLock = "pauses"

// InitPauses sets the pause state in the db to the configured default when none exists.
func InitPauses(ctx context.Context, db *pgxpool.Pool, cfg *config.APIConfigSpec) error {

	// start a transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		err := fmt.Errorf("failed to begin transaction while initializing pauses")
		return err
	}
	defer tx.Rollback(ctx)

	// prevent double inits
	got, err := tryLock(ctx, tx, pauseLock)
	if !got {
		// someone else is messing with pauses. they got it.
		return nil
	}

	ts := time.Now().UnixMicro()

	// insert the configured pause start configuration
	// the NOT EXISTS guard means this is a no-op post first boot
	_, err = tx.Exec(
		ctx,
		"INSERT INTO pauses (timestamp, state, type) SELECT $1, $2, $3 WHERE NOT EXISTS (SELECT 1 FROM pauses WHERE type = $3)",
		ts, cfg.StartPaused, types.CheckPause,
	)
	if err != nil {
		err = fmt.Errorf("failed to initialize %q pause state: %w", types.CheckPause, err)
		return err
	}

	_, err = tx.Exec(
		ctx,
		"INSERT INTO pauses (timestamp, state, type) SELECT $1, $2, $3 WHERE NOT EXISTS (SELECT 1 FROM pauses WHERE type = $3)",
		ts, cfg.StartPaused, types.ScoringPause,
	)
	if err != nil {
		err = fmt.Errorf("failed to initialize %q pause state: %w", types.ScoringPause, err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		err = fmt.Errorf("failed to initialize pause state: %w", err)
		return err
	}

	return nil
}
