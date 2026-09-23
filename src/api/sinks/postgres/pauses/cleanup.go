package pauses

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// deletes runs of two of the same pause state in a row (always deletes the second)
//
// the logic (not running in this order, but this is what it's doing)
// 1. it selects a and b, where b is after a.
// 2. it compares the state, they must both be pauses or both be unpauses.
// 3. it checks that there's no c such that c is between a and b (verifies that these are two in a row)
// 4. it deletes every b -> the duplicate/second.
// 5. will also catch 3, 4, ..., n in a row, because the 5th duplicates the 4th, etc.
const cleanupSql = `
	DELETE FROM pause_states
	WHERE id IN (
		SELECT b.id
		FROM pause_states a
		JOIN pause_states b
			ON b.timestamp > a.timestamp
		AND b.state = a.state
		WHERE NOT EXISTS (
			SELECT 1
			FROM pause_states c
			WHERE c.timestamp > a.timestamp
				AND c.timestamp < b.timestamp
		)
);`

// CleanupPauses ensures that the pauses table is in a tenable spot.
// Run as part of a transaction to verify that future references to it are valid.
func CleanupPauses(ctx context.Context, tx pgx.Tx) error {

	_, err := tx.Exec(ctx, cleanupSql)
	if err != nil {
		err := fmt.Errorf("failed to execute pause cleanup: %w", err)
		return err
	}

	return nil
}
