package scoring_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/HackUCF/quincy/api/config"
	dbagent "github.com/HackUCF/quincy/api/sinks/postgres/agent"
	"github.com/HackUCF/quincy/common/types"
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

	// seed one score so current status has a row
	dbagent.AddScore(ctx, pool, types.Score{
		ServiceName: "http",
		BoxName:     "testbox",
		TeamNum:     1,
		Status:      true,
		Message:     "ok",
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
