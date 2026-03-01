package misc

import (
	"time"

	"github.com/HackUCF/Quincy/api/db/conn"
)

// GetCompDuration returns how long the competition has been running
func GetCompDuration() (time.Duration, error) {
	db := conn.Get()

	var start uint64
	var end uint64

	// get the first and last timestamp in the scores table
	rows, err := db.Query(`
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
