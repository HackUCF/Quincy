package users

import (
	"context"
	"fmt"

	"github.com/HackUCF/quincy/common/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetRandomUser pulls a random user from a specific userlist for a specific team.
// This performs no validation. Will fail with an error if invalid input is given.
func GetRandomUser(
	ctx context.Context,
	db *pgxpool.Pool,
	userListID types.UserListName,
	teamNum types.TeamNum,
) (types.User, error) {

	var u types.User

	query := `
		SELECT username, password, domain, netbios
		FROM scoring_users
		WHERE team_num = $1
			AND user_list = $2
		ORDER BY RANDOM()
		LIMIT 1
	`

	row := db.QueryRow(ctx, query, teamNum, userListID)

	err := row.Scan(&u.Username, &u.Password, &u.DomainName, &u.NetBIOSName)
	if err != nil {
		err = fmt.Errorf("failed to scan or query db: %w", err)
		return u, err
	}

	return u, nil
}
