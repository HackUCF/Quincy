package scoring

import (
	"context"
	"fmt"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/common/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetCurrentServiceStatus returns the entire recent scores table.
// It is sorted by team number, then box ID, then service ID.
func GetCurrentServiceStatus(ctx context.Context, db *pgxpool.Pool, cfg *config.APIConfigSpec) ([]types.Score, error) {

	cap := int(cfg.NumTeams) * len(cfg.Boxes)
	status := make([]types.Score, 0, cap)

	rows, err := db.Query(ctx, `
    SELECT service, box, team_num, status, message, timestamp
    FROM recent_scores
    ORDER BY team_num, box, service
  `)
	if err != nil {
		err = fmt.Errorf("failed to query db: %w", err)
		return status, err
	}

	for rows.Next() {
		var s types.Score
		rows.Scan(&s.ServiceName, &s.BoxName, &s.TeamNum, &s.Status, &s.Message, &s.Timestamp)
		status = append(status, s)
	}

	return status, nil
}
