package scoring

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/HackUCF/quincy/common/types"
)

// AddScore inserts a score into the main score database.
// It also updates the entry for this specific box/service/team combo in the recent scores table.
func AddScore(db *sql.DB, score types.Score) error {
	ts := time.Now().UnixMicro()

	tx, err := db.Begin()
	if err != nil {
		err = fmt.Errorf("failed to begin database transaction: %w", err)
		return err
	}
	defer tx.Rollback()

	// insert the total scores table
	{
		stmt, err := tx.Prepare(`
		  INSERT INTO scores (service, box, team_num, status, message, timestamp)
			VALUES (?, ?, ?, ?, ?, ?)
		`)
		if err != nil {
			err = fmt.Errorf("failed to create prepared statement for scores table: %w", err)
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(
			score.ServiceName,
			score.BoxName,
			score.TeamNum,
			score.Status,
			score.Message,
			ts,
		)
		if err != nil {
			err = fmt.Errorf("failed to insert into scores table: %w", err)
			return err
		}
	}

	// update the entry in the current service status table
	{
		stmt, err := tx.Prepare(`
		  INSERT OR REPLACE INTO recent_scores (service, box, team_num, status, message, timestamp)
			VALUES (?, ?, ?, ?, ?, ?)
		`)
		if err != nil {
			err = fmt.Errorf("failed to create prepared statement for recent_scores table: %w", err)
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(
			score.ServiceName,
			score.BoxName,
			score.TeamNum,
			score.Status,
			score.Message,
			ts,
		)
		if err != nil {
			err = fmt.Errorf("failed to insert into recent_scores table: %w", err)
			return err
		}
	}

	// update score in final scores table
	{
		stmt, err := tx.Prepare(`
			UPDATE final_scores
				SET
					total = total + 1,
					passed = passed + ?
			WHERE service = ? AND box = ? AND team_num = ?
		`)
		if err != nil {
			err = fmt.Errorf("failed to create prepared statement for final_scores table: %w", err)
			return err
		}
		defer stmt.Close()

		var passed int = 0
		if score.Status {
			passed = 1
		}

		res, err := stmt.Exec(
			passed,
			score.ServiceName,
			score.BoxName,
			score.TeamNum,
		)
		if err != nil {
			err = fmt.Errorf("failed to insert into final_scores table: %w", err)
			return err
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			err = fmt.Errorf("failed to get rows affected: %w", err)
			return err
		}
		if rowsAffected == 0 {
			err = fmt.Errorf("no matching row found in final_scores for service=%s, box=%s, team_num=%d",
				score.ServiceName, score.BoxName, score.TeamNum)
			return err
		}
	}

	// execute
	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("failed to commit transaction: %w", err)
		return err
	}

	return nil
}
