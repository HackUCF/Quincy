package scoring

import (
	"context"
	"fmt"

	"github.com/HackUCF/Quincy/src/api/config"
	"github.com/HackUCF/Quincy/src/common/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

// recentScoresQuery pulls the newest numChecks rows for one team's single service.
// It is run once per service so that every service gets its own numChecks,
// even though services are not checked in lockstep.
//
// The equality columns and the sort line up with idx_scores_team_box_service_ts,
// so this is an index top-N: it reads numChecks rows, not the service's history.
const recentScoresQuery = `
  SELECT status, stdout, stderr, timestamp
  FROM scores
  WHERE team_num = $1
    AND box      = $2
    AND service  = $3
  ORDER BY timestamp DESC
  LIMIT $4
`

// GetRecentNScores collects the numChecks most recent scores for every service
// on every box belonging to a single team.
// Results are sorted by box, then service, then newest check first.
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
	capacity *= numChecks // multiply by the number needed

	// create output slice
	recentScores := make([]types.Score, 0, capacity)

	// one query per service, walked in config order so the results come back
	// already sorted by box then service
	for _, box := range cfg.Boxes {
		for _, service := range box.Services {
			err := appendServiceScores(ctx, db, &recentScores, teamNum, box.Name, service.Name, numChecks)
			if err != nil {
				return recentScores, err
			}
		}
	}

	return recentScores, nil
}

// appendServiceScores runs one service's query and appends its rows to dst.
// Split out from the loop above so that the rows can be closed by defer:
// deferring inside the loop would hold every connection until the whole
// request finished, which deadlocks the pool once services outnumber conns.
func appendServiceScores(
	ctx context.Context,
	db *pgxpool.Pool,
	dst *[]types.Score,
	teamNum types.TeamNum,
	boxName types.BoxName,
	serviceName types.ServiceName,
	numChecks int,
) error {
	rows, err := db.Query(ctx, recentScoresQuery, teamNum, boxName, serviceName, numChecks)
	if err != nil {
		return fmt.Errorf("failed to query db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		s := types.Score{
			TeamNum:     teamNum,
			BoxName:     boxName,
			ServiceName: serviceName,
		}
		if err := rows.Scan(&s.Status, &s.Stdout, &s.Stderr, &s.Timestamp); err != nil {
			return fmt.Errorf("failed to scan recent score row: %w", err)
		}
		*dst = append(*dst, s)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("failed to read recent score rows: %w", err)
	}

	return nil
}
