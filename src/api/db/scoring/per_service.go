package scoring

import (
	"database/sql"
	"fmt"

	"github.com/HackUCF/Quincy/common/types"

	_ "github.com/mattn/go-sqlite3"
)

// ServiceScores is a nested map that can be used to quickly reference stats for a specific boxes' specific service.
// very ugly but i cant think of a better way to do this.
type ServiceScores map[types.BoxName]map[types.ServiceName]types.ScoreResult

// GetServiceScores collates and calculates competion scores for each service on each box.
func GetServiceScores(db *sql.DB) (ServiceScores, error) {

	// ugly nested dictionary
	var scores = make(ServiceScores)

	// get totals for each box+service combo
	rows, err := db.Query(`
    SELECT
    	box,
     	service,
      SUM(total) as total,
      SUM(passed) as passed
    FROM final_scores
    GROUP BY box, service
  `)
	if err != nil {
		err = fmt.Errorf("failed to query db: %w", err)
		return scores, err
	}
	defer rows.Close()

	// holding variables, shit gets scanned into them
	var boxID types.BoxName
	var serviceID types.ServiceName
	var total, passed uint64
	// loop through all of the teams
	for rows.Next() {
		// scan into variables
		err := rows.Scan(&boxID, &serviceID, &total, &passed)
		if err != nil {
			err = fmt.Errorf("failed to scan a row: %w", err)
			return scores, err
		}

		// store in results (complicated)
		// create map if it doesn't exist yet
		if _, ok := scores[boxID]; !ok {
			scores[boxID] = make(map[types.ServiceName]types.ScoreResult)
		}
		scores[boxID][serviceID] = getResult(total, passed)
	}

	return scores, nil
}
