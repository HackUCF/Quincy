package scoring_test

import (
	"context"
	"testing"

	"github.com/HackUCF/Quincy/src/api/config"
	"github.com/HackUCF/Quincy/src/api/sinks/postgres/scoring"
	"github.com/HackUCF/Quincy/src/common/types"
)

// recentTestTeam is kept separate from the scores seeded in TestMain so that
// the aggregate tests in this package are not disturbed by the rows below.
const recentTestTeam = types.TeamNum(2)

// seedTimedScores inserts rows straight into the scores table with explicit
// timestamps. AddScore stamps its own time, which is too coarse to assert a
// deterministic order on, so these tests write the column directly.
func seedTimedScores(t *testing.T, service types.ServiceName, count int, baseTS int64) {
	t.Helper()
	ctx := context.Background()

	for i := range count {
		_, err := testPool.Exec(ctx, `
			INSERT INTO scores (service, box, team_num, status, stdout, stderr, timestamp)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, service, "testbox", recentTestTeam, i%2 == 0, "out", "err", baseTS+int64(i))
		if err != nil {
			t.Fatalf("seed score %d for %s: %v", i, service, err)
		}
	}
}

func TestGetRecentNScores_limitsPerService(t *testing.T) {
	seedTimedScores(t, "http", 5, 1_700_000_000_000_000)
	seedTimedScores(t, "ssh", 5, 1_700_000_000_000_000)

	const want = 2
	scores, err := scoring.GetRecentNScores(context.Background(), testPool, testCfg, recentTestTeam, want)
	if err != nil {
		t.Fatalf("GetRecentNScores: %v", err)
	}

	// both services must get their own allowance, not share one budget
	perService := make(map[types.ServiceName]int)
	for _, s := range scores {
		perService[s.ServiceName]++
	}
	for _, service := range []types.ServiceName{"http", "ssh"} {
		if perService[service] != want {
			t.Errorf("service %q returned %d rows, want %d", service, perService[service], want)
		}
	}
}

func TestGetRecentNScores_newestFirstPerService(t *testing.T) {
	const base = 1_700_100_000_000_000
	seedTimedScores(t, "http", 4, base)

	scores, err := scoring.GetRecentNScores(context.Background(), testPool, testCfg, recentTestTeam, 3)
	if err != nil {
		t.Fatalf("GetRecentNScores: %v", err)
	}

	var seen []int64
	for _, s := range scores {
		if s.ServiceName == "http" {
			seen = append(seen, s.Timestamp)
		}
	}
	if len(seen) == 0 {
		t.Fatal("no http rows returned")
	}
	for i := 1; i < len(seen); i++ {
		if seen[i] > seen[i-1] {
			t.Errorf("timestamps not descending: %d came after %d", seen[i], seen[i-1])
		}
	}
	// the newest seeded row must be first
	if seen[0] != base+3 {
		t.Errorf("first timestamp = %d, want newest seeded %d", seen[0], base+3)
	}
}

func TestGetRecentNScores_onlyRequestedTeam(t *testing.T) {
	seedTimedScores(t, "http", 2, 1_700_200_000_000_000)

	scores, err := scoring.GetRecentNScores(context.Background(), testPool, testCfg, recentTestTeam, 10)
	if err != nil {
		t.Fatalf("GetRecentNScores: %v", err)
	}
	if len(scores) == 0 {
		t.Fatal("expected rows for the seeded team")
	}
	for _, s := range scores {
		if s.TeamNum != recentTestTeam {
			t.Errorf("got row for team %d, want only team %d", s.TeamNum, recentTestTeam)
		}
	}
}

func TestGetRecentNScores_sortedByBoxThenService(t *testing.T) {
	seedTimedScores(t, "http", 2, 1_700_300_000_000_000)
	seedTimedScores(t, "ssh", 2, 1_700_300_000_000_000)

	scores, err := scoring.GetRecentNScores(context.Background(), testPool, testCfg, recentTestTeam, 2)
	if err != nil {
		t.Fatalf("GetRecentNScores: %v", err)
	}

	// config order is testbox/http then testbox/ssh, and every row of a
	// service must be contiguous
	var order []types.ServiceName
	for _, s := range scores {
		if len(order) == 0 || order[len(order)-1] != s.ServiceName {
			order = append(order, s.ServiceName)
		}
	}
	want := []types.ServiceName{"http", "ssh"}
	if len(order) != len(want) {
		t.Fatalf("service groups = %v, want %v", order, want)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Errorf("service group %d = %q, want %q", i, order[i], want[i])
		}
	}
}

func TestGetRecentNScores_returnsFewerWhenHistoryIsShort(t *testing.T) {
	// ask for far more than exists; the function must not pad or fail
	scores, err := scoring.GetRecentNScores(context.Background(), testPool, testCfg, recentTestTeam, 1000)
	if err != nil {
		t.Fatalf("GetRecentNScores: %v", err)
	}

	var count int
	err = testPool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM scores WHERE team_num = $1`, recentTestTeam).Scan(&count)
	if err != nil {
		t.Fatalf("count scores: %v", err)
	}
	if len(scores) != count {
		t.Errorf("got %d rows, want all %d rows for the team", len(scores), count)
	}
}

func TestGetRecentNScores_emptyForUnscoredServices(t *testing.T) {
	// a config whose box and service have never been checked
	cfg := &config.APIConfigSpec{
		NumTeams: 2,
		Boxes: []config.BoxSpec{
			{
				Name:     "ghostbox",
				Host:     "10.0.{}.9",
				Services: []types.ServiceSpec{{Name: "nothing", CheckName: "stub"}},
			},
		},
	}

	scores, err := scoring.GetRecentNScores(context.Background(), testPool, cfg, recentTestTeam, 5)
	if err != nil {
		t.Fatalf("GetRecentNScores: %v", err)
	}
	if len(scores) != 0 {
		t.Errorf("got %d rows for a service with no checks, want 0", len(scores))
	}
}
