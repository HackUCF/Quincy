package users

import (
	"context"
	"fmt"

	"github.com/HackUCF/quincy/common/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UpdateUser performs a password change request.
// This performs no validation. Will cause an error if invalid input is given.
// This query can only update the password in an existing row.
func UpdateUser(
	ctx context.Context,
	db *pgxpool.Pool,
	userListID types.UserListName,
	teamNum types.TeamNum,
	user types.User,
) error {

	query := `
		UPDATE scoring_users
		SET password = $1
		WHERE team_num = $2
			AND user_list = $3
			AND username  = $4
	`

	res, err := db.Exec(
		ctx,
		query,
		user.Password,
		teamNum,
		userListID,
		user.Username,
	)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	rows := res.RowsAffected()
	if rows != 1 {
		return fmt.Errorf("wrong number of rows affected: %d. should be 1", rows)
	}

	return nil
}
