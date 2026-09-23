package agent_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/HackUCF/Quincy/src/api/config"
	"github.com/HackUCF/Quincy/src/api/services"
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

	if err := testutil.SeedUsers(ctx, pool, testCfg); err != nil {
		panic(err)
	}
	if err := testutil.SeedScoring(ctx, pool, testCfg); err != nil {
		panic(err)
	}
	if err := testutil.SeedPauses(ctx, pool, testCfg); err != nil {
		panic(err)
	}
	if err := services.InitServices(testCfg); err != nil {
		panic(err)
	}

	testRouter = testutil.NewTestRouter(pool, testCfg)
	os.Exit(m.Run())
}

func TestGetCheck_OK(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/agent/new-check", nil)
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}

	var svc types.Service
	if err := json.Unmarshal(w.Body.Bytes(), &svc); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if svc.Name == "" {
		t.Error("returned service has empty Name")
	}
}

func TestAddScore_OK(t *testing.T) {
	score := types.Score{
		ServiceName: "http",
		BoxName:     "testbox",
		TeamNum:     1,
		Status:      true,
		Stdout:      "test ok",
		Stderr:      "test warn",
	}
	body, _ := json.Marshal(score)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/agent/completed-score",
		bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
}

func TestAddScore_BadJSON(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/agent/completed-score",
		bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}
