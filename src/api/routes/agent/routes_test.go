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
	"github.com/HackUCF/Quincy/src/api/routes"
	"github.com/HackUCF/Quincy/src/api/services"
	"github.com/HackUCF/Quincy/src/api/sinks/postgres/pauses"
	"github.com/HackUCF/Quincy/src/common/middleware"
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

// pauseChecksFor halts check dispatch for the duration of the test and restores
// it afterwards, so a pause cannot leak into a later test in this package.
func pauseChecksFor(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	if _, err := pauses.ChangePauseState(ctx, testPool, types.Paused, types.CheckPause); err != nil {
		t.Fatalf("pause checks: %v", err)
	}
	t.Cleanup(func() {
		if _, err := pauses.ChangePauseState(ctx, testPool, types.Unpaused, types.CheckPause); err != nil {
			t.Fatalf("unpause checks: %v", err)
		}
	})
}

func getCheck(t *testing.T, router *gin.Engine) (*httptest.ResponseRecorder, types.Service) {
	t.Helper()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/agent/new-check", nil)
	router.ServeHTTP(w, req)

	var svc types.Service
	if w.Code == http.StatusOK {
		if err := json.Unmarshal(w.Body.Bytes(), &svc); err != nil {
			t.Fatalf("unmarshal response: %v", err)
		}
	}
	return w, svc
}

func TestGetCheck_checkPauseReturnsNoOp(t *testing.T) {
	pauseChecksFor(t)

	w, svc := getCheck(t, testRouter)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	if !svc.NoOp {
		t.Errorf("no_op = false while checks are paused; body: %s", w.Body.String())
	}
	if svc.Name != "" {
		t.Errorf("no-op assignment carries check details (%q); an agent must get no work", svc.Name)
	}
}

// A no-op must not consume a queue slot, or a pause would silently skip checks.
// MinimalConfig yields a four-entry rotation (2 teams x 2 services on one box),
// so the fifth real check must repeat the first — with no-ops interleaved or not.
func TestGetCheck_checkPauseDoesNotAdvanceTheQueue(t *testing.T) {
	ctx := context.Background()

	ident := func(s types.Service) string {
		return string(s.BoxName) + "/" + string(s.Name) + "/" + string(rune('0'+s.TeamNum))
	}

	_, first := getCheck(t, testRouter)

	// paused and resumed inline: the queue position either side of the pause is
	// the whole assertion, so it cannot wait on a deferred cleanup
	if _, err := pauses.ChangePauseState(ctx, testPool, types.Paused, types.CheckPause); err != nil {
		t.Fatalf("pause checks: %v", err)
	}
	for i := range 2 {
		if _, svc := getCheck(t, testRouter); !svc.NoOp {
			t.Fatalf("call %d while paused returned real work, want a no-op", i+1)
		}
	}
	if _, err := pauses.ChangePauseState(ctx, testPool, types.Unpaused, types.CheckPause); err != nil {
		t.Fatalf("unpause checks: %v", err)
	}

	// three more real checks complete the rotation
	for range 3 {
		getCheck(t, testRouter)
	}

	_, fifth := getCheck(t, testRouter)
	if ident(fifth) != ident(first) {
		t.Errorf("rotation restarted at %s, want %s; the paused calls consumed queue slots",
			ident(fifth), ident(first))
	}
}

// A scoring pause is applied when a result is recorded, not when work is handed
// out, so agents must keep getting real checks through one.
func TestGetCheck_scoringPauseStillServesChecks(t *testing.T) {
	ctx := context.Background()

	if _, err := pauses.ChangePauseState(ctx, testPool, types.Paused, types.ScoringPause); err != nil {
		t.Fatalf("pause scoring: %v", err)
	}
	t.Cleanup(func() {
		if _, err := pauses.ChangePauseState(ctx, testPool, types.Unpaused, types.ScoringPause); err != nil {
			t.Fatalf("unpause scoring: %v", err)
		}
	})

	w, svc := getCheck(t, testRouter)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	if svc.NoOp {
		t.Error("a scoring pause suppressed check dispatch")
	}
	if svc.Name == "" {
		t.Error("returned service has empty Name")
	}
}

// The check-pause lookup must be skipped entirely when no database is
// configured, or the handler dereferences a nil pool and 500s.
func TestGetCheck_servedWithoutTheDatabase(t *testing.T) {
	// a config with no sinks, which is what a database-less deployment loads
	noSinkCfg := *testCfg
	noSinkCfg.Sinks = config.Sinks{}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(
		middleware.Recovery(false),
		func(c *gin.Context) {
			// no db middleware, mirroring a deployment with no sinks
			c.Set(config.CfgKey, &noSinkCfg)
			c.Next()
		},
	)
	routes.RegisterRoutes(r, config.Sinks{})

	w, svc := getCheck(t, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
	if svc.Name == "" {
		t.Errorf("returned service has empty Name; body: %s", w.Body.String())
	}
}
