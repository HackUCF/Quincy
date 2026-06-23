package users_test

import (
	"bytes"
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

	if err := testutil.SeedUsers(ctx, pool, testCfg); err != nil {
		panic(err)
	}
	if err := testutil.SeedScoring(ctx, pool, testCfg); err != nil {
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

func post(path string, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(w, req)
	return w
}

func TestGetAllUsers_OK(t *testing.T) {
	w := get("/api/v1/users")
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
}

func TestSubmitPCR_OK(t *testing.T) {
	body := map[string]any{
		"username":  "alice",
		"password":  "newpass",
		"user_list": "local",
		"team_num":  1,
	}
	w := post("/api/v1/users", body)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200; body: %s", w.Code, w.Body.String())
	}
}

func TestSubmitPCR_UnknownUserList(t *testing.T) {
	body := map[string]any{
		"username":  "alice",
		"password":  "x",
		"user_list": "nonexistent",
		"team_num":  1,
	}
	w := post("/api/v1/users", body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400; body: %s", w.Code, w.Body.String())
	}
}

func TestSubmitPCR_BadJSON(t *testing.T) {
	r := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	testRouter.ServeHTTP(r, req)
	if r.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", r.Code)
	}
}
