package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/HackUCF/quincy/common/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AddScore inserts a score into the main score database.
// It also updates the entry for this specific box/service/team combo in the recent scores table.
func AddScore(ctx context.Context, db *pgxpool.Pool, score types.Score) error {
	ts := time.Now().UnixMicro()

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin database transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `
		INSERT INTO scores (service, box, team_num, status, message, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, score.ServiceName, score.BoxName, score.TeamNum, score.Status, score.Message, ts)
	if err != nil {
		return fmt.Errorf("failed to insert into scores table: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("failed to insert into scores table: no rows affected")
	}

	tag, err = tx.Exec(ctx, `
		INSERT INTO recent_scores (service, box, team_num, status, message, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (team_num, service, box) DO UPDATE SET
			status    = EXCLUDED.status,
			message   = EXCLUDED.message,
			timestamp = EXCLUDED.timestamp
	`, score.ServiceName, score.BoxName, score.TeamNum, score.Status, score.Message, ts)
	if err != nil {
		return fmt.Errorf("failed to insert into recent_scores table: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("failed to upsert into recent_scores table: no rows affected")
	}

	tag, err = tx.Exec(ctx, `
		UPDATE final_scores
		SET total  = total + 1,
		    passed = passed + CASE WHEN $1 THEN 1 ELSE 0 END
		WHERE service = $2 AND box = $3 AND team_num = $4
	`, score.Status, score.ServiceName, score.BoxName, score.TeamNum)
	if err != nil {
		return fmt.Errorf("failed to update final_scores table: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no matching row in final_scores for service=%s, box=%s, team_num=%d",
			score.ServiceName, score.BoxName, score.TeamNum)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
