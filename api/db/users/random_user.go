package users

import (
	"fmt"

	"github.com/HackUCF/Quincy/api/db/conn"
	"github.com/HackUCF/Quincy/common/types"
)

// GetRandomUser pulls a random user from a specific userlist for a specific team.
// This performs no validation. Will fail with an error if invalid input is given.
func GetRandomUser(userListID types.UserListName, teamNum types.TeamNum) (types.User, error) {
	db := conn.Get()
	var u types.User

	query := `
	  SELECT username, password, domain, netbios
		FROM scoring_users
		WHERE team_num = ?
		  AND user_list = ?
		ORDER BY RANDOM()
		LIMIT 1
	`

	row := db.QueryRow(query, teamNum, userListID)

	err := row.Scan(&u.Username, &u.Password, &u.DomainName, &u.NetBIOSName)
	if err != nil {
		err = fmt.Errorf("failed to scan or query db: %w", err)
		return u, err
	}

	return u, nil
}
