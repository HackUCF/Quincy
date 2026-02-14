package users

import (
	"fmt"

	"github.com/HackUCF/Quincy/api/db/conn"
	"github.com/HackUCF/Quincy/common/types"
)

type AllUsers map[types.TeamNum]map[types.UserListID][]types.User

func GetAllUsers() (AllUsers, error) {
	allUsers := make(AllUsers)

	query := `
	  SELECT team_num, user_list, username, password, domain, netbios
		FROM scoring_users
	`
	rows, err := conn.Get().Query(query)
	if err != nil {
		return allUsers, fmt.Errorf("failed to list users in db: %w", err)
	}

	// loop vars
	var t types.TeamNum
	var ul types.UserListID
	var user types.User

	for rows.Next() {
		err := rows.Scan(&t, &ul, &user.Username, &user.Password, &user.DomainName, &user.NetBIOSName)
		if err != nil {
			return allUsers, fmt.Errorf("error scanning row: %w", err)
		}

		if _, ok := allUsers[t]; !ok {
			allUsers[t] = make(map[types.UserListID][]types.User)
		}

		if _, ok := allUsers[t][ul]; !ok {
			allUsers[t][ul] = make([]types.User, 0)
		}

		allUsers[t][ul] = append(allUsers[t][ul], user)
	}

	return allUsers, nil
}
