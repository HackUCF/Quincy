/*
Package users contains the DB functions and all SQL queries involving scoring users.
This package exports an init function (InitUsers) that is required for proper use.
It requires the DB conneciton to have already been initialized.
*/
package users

import (
	"fmt"

	"github.com/HackUCF/Quincy/api/config"
	"github.com/HackUCF/Quincy/api/db/conn"
)

// InitUsers fills the users table in the database with default usernames and passwords.
// Users that have had passwords changed in a previous run of the API will be ignored.
// This keeps consistent state between restarts.
// This requires the config to be loaded.
func InitUsers(cfg *config.APIConfigSpec) error {
	db := conn.Get()

	query := `
	  INSERT OR IGNORE INTO scoring_users (team_num, user_list, username, password, domain, netbios)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin database transaction: %w", err)
	}

	for _, t := range config.TeamRange {
		for _, ul := range cfg.UserLists {
			for _, u := range ul.Users {
				// for every team, userlist, and user
				_, err := tx.Exec(
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

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("could not commit database transaction: %w", err)
	}

	return nil
}
