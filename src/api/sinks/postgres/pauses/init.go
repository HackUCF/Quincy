package pauses

import (
	"context"
	"fmt"
	"time"

	"github.com/HackUCF/Quincy/src/api/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InitPauses sets the pause state in the db to the configured default when none exists.
func InitPauses(ctx context.Context, db *pgxpool.Pool, cfg *config.APIConfigSpec) error {

	ts := time.Now().UnixMicro()

	// insert the configured pause start configuration
	// the NOT EXISTS guard means this is a no-op post first boot
	_, err := db.Exec(
		ctx,
		"INSERT INTO pause_states (timestamp, state) SELECT $1, $2 WHERE NOT EXISTS (SELECT 1 FROM pause_states)",
		ts, cfg.StartPaused,
	)
	if err != nil {
		err = fmt.Errorf("failed to initialize pause state: %w", err)
		return err
	}

	return nil
}
