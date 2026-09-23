package agent_test

import (
	"context"
	"os"
	"testing"

	"github.com/HackUCF/Quincy/src/api/config"
	dbagent "github.com/HackUCF/Quincy/src/api/sinks/postgres/agent"
	"github.com/HackUCF/Quincy/src/api/sinks/postgres/pauses"
	"github.com/HackUCF/Quincy/src/common/types"
	"github.com/HackUCF/Quincy/src/testutil"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testPool *pgxpool.Pool
var testCfg *config.APIConfigSpec

func TestMain(m *testing.M) {
	ctx := context.Background()

	pool, cleanup, err := testutil.NewTestDB(ctx)
	if err != nil {
		testutil.SkipDBTests(err)
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
	if err := testutil.SeedPauses(ctx, pool, testCfg); err != nil {
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
		Stdout:      "ok",
		Stderr:      "warn",
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
		Stdout:      "",
		Stderr:      "timed out",
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

func TestAddScore_persistsStdoutAndStderr(t *testing.T) {
	ctx := context.Background()
	score := types.Score{
		ServiceName: "http",
		BoxName:     "testbox",
		TeamNum:     2,
		Status:      true,
		Stdout:      "HTTP 200 OK",
		Stderr:      "cert nearly expired",
	}
	if err := dbagent.AddScore(ctx, testPool, score); err != nil {
		t.Fatalf("AddScore: %v", err)
	}

	var stdout, stderr string
	if err := testPool.QueryRow(ctx,
		`SELECT stdout, stderr FROM scores
		 WHERE service = $1 AND box = $2 AND team_num = $3
		 ORDER BY id DESC LIMIT 1`,
		"http", "testbox", 2,
	).Scan(&stdout, &stderr); err != nil {
		t.Fatalf("query scores: %v", err)
	}
	if stdout != "HTTP 200 OK" {
		t.Errorf("scores.stdout = %q, want %q", stdout, "HTTP 200 OK")
	}
	if stderr != "cert nearly expired" {
		t.Errorf("scores.stderr = %q, want %q", stderr, "cert nearly expired")
	}

	if err := testPool.QueryRow(ctx,
		`SELECT stdout, stderr FROM recent_scores
		 WHERE service = $1 AND box = $2 AND team_num = $3`,
		"http", "testbox", 2,
	).Scan(&stdout, &stderr); err != nil {
		t.Fatalf("query recent_scores: %v", err)
	}
	if stdout != "HTTP 200 OK" {
		t.Errorf("recent_scores.stdout = %q, want %q", stdout, "HTTP 200 OK")
	}
	if stderr != "cert nearly expired" {
		t.Errorf("recent_scores.stderr = %q, want %q", stderr, "cert nearly expired")
	}
}

func TestAddScore_upsertOverwritesStdoutAndStderr(t *testing.T) {
	ctx := context.Background()
	base := types.Score{ServiceName: "ssh", BoxName: "testbox", TeamNum: 1}

	first := base
	first.Status = true
	first.Stdout = "first stdout"
	first.Stderr = "first stderr"
	if err := dbagent.AddScore(ctx, testPool, first); err != nil {
		t.Fatalf("AddScore (first): %v", err)
	}

	second := base
	second.Status = false
	second.Stdout = "second stdout"
	second.Stderr = "second stderr"
	if err := dbagent.AddScore(ctx, testPool, second); err != nil {
		t.Fatalf("AddScore (second): %v", err)
	}

	var stdout, stderr string
	var status bool
	if err := testPool.QueryRow(ctx,
		`SELECT stdout, stderr, status FROM recent_scores
		 WHERE service = $1 AND box = $2 AND team_num = $3`,
		"ssh", "testbox", 1,
	).Scan(&stdout, &stderr, &status); err != nil {
		t.Fatalf("query recent_scores: %v", err)
	}
	if stdout != "second stdout" {
		t.Errorf("stdout = %q, want %q", stdout, "second stdout")
	}
	if stderr != "second stderr" {
		t.Errorf("stderr = %q, want %q", stderr, "second stderr")
	}
	if status {
		t.Error("status should be false after upsert with a failing score")
	}

	// both checks must still be recorded in the append-only scores table
	var count int
	if err := testPool.QueryRow(ctx,
		`SELECT COUNT(*) FROM scores WHERE service = $1 AND box = $2 AND team_num = $3`,
		"ssh", "testbox", 1,
	).Scan(&count); err != nil {
		t.Fatalf("query scores: %v", err)
	}
	if count < 2 {
		t.Errorf("scores rows = %d, want at least 2", count)
	}
}

func TestAddScore_emptyStdoutAndStderr(t *testing.T) {
	ctx := context.Background()
	score := types.Score{
		ServiceName: "ssh",
		BoxName:     "testbox",
		TeamNum:     2,
		Status:      true,
	}
	if err := dbagent.AddScore(ctx, testPool, score); err != nil {
		t.Fatalf("AddScore: %v", err)
	}

	var stdout, stderr string
	if err := testPool.QueryRow(ctx,
		`SELECT stdout, stderr FROM recent_scores
		 WHERE service = $1 AND box = $2 AND team_num = $3`,
		"ssh", "testbox", 2,
	).Scan(&stdout, &stderr); err != nil {
		t.Fatalf("query recent_scores: %v", err)
	}
	if stdout != "" || stderr != "" {
		t.Errorf("stdout/stderr = %q/%q, want empty/empty", stdout, stderr)
	}
}

// pauseFor pauses the competition for the duration of the test and restores the
// running state afterwards, so a pause cannot leak into a later test.
func pauseFor(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	if _, err := pauses.ChangePauseState(ctx, testPool, types.Paused); err != nil {
		t.Fatalf("pause: %v", err)
	}
	t.Cleanup(func() {
		if _, err := pauses.ChangePauseState(ctx, testPool, types.Unpaused); err != nil {
			t.Fatalf("unpause: %v", err)
		}
	})
}

func TestAddScore_pausedArchivesButDoesNotCount(t *testing.T) {
	ctx := context.Background()

	testPool.Exec(ctx,
		`UPDATE final_scores SET passed = 0, total = 0 WHERE service = $1 AND box = $2 AND team_num = $3`,
		"http", "testbox", 2,
	)

	pauseFor(t)

	score := types.Score{
		ServiceName: "http",
		BoxName:     "testbox",
		TeamNum:     2,
		Status:      true,
		Stdout:      "paused run",
		Stderr:      "",
	}
	if err := dbagent.AddScore(ctx, testPool, score); err != nil {
		t.Fatalf("AddScore: %v", err)
	}

	// the result is still history and still the current status...
	var count int
	if err := testPool.QueryRow(ctx,
		`SELECT COUNT(*) FROM scores WHERE service = $1 AND box = $2 AND team_num = $3 AND stdout = $4`,
		"http", "testbox", 2, "paused run",
	).Scan(&count); err != nil {
		t.Fatalf("query scores: %v", err)
	}
	if count != 1 {
		t.Errorf("scores holds %d rows for the paused check, want 1 — a pause must not discard results", count)
	}

	var stdout string
	if err := testPool.QueryRow(ctx,
		`SELECT stdout FROM recent_scores WHERE service = $1 AND box = $2 AND team_num = $3`,
		"http", "testbox", 2,
	).Scan(&stdout); err != nil {
		t.Fatalf("query recent_scores: %v", err)
	}
	if stdout != "paused run" {
		t.Errorf("recent_scores.stdout = %q, want %q", stdout, "paused run")
	}

	// ...but it costs the team nothing either way
	var passed, total int
	if err := testPool.QueryRow(ctx,
		`SELECT passed, total FROM final_scores WHERE service = $1 AND box = $2 AND team_num = $3`,
		"http", "testbox", 2,
	).Scan(&passed, &total); err != nil {
		t.Fatalf("query final_scores: %v", err)
	}
	if total != 0 || passed != 0 {
		t.Errorf("final_scores = passed %d / total %d after a paused check, want 0 / 0", passed, total)
	}
}

func TestAddScore_unpausingResumesCounting(t *testing.T) {
	ctx := context.Background()

	testPool.Exec(ctx,
		`UPDATE final_scores SET passed = 0, total = 0 WHERE service = $1 AND box = $2 AND team_num = $3`,
		"ssh", "testbox", 1,
	)

	score := types.Score{
		ServiceName: "ssh",
		BoxName:     "testbox",
		TeamNum:     1,
		Status:      true,
		Stdout:      "ok",
		Stderr:      "",
	}

	// the same check submitted once on each side of an unpause
	if _, err := pauses.ChangePauseState(ctx, testPool, types.Paused); err != nil {
		t.Fatalf("pause: %v", err)
	}
	if err := dbagent.AddScore(ctx, testPool, score); err != nil {
		t.Fatalf("AddScore (paused): %v", err)
	}
	if _, err := pauses.ChangePauseState(ctx, testPool, types.Unpaused); err != nil {
		t.Fatalf("unpause: %v", err)
	}

	if err := dbagent.AddScore(ctx, testPool, score); err != nil {
		t.Fatalf("AddScore (running): %v", err)
	}

	var passed, total int
	if err := testPool.QueryRow(ctx,
		`SELECT passed, total FROM final_scores WHERE service = $1 AND box = $2 AND team_num = $3`,
		"ssh", "testbox", 1,
	).Scan(&passed, &total); err != nil {
		t.Fatalf("query final_scores: %v", err)
	}
	if passed != 1 || total != 1 {
		t.Errorf("final_scores = passed %d / total %d, want 1 / 1 — only the unpaused check counts", passed, total)
	}
}
