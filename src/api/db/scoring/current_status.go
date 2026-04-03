package scoring

import (
	"database/sql"
	"fmt"

	"github.com/HackUCF/Quincy/api/config"
	"github.com/HackUCF/Quincy/common/types"
)

// GetCurrentServiceStatus returns the entire recent scores table.
// It is sorted by team number, then box ID, then service ID.
func GetCurrentServiceStatus(db *sql.DB, cfg *config.APIConfigSpec) ([]types.Score, error) {

	cap := int(cfg.NumTeams) * len(cfg.Boxes)
	status := make([]types.Score, 0, cap)

	rows, err := db.Query(`
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
