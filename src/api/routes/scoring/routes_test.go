package scoring_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/HackUCF/Quincy/src/api/config"
	dbagent "github.com/HackUCF/Quincy/src/api/sinks/postgres/agent"
	"github.com/HackUCF/Quincy/src/common/types"
	"github.com/HackUCF/Quincy/src/testutil"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testPool *pgxpool.Pool
var testCfg *config.APIConfigSpec
var testRouter *gin.Engine

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

	// seed one score so current status has a row
	dbagent.AddScore(ctx, pool, types.Score{
		ServiceName: "http",
		BoxName:     "testbox",
		TeamNum:     1,
		Status:      true,
		Stdout:      "ok",
	})

	testRouter = testutil.NewTestRouter(pool, testCfg)
	os.Exit(m.Run())
}

func get(router *gin.Engine, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, path, nil)
	router.ServeHTTP(w, req)
	return w
}

func TestGetTeamScores_OK(t *testing.T) {
	w := get(testRouter, "/api/v1/scores/team")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	var result map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
}

func TestGetBoxScores_OK(t *testing.T) {
	w := get(testRouter, "/api/v1/scores/box")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
}

func TestGetServiceScores_OK(t *testing.T) {
	w := get(testRouter, "/api/v1/scores/service")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
}

func TestGetCurrentChecks_OK(t *testing.T) {
	w := get(testRouter, "/api/v1/scores/current")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
}

func TestGetDetailedScores_OK(t *testing.T) {
	w := get(testRouter, "/api/v1/scores/detailed")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
}

// seedRecent inserts rows for team 2 with explicit timestamps so the recent
// endpoint has a history to page through. Team 2 is used so the scores seeded
// for team 1 in TestMain stay untouched.
func seedRecent(t *testing.T, service types.ServiceName, count int) {
	t.Helper()
	ctx := context.Background()

	for i := range count {
		_, err := testPool.Exec(ctx, `
			INSERT INTO scores (service, box, team_num, status, stdout, stderr, timestamp)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, service, "testbox", 2, true, "out", "", 1_700_000_000_000_000+int64(i))
		if err != nil {
			t.Fatalf("seed recent score: %v", err)
		}
	}
}

func TestGetRecentChecks_OK(t *testing.T) {
	seedRecent(t, "http", 3)

	w := get(testRouter, "/api/v1/scores/recent/2")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}

	var scores []types.Score
	if err := json.Unmarshal(w.Body.Bytes(), &scores); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(scores) == 0 {
		t.Error("expected at least one recent check")
	}
}

func TestGetRecentChecks_respectsNParam(t *testing.T) {
	seedRecent(t, "http", 5)

	w := get(testRouter, "/api/v1/scores/recent/2?n=1")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}

	var scores []types.Score
	if err := json.Unmarshal(w.Body.Bytes(), &scores); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	perService := make(map[types.ServiceName]int)
	for _, s := range scores {
		perService[s.ServiceName]++
	}
	for service, n := range perService {
		if n > 1 {
			t.Errorf("service %q returned %d rows, want at most 1", service, n)
		}
	}
}

func TestGetRecentChecks_rejectsBadInput(t *testing.T) {
	cases := []struct {
		name string
		path string
	}{
		{"non-numeric team", "/api/v1/scores/recent/abc"},
		{"team zero", "/api/v1/scores/recent/0"},
		{"team above configured count", "/api/v1/scores/recent/99"},
		{"non-numeric n", "/api/v1/scores/recent/1?n=abc"},
		{"n zero", "/api/v1/scores/recent/1?n=0"},
		{"n above cap", "/api/v1/scores/recent/1?n=101"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := get(testRouter, tc.path)
			if w.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want 400; body: %s", w.Code, w.Body.String())
			}
		})
	}
}
