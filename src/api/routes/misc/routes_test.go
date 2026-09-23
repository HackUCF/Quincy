package misc_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/HackUCF/Quincy/src/api/config"
	"github.com/HackUCF/Quincy/src/api/services"
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

	if err := services.InitServices(testCfg); err != nil {
		panic(err)
	}

	// the pause routes read this table on every request
	if err := testutil.SeedPauses(ctx, pool, testCfg); err != nil {
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

func post(path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, path, nil)
	testRouter.ServeHTTP(w, req)
	return w
}

// pauseState reads the pause status endpoint and returns the reported flag.
func pauseState(t *testing.T) bool {
	t.Helper()

	w := get("/api/v1/pause-status")
	if w.Code != http.StatusOK {
		t.Fatalf("pause-status status = %d, want 200; body: %s", w.Code, w.Body.String())
	}

	var body struct {
		IsPaused bool   `json:"is_paused"`
		Since    string `json:"since"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("pause-status response is not JSON: %v", err)
	}
	if body.Since == "" {
		t.Error("pause-status response missing since field")
	}
	return body.IsPaused
}

// restoreRunning leaves the competition unpaused for later tests regardless of
// what this one did to it.
func restoreRunning(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		if w := post("/api/v1/unpause"); w.Code != http.StatusOK && w.Code != http.StatusTeapot {
			t.Fatalf("restoring the running state: status = %d; body: %s", w.Code, w.Body.String())
		}
	})
}

func TestPauseStatus_reportsRunningBeforeAnyPause(t *testing.T) {
	if pauseState(t) {
		t.Error("is_paused = true before anything paused the competition")
	}
}

func TestPause_thenUnpause(t *testing.T) {
	restoreRunning(t)

	if w := post("/api/v1/pause"); w.Code != http.StatusOK {
		t.Fatalf("pause status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	if !pauseState(t) {
		t.Error("is_paused = false after a successful pause")
	}

	if w := post("/api/v1/unpause"); w.Code != http.StatusOK {
		t.Fatalf("unpause status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	if pauseState(t) {
		t.Error("is_paused = true after a successful unpause")
	}
}

func TestPause_redundantCallsAreTeapots(t *testing.T) {
	restoreRunning(t)

	// unpausing while already running changes nothing
	if w := post("/api/v1/unpause"); w.Code != http.StatusTeapot {
		t.Errorf("unpause while running = %d, want 418; body: %s", w.Code, w.Body.String())
	}

	if w := post("/api/v1/pause"); w.Code != http.StatusOK {
		t.Fatalf("pause status = %d, want 200; body: %s", w.Code, w.Body.String())
	}

	// and pausing while already paused does too
	if w := post("/api/v1/pause"); w.Code != http.StatusTeapot {
		t.Errorf("pause while paused = %d, want 418; body: %s", w.Code, w.Body.String())
	}
	if !pauseState(t) {
		t.Error("is_paused = false after a redundant pause — the state must be left alone")
	}
}

func TestPauseStatus_sinceMovesWithTheTransition(t *testing.T) {
	restoreRunning(t)

	sinceOf := func() string {
		t.Helper()
		w := get("/api/v1/pause-status")
		var body struct {
			Since string `json:"since"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatalf("pause-status response is not JSON: %v", err)
		}
		return body.Since
	}

	before := sinceOf()

	if w := post("/api/v1/pause"); w.Code != http.StatusOK {
		t.Fatalf("pause status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	after := sinceOf()
	if after == before {
		t.Errorf("since = %s before and after a real transition, want the moment of the pause", after)
	}

	// a rejected transition must not move it
	if w := post("/api/v1/pause"); w.Code != http.StatusTeapot {
		t.Fatalf("pause while paused = %d, want 418", w.Code)
	}
	if got := sinceOf(); got != after {
		t.Errorf("since moved from %s to %s on a rejected pause", after, got)
	}
}
