package competition_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/HackUCF/Quincy/src/api/config"
	"github.com/HackUCF/Quincy/src/api/routes"
	"github.com/HackUCF/Quincy/src/api/services"
	"github.com/HackUCF/Quincy/src/common/middleware"
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

	// every handler in this package reads this table on every request
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

func post(path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, path, nil)
	testRouter.ServeHTTP(w, req)
	return w
}

type statusBody struct {
	ScoringPaused bool      `json:"scoring_is_paused"`
	ScoringSince  time.Time `json:"scoring_since"`
	ChecksPaused  bool      `json:"checks_are_paused"`
	ChecksSince   time.Time `json:"checks_since"`
}

// pauseStatus reads the status endpoint and returns both reported states.
func pauseStatus(t *testing.T) statusBody {
	t.Helper()

	w := get("/api/v1/comp/pause-status")
	if w.Code != http.StatusOK {
		t.Fatalf("pause-status status = %d, want 200; body: %s", w.Code, w.Body.String())
	}

	var body statusBody
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("pause-status response is not JSON: %v", err)
	}
	if body.ScoringSince.IsZero() || body.ChecksSince.IsZero() {
		t.Errorf("pause-status missing a since field: %s", w.Body.String())
	}
	return body
}

// restoreRunning leaves both pauses off for later tests regardless of what this
// one did to them.
func restoreRunning(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		for _, path := range []string{"/api/v1/comp/unpause-scoring", "/api/v1/comp/unpause-checks"} {
			if w := post(path); w.Code != http.StatusOK && w.Code != http.StatusTeapot {
				t.Fatalf("restoring the running state via %s: status = %d; body: %s", path, w.Code, w.Body.String())
			}
		}
	})
}

func TestPauseStatus_reportsRunningBeforeAnyPause(t *testing.T) {
	got := pauseStatus(t)
	if got.ScoringPaused || got.ChecksPaused {
		t.Errorf("fresh competition reports scoring=%v checks=%v, want both running", got.ScoringPaused, got.ChecksPaused)
	}
}

func TestPauseScoring_thenUnpause(t *testing.T) {
	restoreRunning(t)

	if w := post("/api/v1/comp/pause-scoring"); w.Code != http.StatusOK {
		t.Fatalf("pause-scoring status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	if got := pauseStatus(t); !got.ScoringPaused {
		t.Error("scoring reports running after pause-scoring")
	}

	if w := post("/api/v1/comp/unpause-scoring"); w.Code != http.StatusOK {
		t.Fatalf("unpause-scoring status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	if got := pauseStatus(t); got.ScoringPaused {
		t.Error("scoring reports paused after unpause-scoring")
	}
}

func TestPauseChecks_thenUnpause(t *testing.T) {
	restoreRunning(t)

	if w := post("/api/v1/comp/pause-checks"); w.Code != http.StatusOK {
		t.Fatalf("pause-checks status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	if got := pauseStatus(t); !got.ChecksPaused {
		t.Error("checks report running after pause-checks")
	}

	if w := post("/api/v1/comp/unpause-checks"); w.Code != http.StatusOK {
		t.Fatalf("unpause-checks status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	if got := pauseStatus(t); got.ChecksPaused {
		t.Error("checks report paused after unpause-checks")
	}
}

// The two pauses share a table and a handler, so the obvious failure is one
// route moving the other's state.
func TestPauses_areIndependent(t *testing.T) {
	restoreRunning(t)

	if w := post("/api/v1/comp/pause-checks"); w.Code != http.StatusOK {
		t.Fatalf("pause-checks status = %d, want 200; body: %s", w.Code, w.Body.String())
	}

	got := pauseStatus(t)
	if !got.ChecksPaused {
		t.Error("checks report running after pause-checks")
	}
	if got.ScoringPaused {
		t.Error("pausing checks also paused scoring")
	}
}

func TestPause_redundantCallsAreTeapots(t *testing.T) {
	restoreRunning(t)

	// already running
	if w := post("/api/v1/comp/unpause-scoring"); w.Code != http.StatusTeapot {
		t.Errorf("unpause while running = %d, want 418; body: %s", w.Code, w.Body.String())
	}

	if w := post("/api/v1/comp/pause-scoring"); w.Code != http.StatusOK {
		t.Fatalf("pause-scoring status = %d, want 200; body: %s", w.Code, w.Body.String())
	}

	// already paused
	if w := post("/api/v1/comp/pause-scoring"); w.Code != http.StatusTeapot {
		t.Errorf("pause while paused = %d, want 418; body: %s", w.Code, w.Body.String())
	}

	if got := pauseStatus(t); !got.ScoringPaused {
		t.Error("a rejected redundant pause changed the state")
	}
}

func TestPauseStatus_sinceMovesWithTheTransition(t *testing.T) {
	restoreRunning(t)

	before := pauseStatus(t)

	time.Sleep(2 * time.Millisecond)

	if w := post("/api/v1/comp/pause-scoring"); w.Code != http.StatusOK {
		t.Fatalf("pause-scoring status = %d, want 200; body: %s", w.Code, w.Body.String())
	}

	after := pauseStatus(t)
	if !after.ScoringSince.After(before.ScoringSince) {
		t.Errorf("scoring_since did not advance: %v then %v", before.ScoringSince, after.ScoringSince)
	}
	if !after.ChecksSince.Equal(before.ChecksSince) {
		t.Errorf("checks_since moved on a scoring-only transition: %v then %v", before.ChecksSince, after.ChecksSince)
	}
}

// noSinkRouter builds the real route tree with no sinks configured, which is how
// a deployment without the database runs.
func noSinkRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(
		middleware.Recovery(false),
		func(c *gin.Context) {
			// the config middleware is applied whether or not a sink exists,
			// so mirror it here; only the db middleware is sink-dependent
			c.Set(config.CfgKey, testCfg)
			c.Next()
		},
	)
	routes.RegisterRoutes(r, config.Sinks{})
	return r
}

func TestPauseRoutes_return501WithoutTheDatabase(t *testing.T) {
	router := noSinkRouter()

	for _, tc := range []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/api/v1/comp/pause-scoring"},
		{http.MethodPost, "/api/v1/comp/unpause-scoring"},
		{http.MethodPost, "/api/v1/comp/pause-checks"},
		{http.MethodPost, "/api/v1/comp/unpause-checks"},
		{http.MethodGet, "/api/v1/comp/pause-status"},
	} {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.path, nil)
			router.ServeHTTP(w, req)

			// these handlers read the pool with no nil guard, so if the gate
			// is ever dropped this returns a recovered 500 instead
			if w.Code != http.StatusNotImplemented {
				t.Errorf("status = %d, want 501; body: %s", w.Code, w.Body.String())
			}
		})
	}
}
