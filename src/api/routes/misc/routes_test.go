package misc_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

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

	// seeds the service queue and the competition pause state, which the
	// pause routes read on every request
	if err := services.InitServices(testCfg); err != nil {
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

// The pause state is a package-level global in services, so this walks the whole
// cycle in one test rather than splitting it across order-dependent cases.
func TestPauseUnpause_cycle(t *testing.T) {
	// leave the competition running no matter how this test exits, or every
	// later test that pulls a check gets a no-op
	t.Cleanup(func() { services.Unpause() })

	if services.IsPaused() {
		t.Fatal("competition is already paused before the test started")
	}

	if w := post("/api/v1/pause"); w.Code != http.StatusOK {
		t.Fatalf("first pause: status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	if !services.IsPaused() {
		t.Error("competition is not paused after a successful pause")
	}

	w := post("/api/v1/pause")
	if w.Code != http.StatusTeapot {
		t.Errorf("second pause: status = %d, want 418; body: %s", w.Code, w.Body.String())
	}
	if !services.IsPaused() {
		t.Error("a rejected pause changed the state")
	}

	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("418 response is not JSON: %v", err)
	}
	if _, ok := body["message"]; !ok {
		t.Error("418 response missing message field")
	}

	if w := post("/api/v1/unpause"); w.Code != http.StatusOK {
		t.Fatalf("first unpause: status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	if services.IsPaused() {
		t.Error("competition is still paused after a successful unpause")
	}

	if w := post("/api/v1/unpause"); w.Code != http.StatusTeapot {
		t.Errorf("second unpause: status = %d, want 418; body: %s", w.Code, w.Body.String())
	}
	if services.IsPaused() {
		t.Error("a rejected unpause changed the state")
	}
}

func TestPause_agentReceivesNoOp(t *testing.T) {
	t.Cleanup(func() { services.Unpause() })

	if w := post("/api/v1/pause"); w.Code != http.StatusOK {
		t.Fatalf("pause: status = %d, want 200", w.Code)
	}

	w := get("/api/v1/agent/new-check")
	if w.Code != http.StatusOK {
		t.Fatalf("new-check: status = %d, want 200; body: %s", w.Code, w.Body.String())
	}

	var svc map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &svc); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if noOp, _ := svc["no_op"].(bool); !noOp {
		t.Errorf("agent got a real check while paused: %s", w.Body.String())
	}
}

func TestPauseStatus_reflectsState(t *testing.T) {
	t.Cleanup(func() { services.Unpause() })

	check := func(wantPaused bool) {
		t.Helper()

		w := get("/api/v1/pause-status")
		if w.Code != http.StatusOK {
			t.Fatalf("pause-status: status = %d, want 200; body: %s", w.Code, w.Body.String())
		}

		var body map[string]any
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatalf("response is not JSON: %v", err)
		}

		got, ok := body["is_paused"].(bool)
		if !ok {
			t.Fatalf("response missing is_paused field: %s", w.Body.String())
		}
		if got != wantPaused {
			t.Errorf("is_paused = %v, want %v", got, wantPaused)
		}

		since, ok := body["since"].(string)
		if !ok || since == "" {
			t.Errorf("response missing since field: %s", w.Body.String())
		}
		if _, err := time.Parse(time.RFC3339, since); err != nil {
			t.Errorf("since = %q is not an RFC 3339 timestamp: %v", since, err)
		}
	}

	check(false)

	if w := post("/api/v1/pause"); w.Code != http.StatusOK {
		t.Fatalf("pause: status = %d, want 200", w.Code)
	}
	check(true)

	if w := post("/api/v1/unpause"); w.Code != http.StatusOK {
		t.Fatalf("unpause: status = %d, want 200", w.Code)
	}
	check(false)
}
