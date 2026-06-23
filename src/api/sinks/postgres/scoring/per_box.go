package scoring

import (
	"context"
	"fmt"

	"github.com/HackUCF/quincy/common/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

// BoxScores is a map that can be used to quickly reference stats for a specific box.
type BoxScores map[types.BoxName]types.ScoreResult

// GetBoxScores collates and calculates competion scores for all teams.
func GetBoxScores(ctx context.Context, db *pgxpool.Pool) (BoxScores, error) {

	// dict of box id to scoring results
	var scores = make(BoxScores)

	// get totals for each box
	// outputs n rows, where n is the number of boxes
	rows, err := db.Query(ctx, `
    SELECT
    	box,
      SUM(total) as total,
      SUM(passed) as passed
    FROM final_scores
    GROUP BY box
  `)
	if err != nil {
		err = fmt.Errorf("failed to query db: %w", err)
		return scores, err
	}
	defer rows.Close()

	// holding variables, shit gets scanned into them
	var boxID types.BoxName
	var total, passed uint64

	// loop through all of the teams
	for rows.Next() {
		// scan into variables
		err := rows.Scan(&boxID, &total, &passed)
		if err != nil {
			err = fmt.Errorf("failed to scan a row: %w", err)
			return scores, err
		}

		// store in results
		scores[boxID] = getResult(total, passed)
	}

	return scores, nil
}
