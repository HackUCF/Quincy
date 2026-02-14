package scoring

import (
	"fmt"

	"github.com/HackUCF/Quincy/api/config"
	"github.com/HackUCF/Quincy/api/db/conn"
	"github.com/HackUCF/Quincy/common/types"
)

// GetCurrentServiceStatus returns the entire recent scores table.
// It is sorted by team number, then box ID, then service ID.
func GetCurrentServiceStatus() ([]types.Score, error) {
	c, err := config.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	cap := int(c.NumTeams) * len(c.Boxes)
	status := make([]types.Score, 0, cap)
	db := conn.Get()

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
		rows.Scan(&s.ServiceID, &s.BoxID, &s.TeamNum, &s.Status, &s.Message, &s.Timestamp)
		status = append(status, s)
	}

	return status, nil
}
