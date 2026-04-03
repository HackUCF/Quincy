package scoring

import (
	"database/sql"
	"fmt"

	"github.com/HackUCF/Quincy/common/types"

	_ "github.com/mattn/go-sqlite3"
)

// TeamScores is a map that can be used to quickly reference a specific teams stats.
type TeamScores map[types.TeamNum]types.ScoreResult

// GetTeamScores collates and calculates competion scores for all teams.
func GetTeamScores(db *sql.DB) (TeamScores, error) {

	// dict of team number to scoring results
	var scores = make(TeamScores)

	// get totals for each team
	// outputs n rows, where n is the number of teams
	rows, err := db.Query(`
    SELECT
    	team_num,
      SUM(total) as total,
      SUM(passed) as passed
    FROM final_scores
    GROUP BY team_num
  `)
	if err != nil {
		err = fmt.Errorf("failed to query db: %w", err)
		return scores, err
	}
	defer rows.Close()

	// holding variables, shit gets scanned into them
	var teamNum types.TeamNum
	var total, passed uint64

	// loop through all of the teams
	for rows.Next() {
		// scan into variables
		err := rows.Scan(&teamNum, &total, &passed)
		if err != nil {
			err = fmt.Errorf("failed to scan a row: %w", err)
			return scores, err
		}

		// store in results
		scores[teamNum] = getResult(total, passed)
	}

	return scores, nil
}
