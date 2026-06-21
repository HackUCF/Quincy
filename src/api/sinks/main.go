package sinks

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/api/sinks/postgres"
	"github.com/HackUCF/quincy/api/sinks/postgres/agent"
	"github.com/HackUCF/quincy/common/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InitSinks calls the correct initialization functions for each of the enabled sinks.
func InitSinks(cfg *config.APIConfigSpec) error {

	// if postgres is configured
	if cfg.Sinks.PGConfig != (config.PGConfig{}) {
		postgres.InitDB(cfg)
	}

	return nil
}

// Add score calls the right
func AddScore(
	ctx context.Context,
	sinks config.Sinks,
	db *pgxpool.Pool,
	score types.Score,
) error {

	if sinks.DBEnabled() {
		agent.AddScore(ctx, db, score)
	}

	return nil
}

// GetRandomUser returns a random user+password.
// If the database sink is configured it pulls from there,
// otherwise it just pulls the default creds from the config.
func GetRandomUser(
	ctx context.Context,
	cfg config.APIConfigSpec,
	db *pgxpool.Pool,
	userListID types.UserListName,
	teamNum types.TeamNum,
) (types.User, error) {

	// return a user from the db
	// this supports password changes
	if cfg.Sinks.DBEnabled() {
		return agent.GetRandomUser(ctx, db, userListID, teamNum)
	}

	// otherwise pull random user from the config file

	// loop through the user lists
	for _, ul := range cfg.UserLists {

		// find the requests one
		if ul.Name == userListID {

			if len(ul.Users) == 0 {

				// fail if it's empty
				return types.User{}, fmt.Errorf("user list %q is empty", userListID)
			}

			// otherrwise grab a random one
			u := ul.Users[rand.IntN(len(ul.Users))]
			return types.User{
				Username:    u.Username,
				Password:    u.Password,
				DomainName:  ul.DomainName,
				NetBIOSName: ul.NetBIOSName,
			}, nil
		}
	}

	// error if userlist doesn't exist
	return types.User{}, fmt.Errorf("user list %q not found", userListID)
}
