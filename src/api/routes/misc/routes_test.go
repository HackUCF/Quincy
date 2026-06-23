package misc_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
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

	testRouter = testutil.NewTestRouter(pool, testCfg)
	os.Exit(m.Run())
}

func get(path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, path, nil)
	testRouter.ServeHTTP(w, req)
	return w
}

func TestGetConfig_OK(t *testing.T) {
	w := get("/api/v1/config")
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if _, ok := body["num_teams"]; !ok {
		t.Error("response missing num_teams field")
	}
}

func TestNoRoute_Returns404(t *testing.T) {
	w := get("/api/v1/nonexistent")
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestNoRoute_HasSimilarRoutes(t *testing.T) {
	// /api/v1/scores is a prefix of registered routes like /api/v1/scores/team
	w := get("/api/v1/scores")
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	routes, ok := body["similar_routes"]
	if !ok {
		t.Fatal("response missing similar_routes field")
	}
	arr, ok := routes.([]any)
	if !ok || len(arr) == 0 {
		t.Error("expected non-empty similar_routes")
	}
}

func TestNoRoute_HasMessageField(t *testing.T) {
	w := get("/totally-unknown-path")
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if _, ok := body["message"]; !ok {
		t.Error("response missing message field")
	}
}
