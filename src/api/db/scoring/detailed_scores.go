package scoring

import (
	"database/sql"
	"fmt"

	"github.com/HackUCF/quincy/common/types"

	_ "github.com/mattn/go-sqlite3"
)

// DetailedScores is a double-nested map that can be used to quickly reference stats for a specific teams' boxes' services.
// EXTREMELY ugly but i cant think of a better way to do this.
// this is a map[]map[]map[]. team_num -> box -> service.
type DetailedScores map[types.TeamNum]ServiceScores

// GetDetailedScores pulls the full final scoring data as it is stored in the db.
func GetDetailedScores(db *sql.DB) (DetailedScores, error) {
	// ugly nested dictionary
	var scores = make(DetailedScores)

	// get every row from db
	rows, err := db.Query(`
    SELECT
      team_num,
    	box,
     	service,
      total,
      passed
    FROM final_scores
  `)
	if err != nil {
		err = fmt.Errorf("failed to query db: %w", err)
		return scores, err
	}
	defer rows.Close()

	// holding variables, shit gets scanned into them
	var teamNum types.TeamNum
	var boxID types.BoxName
	var serviceID types.ServiceName
	var total, passed uint64

	// loop through all of the teams
	for rows.Next() {
		// scan into variables
		err := rows.Scan(&teamNum, &boxID, &serviceID, &total, &passed)
		if err != nil {
			err = fmt.Errorf("failed to scan a row: %w", err)
			return scores, err
		}

		// store in results (stupid complicated)

		// create map for this team if it doesn't exist yet
		if _, ok := scores[teamNum]; !ok {
			scores[teamNum] = make(ServiceScores)
		}

		// create map for this box if it doesn't exist yet
		if _, ok := scores[teamNum][boxID]; !ok {
			scores[teamNum][boxID] = make(map[types.ServiceName]types.ScoreResult)
		}

		// insert the result 😭
		scores[teamNum][boxID][serviceID] = getResult(total, passed)
	}

	return scores, nil
}
