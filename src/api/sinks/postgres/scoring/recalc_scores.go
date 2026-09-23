package scoring

import (
	"context"

	"github.com/HackUCF/Quincy/src/api/config"
	"github.com/jackc/pgx/v5"
)

// RedoFinalScores udpates the final scores table retroactively, ensuring pauses in scoring are accurately represented.
func RedoFinalScores(ctx context.Context, tx pgx.Tx, cfg *config.APIConfigSpec) (changed bool, err error) {

	// for _, teamNum := range config.TeamRange {

	// 	for _, box := range cfg.Boxes {

	// 		for _, svc := range box.Services {

	// 			_, err := tx.Exec(ctx, `
	// 				UPDATE final_scores
	// 				SET
	// 			`, teamNum, svc)
	// 			if err != nil {
	// 				return fmt.Errorf("failed to recalculate final scores")
	// 			}
	// 		}
	// 	}
	// }

	return changed, nil
}
