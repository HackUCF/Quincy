package scoring_test

import (
	"context"
	"os"
	"testing"

	dbagent "github.com/HackUCF/quincy/api/db/agent"
	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/api/db/scoring"
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

	if err := testutil.SeedScoring(ctx, pool, testCfg); err != nil {
		panic(err)
	}

	// seed one passing score for team 1, box testbox, service http
	_ = dbagent.AddScore(ctx, pool, types.Score{
		ServiceName: "http",
		BoxName:     "testbox",
		TeamNum:     1,
		Status:      true,
		Message:     "ok",
	})

	os.Exit(m.Run())
}

func TestGetTeamScores_returnsAllTeams(t *testing.T) {
	scores, err := scoring.GetTeamScores(context.Background(), testPool)
	if err != nil {
		t.Fatalf("GetTeamScores: %v", err)
	}
	if len(scores) == 0 {
		t.Error("expected scores for at least one team")
	}
}

func TestGetBoxScores_returnsResult(t *testing.T) {
	scores, err := scoring.GetBoxScores(context.Background(), testPool)
	if err != nil {
		t.Fatalf("GetBoxScores: %v", err)
	}
	if len(scores) == 0 {
		t.Error("expected box scores")
	}
}

func TestGetServiceScores_returnsResult(t *testing.T) {
	scores, err := scoring.GetServiceScores(context.Background(), testPool)
	if err != nil {
		t.Fatalf("GetServiceScores: %v", err)
	}
	if len(scores) == 0 {
		t.Error("expected service scores")
	}
}

func TestGetCurrentServiceStatus_returnsRows(t *testing.T) {
	rows, err := scoring.GetCurrentServiceStatus(context.Background(), testPool, testCfg)
	if err != nil {
		t.Fatalf("GetCurrentServiceStatus: %v", err)
	}
	// at least one row from seeded score
	if len(rows) == 0 {
		t.Error("expected at least one row in current status")
	}
}

func TestGetDetailedScores_returnsAllTeams(t *testing.T) {
	scores, err := scoring.GetDetailedScores(context.Background(), testPool)
	if err != nil {
		t.Fatalf("GetDetailedScores: %v", err)
	}
	if len(scores) == 0 {
		t.Error("expected detailed scores for at least one team")
	}
}

func TestGetTeamScores_team1HasPassedChecks(t *testing.T) {
	scores, err := scoring.GetTeamScores(context.Background(), testPool)
	if err != nil {
		t.Fatalf("GetTeamScores: %v", err)
	}
	r, ok := scores[types.TeamNum(1)]
	if !ok {
		t.Fatal("no score for team 1")
	}
	if r.ChecksPassed == 0 {
		t.Error("team 1 should have at least one passed check from seeded score")
	}
}
