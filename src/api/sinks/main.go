package sinks

import (
	"context"

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
	return types.User{}, nil
}
