package testutil

import (
	"context"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/api/db/scoring"
	"github.com/HackUCF/quincy/api/db/users"
	"github.com/HackUCF/quincy/common/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MinimalConfig returns a minimal APIConfigSpec for tests: 2 teams, 1 box, 2 services.
func MinimalConfig() *config.APIConfigSpec {
	return &config.APIConfigSpec{
		NumTeams: 2,
		Boxes: []config.BoxSpec{
			{
				Name: "testbox",
				Host: "10.0.{}.1",
				Services: []types.ServiceSpec{
					{Name: "http", CheckName: "stub"},
					{Name: "ssh", CheckName: "stub", UserList: "local"},
				},
			},
		},
		UserLists: []config.UserListSpec{
			{
				Name: "local",
				Users: []config.UserSpec{
					{Username: "alice", Password: "pass1"},
					{Username: "bob", Password: "pass2"},
				},
			},
		},
		HTTP: config.HTTPSpec{Host: "127.0.0.1", Port: 8888},
	}
}

// SetupConfig sets config.TeamRange from cfg without calling config.SetConfig.
func SetupConfig(cfg *config.APIConfigSpec) {
	config.TeamRange = make([]types.TeamNum, 0, cfg.NumTeams)
	for i := types.TeamNum(1); i <= cfg.NumTeams; i++ {
		config.TeamRange = append(config.TeamRange, i)
	}
}

// SeedUsers inserts default users into scoring_users. Requires SetupConfig to have been called.
func SeedUsers(ctx context.Context, db *pgxpool.Pool, cfg *config.APIConfigSpec) error {
	return users.InitUsers(ctx, db, cfg)
}

// SeedScoring inserts zero-count rows into final_scores. Requires SetupConfig to have been called.
func SeedScoring(ctx context.Context, db *pgxpool.Pool, cfg *config.APIConfigSpec) error {
	return scoring.InitScoring(ctx, db, cfg)
}
