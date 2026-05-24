/*
Package users contains the DB functions and all SQL queries involving scoring users.
This package exports an init function (InitUsers) that is required for proper use.
It requires the DB conneciton to have already been initialized.
*/
package users

import (
	"context"
	"fmt"

	"github.com/HackUCF/quincy/api/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InitUsers fills the users table in the database with default usernames and passwords.
// Users that have had passwords changed in a previous run of the API will be ignored.
// This keeps consistent state between restarts.
// This requires the config to be loaded.
func InitUsers(
	ctx context.Context,
	db *pgxpool.Pool,
	cfg *config.APIConfigSpec,
) error {

	query := `
		INSERT INTO scoring_users (team_num, user_list, username, password, domain, netbios)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT DO NOTHING
	`

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("could not begin database transaction: %w", err)
	}

	for _, t := range config.TeamRange {
		for _, ul := range cfg.UserLists {
			for _, u := range ul.Users {
				// for every team, userlist, and user
				_, err := tx.Exec(
					ctx,
					query,
					t,
					ul.Name,
					u.Username,
					u.Password,
					ul.DomainName,
					ul.NetBIOSName,
				)
				if err != nil {
					return fmt.Errorf("could not insert into table: %w", err)
				}
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("could not commit database transaction: %w", err)
	}

	return nil
}
