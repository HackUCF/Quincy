package graphs_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/HackUCF/Quincy/src/api/config"
	"github.com/HackUCF/Quincy/src/api/sinks/postgres/graphs"
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

	if err := testutil.SeedScoring(ctx, pool, testCfg); err != nil {
		panic(err)
	}
	if err := testutil.SeedPauses(ctx, pool, testCfg); err != nil {
		panic(err)
	}

	// Seed scores with distinct timestamps so GetScoresData/GetCompDuration work.
	// bucketUs = (2000000000 - 1000000000) / 100 = 10000000 µs = 10s, safe for bucketing.
	_, err = pool.Exec(ctx, `
		INSERT INTO scores (service, box, team_num, status, stdout, stderr, timestamp)
		VALUES ('http', 'testbox', 1, true, 'ok', '', 1000000000),
		       ('http', 'testbox', 1, true, 'ok', '', 2000000000)
	`)
	if err != nil {
		panic(err)
	}

	// Seed recent_scores for scoreboard/heatmap tests.
	_, err = pool.Exec(ctx, `
		INSERT INTO recent_scores (service, box, team_num, status, stdout, stderr, timestamp)
		VALUES ('http', 'testbox', 1, true,  'pass', '',        1000000000),
		       ('ssh',  'testbox', 1, false, '',     'failerr', 2000000000)
	`)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestGetScoreboardData_returnsData(t *testing.T) {
	data, err := graphs.GetScoreboardData(context.Background(), testPool)
	if err != nil {
		t.Fatalf("GetScoreboardData: %v", err)
	}
	if string(data.XLabels) == "[]" || string(data.XLabels) == "null" {
		t.Error("expected non-empty XLabels")
	}
}

func TestGetHeatmapData_returnsData(t *testing.T) {
	data, err := graphs.GetHeatmapData(context.Background(), testPool)
	if err != nil {
		t.Fatalf("GetHeatmapData: %v", err)
	}
	if data == nil {
		t.Error("expected non-nil HeatmapData")
	}
}

func TestGetStandingsData_returnsTeams(t *testing.T) {
	data, err := graphs.GetStandingsData(context.Background(), testPool)
	if err != nil {
		t.Fatalf("GetStandingsData: %v", err)
	}
	if string(data.Labels) == "[]" || string(data.Labels) == "null" {
		t.Error("expected non-empty Labels")
	}
}

func TestGetScoresData_returnsData(t *testing.T) {
	data, err := graphs.GetScoresData(context.Background(), testPool)
	if err != nil {
		t.Fatalf("GetScoresData: %v", err)
	}
	if data == nil {
		t.Error("expected non-nil ScoresData")
	}
}

// scoreboard datapoints carry stdout and stderr joined into a single "msg" field.
func TestGetScoreboardData_messageJoinsStdoutAndStderr(t *testing.T) {
	data, err := graphs.GetScoreboardData(context.Background(), testPool)
	if err != nil {
		t.Fatalf("GetScoreboardData: %v", err)
	}

	var points []struct {
		X   string `json:"x"`
		Msg string `json:"msg"`
	}
	if err := json.Unmarshal([]byte(data.DataPoints), &points); err != nil {
		t.Fatalf("unmarshal DataPoints: %v", err)
	}

	want := map[string]string{
		"testbox-http": "pass" + "\n\n\n" + "",
		"testbox-ssh":  "" + "\n\n\n" + "failerr",
	}
	seen := 0
	for _, p := range points {
		w, ok := want[p.X]
		if !ok {
			continue
		}
		seen++
		if p.Msg != w {
			t.Errorf("point %s msg = %q, want %q", p.X, p.Msg, w)
		}
	}
	if seen != len(want) {
		t.Errorf("matched %d datapoints, want %d (got %+v)", seen, len(want), points)
	}
}
