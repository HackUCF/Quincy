package pauses

/*
	i'm stubbing this code

	it's a really cool idea, let the organizers retroactively add pauses to the competition.
	the hard part is the final_scores table. it would need to be updated/tainted/recalculated

	it could be expanded to also take a box, service, and/or team number.
	then you could exempt certain boxes/checks for certain periods.
	AND you'd get to view/audit all those actions bc they're stored in the db.
*/

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

	// ensure the pause table is in a valid state
	err = CleanupPauses(ctx, tx)
	if err != nil {
		return err
	}

	// insert pause start and end
	tag, err := tx.Exec(
		ctx,
		"INSERT INTO pause_states (timestamp, state) VALUES ($1, $2), ($3, $4);",
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

	// update the final scores table

	// commit transaction
	if err = tx.Commit(ctx); err != nil {
		err = fmt.Errorf("failed to commit transaction: %w", err)
		return err
	}

	return nil
}
