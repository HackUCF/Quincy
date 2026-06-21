package users_test

import (
	"context"
	"os"
	"testing"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/api/db/users"
	"github.com/HackUCF/quincy/common/types"
	"github.com/HackUCF/quincy/testutil"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testPool *pgxpool.Pool
var testCfg *config.APIConfigSpec

func TestMain(m *testing.M) {
	ctx := context.Background()

	pool, cleanup, err := testutil.NewTestDB(ctx)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	testPool = pool
	testCfg = testutil.MinimalConfig()
	testutil.SetupConfig(testCfg)

	if err := testutil.SeedUsers(ctx, pool, testCfg); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestInitUsers_idempotent(t *testing.T) {
	ctx := context.Background()
	// calling InitUsers again should not error (ON CONFLICT DO NOTHING)
	if err := users.InitUsers(ctx, testPool, testCfg); err != nil {
		t.Fatalf("InitUsers (second call): %v", err)
	}
}

func TestGetAllUsers_returnsSeededUsers(t *testing.T) {
	ctx := context.Background()
	all, err := users.GetAllUsers(ctx, testPool)
	if err != nil {
		t.Fatalf("GetAllUsers: %v", err)
	}

	// MinimalConfig has 2 teams × 1 userlist × 2 users = 4 rows
	total := 0
	for _, lists := range all {
		for _, us := range lists {
			total += len(us)
		}
	}
	if total == 0 {
		t.Error("expected seeded users, got none")
	}
}

func TestGetAllUsers_containsExpectedUsername(t *testing.T) {
	ctx := context.Background()
	all, err := users.GetAllUsers(ctx, testPool)
	if err != nil {
		t.Fatalf("GetAllUsers: %v", err)
	}

	found := false
	for _, lists := range all {
		for _, us := range lists {
			for _, u := range us {
				if u.Username == "alice" {
					found = true
				}
			}
		}
	}
	if !found {
		t.Error("expected user 'alice' in results")
	}
}

func TestUpdateUser_changesPassword(t *testing.T) {
	ctx := context.Background()
	updated := types.User{Username: "alice", Password: "newpass"}
	if err := users.UpdateUser(ctx, testPool, "local", 1, updated); err != nil {
		t.Fatalf("UpdateUser: %v", err)
	}

	var password string
	err := testPool.QueryRow(ctx,
		`SELECT password FROM scoring_users WHERE team_num = $1 AND user_list = $2 AND username = $3`,
		1, "local", "alice",
	).Scan(&password)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if password != "newpass" {
		t.Errorf("password = %q, want newpass", password)
	}
}
