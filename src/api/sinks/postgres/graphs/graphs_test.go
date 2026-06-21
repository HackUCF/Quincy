package graphs_test

import (
	"context"
	"os"
	"testing"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/api/db/graphs"
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

	if err := testutil.SeedScoring(ctx, pool, testCfg); err != nil {
		panic(err)
	}

	// Seed scores with distinct timestamps so GetScoresData/GetCompDuration work.
	// bucketUs = (2000000000 - 1000000000) / 100 = 10000000 µs = 10s, safe for bucketing.
	_, err = pool.Exec(ctx, `
		INSERT INTO scores (service, box, team_num, status, message, timestamp)
		VALUES ('http', 'testbox', 1, true, 'ok', 1000000000),
		       ('http', 'testbox', 1, true, 'ok', 2000000000)
	`)
	if err != nil {
		panic(err)
	}

	// Seed recent_scores for scoreboard/heatmap tests.
	_, err = pool.Exec(ctx, `
		INSERT INTO recent_scores (service, box, team_num, status, message, timestamp)
		VALUES ('http', 'testbox', 1, true,  'pass', 1000000000),
		       ('ssh',  'testbox', 1, false, 'fail', 2000000000)
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
