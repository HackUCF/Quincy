package graphs_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/HackUCF/quincy/api/config"
	"github.com/HackUCF/quincy/testutil"
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

	// Seed scores with distinct timestamps — required for GetScores (GetCompDuration + bucket math).
	_, err = pool.Exec(ctx, `
		INSERT INTO scores (service, box, team_num, status, message, timestamp)
		VALUES ('http', 'testbox', 1, true, 'ok', 1000000000),
		       ('http', 'testbox', 1, true, 'ok', 2000000000)
	`)
	if err != nil {
		panic(err)
	}

	// Seed recent_scores for scoreboard.
	_, err = pool.Exec(ctx, `
		INSERT INTO recent_scores (service, box, team_num, status, message, timestamp)
		VALUES ('http', 'testbox', 1, true, 'pass', 1000000000),
		       ('ssh',  'testbox', 1, false, 'fail', 2000000000)
	`)
	if err != nil {
		panic(err)
	}

	testRouter = testutil.NewTestRouter(pool, testCfg)
	os.Exit(m.Run())
}

func get(path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, path, nil)
	testRouter.ServeHTTP(w, req)
	return w
}

func TestGetScoreboard_OK(t *testing.T) {
	w := get("/api/v1/graphs/scoreboard")
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q, want text/html", ct)
	}
}

func TestGetStandings_OK(t *testing.T) {
	w := get("/api/v1/graphs/standings")
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q, want text/html", ct)
	}
}

func TestGetHeatmap_OK(t *testing.T) {
	w := get("/api/v1/graphs/heatmap")
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q, want text/html", ct)
	}
}

func TestGetScores_OK(t *testing.T) {
	w := get("/api/v1/graphs/scores")
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	if ct := w.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q, want text/html", ct)
	}
}
