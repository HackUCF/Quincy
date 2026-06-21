package agent_test

import (
	"context"
	"os"
	"testing"

	"github.com/HackUCF/quincy/api/config"
	dbagent "github.com/HackUCF/quincy/api/sinks/postgres/agent"
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
	if err := testutil.SeedScoring(ctx, pool, testCfg); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestAddScore_insertsToAllTables(t *testing.T) {
	ctx := context.Background()
	score := types.Score{
		ServiceName: "http",
		BoxName:     "testbox",
		TeamNum:     1,
		Status:      true,
		Message:     "ok",
	}
	if err := dbagent.AddScore(ctx, testPool, score); err != nil {
		t.Fatalf("AddScore: %v", err)
	}

	var count int
	if err := testPool.QueryRow(ctx,
		`SELECT COUNT(*) FROM scores WHERE service = $1 AND box = $2 AND team_num = $3`,
		"http", "testbox", 1,
	).Scan(&count); err != nil {
		t.Fatalf("query scores: %v", err)
	}
	if count == 0 {
		t.Error("no row in scores table after AddScore")
	}

	var status bool
	if err := testPool.QueryRow(ctx,
		`SELECT status FROM recent_scores WHERE service = $1 AND box = $2 AND team_num = $3`,
		"http", "testbox", 1,
	).Scan(&status); err != nil {
		t.Fatalf("query recent_scores: %v", err)
	}
	if !status {
		t.Error("recent_scores.status should be true")
	}

	var passed, total int
	if err := testPool.QueryRow(ctx,
		`SELECT passed, total FROM final_scores WHERE service = $1 AND box = $2 AND team_num = $3`,
		"http", "testbox", 1,
	).Scan(&passed, &total); err != nil {
		t.Fatalf("query final_scores: %v", err)
	}
	if total == 0 {
		t.Error("final_scores.total should be > 0 after AddScore")
	}
	if passed == 0 {
		t.Error("final_scores.passed should be > 0 for a passing score")
	}
}

func TestAddScore_failDoesNotIncrementPassed(t *testing.T) {
	ctx := context.Background()

	// reset counters for this service first
	testPool.Exec(ctx,
		`UPDATE final_scores SET passed = 0, total = 0 WHERE service = $1 AND box = $2 AND team_num = $3`,
		"ssh", "testbox", 2,
	)

	score := types.Score{
		ServiceName: "ssh",
		BoxName:     "testbox",
		TeamNum:     2,
		Status:      false,
		Message:     "timed out",
	}
	if err := dbagent.AddScore(ctx, testPool, score); err != nil {
		t.Fatalf("AddScore: %v", err)
	}

	var passed, total int
	if err := testPool.QueryRow(ctx,
		`SELECT passed, total FROM final_scores WHERE service = $1 AND box = $2 AND team_num = $3`,
		"ssh", "testbox", 2,
	).Scan(&passed, &total); err != nil {
		t.Fatalf("query final_scores: %v", err)
	}
	if total != 1 {
		t.Errorf("total = %d, want 1", total)
	}
	if passed != 0 {
		t.Errorf("passed = %d, want 0 for a failing score", passed)
	}
}

func TestGetRandomUser_returnsUser(t *testing.T) {
	ctx := context.Background()
	u, err := dbagent.GetRandomUser(ctx, testPool, "local", 1)
	if err != nil {
		t.Fatalf("GetRandomUser: %v", err)
	}
	if u.Username == "" {
		t.Error("returned empty username")
	}
}

func TestGetRandomUser_unknownList(t *testing.T) {
	ctx := context.Background()
	_, err := dbagent.GetRandomUser(ctx, testPool, "nonexistent", 1)
	if err == nil {
		t.Fatal("expected error for unknown user list, got nil")
	}
}
