package misc

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// GetCompDuration returns how long the competition has been running
func GetCompDuration(ctx context.Context, db *pgxpool.Pool) (time.Duration, error) {

	var start uint64
	var end uint64

	// get the first and last timestamp in the scores table
	rows, err := db.Query(ctx, `
		SELECT * FROM (SELECT timestamp FROM scores ORDER BY timestamp ASC LIMIT 1)
		UNION ALL
		SELECT * FROM (SELECT timestamp FROM scores ORDER BY timestamp DESC LIMIT 1);
	`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	// scan rows
	rows.Next()
	err = rows.Scan(&start)
	if err != nil {
		return 0, err
	}
	rows.Next()
	err = rows.Scan(&end)
	if err != nil {
		return 0, err
	}

	// convert to duration and return
	microSeconds := time.Duration(end-start) * time.Microsecond
	return microSeconds, nil
}
