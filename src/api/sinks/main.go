package sinks

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/api/sinks/opentelemetry"
	"github.com/HackUCF/quincy/api/sinks/postgres"
	"github.com/HackUCF/quincy/api/sinks/postgres/agent"
	"github.com/HackUCF/quincy/common/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

// InitSinks calls the correct initialization functions for each of the enabled sinks.
func InitSinks(cfg *config.APIConfigSpec) error {

	if cfg.Sinks.DBEnabled() {
		if err := postgres.InitDB(cfg); err != nil {
			return fmt.Errorf("postgres sink: %w", err)
		}
	}

	if cfg.Sinks.OTelEnabled() {
		if _, err := opentelemetry.Init(context.Background(), cfg.Sinks.OTelConfig); err != nil {
			return fmt.Errorf("opentelemetry sink: %w", err)
		}
	}

	return nil
}

// AddScore sends the score to every configured sink.
func AddScore(
	ctx context.Context,
	sinks config.Sinks,
	db *pgxpool.Pool,
	score types.Score,
) error {

	if sinks.DBEnabled() {
		if err := agent.AddScore(ctx, db, score); err != nil {
			return err
		}
	}

	if sinks.OTelEnabled() {
		if err := opentelemetry.AddScore(ctx, score); err != nil {
			return err
		}
	}

	return nil
}

// GetRandomUser returns a random user+password.
// If the database sink is configured it pulls from there,
// otherwise it just pulls the default creds from the config.
func GetRandomUser(
	ctx context.Context,
	cfg *config.APIConfigSpec,
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

		// find the requested one
		if ul.Name == userListID {

			if len(ul.Users) == 0 {

				// fail if it's empty
				return types.User{}, fmt.Errorf("user list %q is empty", userListID)
			}

			// otherwise grab a random one
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
