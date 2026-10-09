package scoring

import (
	"context"
	"fmt"

	"github.com/HackUCF/Quincy/src/api/config"
	"github.com/HackUCF/Quincy/src/common/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetRecentNScores collects
func GetRecentNScores(
	ctx context.Context,
	db *pgxpool.Pool,
	cfg *config.APIConfigSpec,
	teamNum types.TeamNum,
	numChecks int,
) ([]types.Score, error) {

	// get final list capacity
	var capacity int = 0
	for _, box := range cfg.Boxes {
		capacity += len(box.Services) // count all services
	}
	capacity *= numChecks // multiple by the number needed

	// create output slice
	recentScores := make([]types.Score, 0, capacity)

	rows, err := db.Query(ctx, `
    SELECT service, box, status, stdout, stderr, timestamp
    FROM scores
    ORDER BY timestamp DESC, box, service
		WHERE team_num = $1
		LIMIT $2
  `, teamNum, numChecks)
	if err != nil {
		err = fmt.Errorf("failed to query db: %w", err)
		return recentScores, err
	}

	for rows.Next() {
		var s types.Score
		s.TeamNum = teamNum
		if err := rows.Scan(&s.ServiceName, &s.BoxName, &s.Status, &s.Stdout, &s.Stderr, &s.Timestamp); err != nil {
			err = fmt.Errorf("failed to scan recent score row: %w", err)
			return recentScores, err
		}
		recentScores = append(recentScores, s)
	}
	return nil, nil

}
