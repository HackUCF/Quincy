package users

import (
	"database/sql"
	"fmt"

	"github.com/HackUCF/quincy/common/types"
)

// UpdateUser performs a password change request.
// This performs no validation. Will cause an error if invalid input is given.
// This query can only update the password in an existing row.
func UpdateUser(db *sql.DB, userListID types.UserListName, teamNum types.TeamNum, user types.User) error {

	query := `
	  UPDATE scoring_users
		SET password = ?
		WHERE team_num = ?
		  AND user_list = ?
			AND username = ?
  `

	res, err := db.Exec(query, user.Password, teamNum, userListID, user.Username)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows: %w", err)
	}
	if rows != 1 {
		return fmt.Errorf("wrong number of rows affected: %d. should be 1", rows)
	}

	return nil
}
