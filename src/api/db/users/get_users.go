package users

import (
	"context"
	"fmt"

	"github.com/HackUCF/quincy/common/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AllUsers map[types.TeamNum]map[types.UserListName][]types.User

func GetAllUsers(ctx context.Context, db *pgxpool.Pool) (AllUsers, error) {

	allUsers := make(AllUsers)

	query := `
		SELECT team_num, user_list, username, password, domain, netbios
		FROM scoring_users
	`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return allUsers, fmt.Errorf("failed to list users in db: %w", err)
	}

	// loop vars
	var t types.TeamNum
	var ul types.UserListName
	var user types.User

	for rows.Next() {
		err := rows.Scan(&t, &ul, &user.Username, &user.Password, &user.DomainName, &user.NetBIOSName)
		if err != nil {
			return allUsers, fmt.Errorf("error scanning row: %w", err)
		}

		if _, ok := allUsers[t]; !ok {
			allUsers[t] = make(map[types.UserListName][]types.User)
		}

		if _, ok := allUsers[t][ul]; !ok {
			allUsers[t][ul] = make([]types.User, 0)
		}

		allUsers[t][ul] = append(allUsers[t][ul], user)
	}

	return allUsers, nil
}
