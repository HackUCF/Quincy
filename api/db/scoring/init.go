package scoring

import (
	"fmt"
	"math"

	"github.com/HackUCF/Quincy/api/config"
	"github.com/HackUCF/Quincy/api/db/conn"
	"github.com/HackUCF/Quincy/common/types"
)

// InitScoring creates rows in the final scores table for each team/box/service combo.
func InitScoring(cfg *config.APIConfigSpec) error {
	db := conn.Get()

	query := `
	  INSERT OR IGNORE INTO final_scores (service, box, team_num, total, passed)
		VALUES (?, ?, ?, 0, 0)
	`

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin database transaction: %w", err)
	}

	for _, t := range config.TeamRange {
		for _, box := range cfg.Boxes {
			for _, svc := range box.Services {
				// for every team, box, and service
				_, err := tx.Exec(
					query,
					svc.ID,
					box.ID,
					t,
				)
				if err != nil {
					return fmt.Errorf("could not insert into table: %w", err)
				}
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("could not commit database transaction: %w", err)
	}

	return nil
}

// some helper functions:

// round a float to n decimal places
// crazy this isn't in the standard library
func round(f float64, n int) float64 {
	factor := math.Pow(10, float64(n))
	return math.Round(f*factor) / factor
	// for example: n = 2
	// multiply by 100, round to nearest int, then divide by 100
}

func getResult(total uint64, passed uint64) types.ScoreResult {
	var result types.ScoreResult

	// prevent division by 0
	if total == 0 {
		// the default value is full of zeros, which happens to be correct
		return result
	}

	// calculate uptime
	pcent := float64(passed) / float64(total) * 100

	// return my struct
	return types.ScoreResult{
		TotalChecks:   total,
		ChecksPassed:  passed,
		ChecksFailed:  total - passed,
		UptimePercent: round(pcent, 2),
	}
}
