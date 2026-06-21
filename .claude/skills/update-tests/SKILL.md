---
name: update-tests
description: Add, update, or fix tests for Quincy. Use when adding new features, changing existing behavior, or when asked to write or update tests.
argument-hint: "optional: path or area to focus on, e.g. api/routes/scoring or agent"
---

You are writing or updating tests for the Quincy scoring engine. The test harness is established — follow its patterns exactly. Do not invent new conventions.

## Test Harness Layout

```
src/
  testutil/              -- shared infrastructure for DB-backed tests
    db.go                -- NewTestDB(): spins up postgres:18.4 via testcontainers-go, applies schema
    fixtures.go          -- MinimalConfig(), SetupConfig(), SeedUsers(), SeedScoring()
    helpers.go           -- NewTestRouter(): Gin engine with DB+config injected, all routes registered

  api/config/
    validation_test.go   -- unit, package config (white-box)
  api/services/
    services_test.go     -- unit, package services (white-box, uses resetForTest)
  api/db/scoring/
    scoring_test.go      -- unit, package scoring (white-box, tests round/getResult)
    scoring_integration_test.go -- DB-backed, package scoring_test
  api/db/agent/
    agent_test.go        -- DB-backed, package agent_test
  api/db/users/
    users_test.go        -- DB-backed, package users_test
  api/routes/agent/
    routes_test.go       -- DB-backed, package agent_test
  api/routes/scoring/
    routes_test.go       -- DB-backed, package scoring_test
  agent/
    checks_test.go       -- unit, package agent (white-box, uses stub scripts on PATH)
    helpers_test.go      -- unit, package agent (white-box)
```

## Rule 1: Unit vs DB-backed

**Unit** — no Docker required, runs with `go test ./...`:
- Pure logic: math helpers, validation, config parsing, type construction
- Agent check execution: stub shell scripts via PATH manipulation
- Service queue: package-level state reset via `resetForTest()`
- Anything where you can pass `nil` for the DB pool safely

**DB-backed** — needs Docker, also runs with `go test ./...` (testcontainers spins up Postgres automatically):
- Anything that executes SQL (all `api/db/` query functions)
- HTTP route handlers (via `httptest` + `NewTestRouter`)
- Anything that calls `testutil.NewTestDB`

**Deciding:** if the code under test touches `*pgxpool.Pool` in a real way, it's DB-backed. If you can stub the DB or avoid it entirely, it's unit.

No build tags. All tests — unit and DB-backed — run with a single `go test ./...` from `src/`. Testcontainers handles Postgres container lifecycle automatically. Requires Docker to be available; set `DOCKER_HOST` if using a non-standard socket.

## Rule 2: Package Naming

- **White-box** (`package foo`): use when testing unexported functions. Required for `api/config`, `api/services`, `api/db/scoring` (unit), `agent`.
- **Black-box** (`package foo_test`): use for DB-backed tests and any test that only needs exported APIs. Required for all `api/db/*/` DB-backed tests, all `api/routes/*/` tests.

Never use `package foo` for DB-backed tests — the `testutil` import would create coupling that makes package-level state harder to reason about.

## Rule 3: No Build Tags

Do not add `//go:build` tags to test files. All tests run unconditionally with `go test ./...`.

## Rule 4: TestMain Pattern (DB-backed Tests)

Every DB-backed package needs exactly one `TestMain`. Put it at the top of the package's test file (or a `main_test.go` if there are multiple test files):

```go
var testPool *pgxpool.Pool
var testCfg *config.APIConfigSpec

func TestMain(m *testing.M) {
    ctx := context.Background()

    pool, cleanup, err := testutil.NewTestDB(ctx)
    if err != nil {
        panic(err)
    }
    defer cleanup()

    testPool = pool
    testCfg = testutil.MinimalConfig()
    testutil.SetupConfig(testCfg) // sets config.TeamRange — must happen before SeedUsers/SeedScoring

    if err := testutil.SeedUsers(ctx, pool, testCfg); err != nil {
        panic(err)
    }
    if err := testutil.SeedScoring(ctx, pool, testCfg); err != nil {
        panic(err)
    }

    os.Exit(m.Run())
}
```

For **route tests** that need the agent service queue, also add:
```go
if err := services.InitServices(testCfg); err != nil {
    panic(err)
}
testRouter = testutil.NewTestRouter(pool, testCfg)
```

`testutil.SetupConfig` MUST be called before `SeedUsers` or `SeedScoring` — both use `config.TeamRange` internally.

## Rule 5: MinimalConfig Shape

`testutil.MinimalConfig()` returns: 2 teams, 1 box (`testbox`), 2 services (`http` (no creds), `ssh` (userlist `local`)), 1 user list (`local`) with users `alice`/`bob`.

If your test needs a different config shape, construct `*config.APIConfigSpec` directly in the test — do not modify `MinimalConfig()`. Always call `testutil.SetupConfig(cfg)` after.

## Rule 6: Unit Tests for Agent Checks (Stub Scripts)

`runCheck` and `dumpService` are unexported. Test them from `package agent`.

Stub scripts go in `t.TempDir()`, prepended to PATH with `t.Setenv` (auto-restored):

```go
func addStubToPath(t *testing.T, name, content string) {
    t.Helper()
    dir := t.TempDir()
    path := filepath.Join(dir, name)
    if err := os.WriteFile(path, []byte(content), 0755); err != nil {
        t.Fatalf("write stub %q: %v", name, err)
    }
    t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}
```

Standard stubs:
- `"#!/bin/sh\nexit 0\n"` → pass
- `"#!/bin/sh\nexit 1\n"` → fail
- `"#!/bin/sh\nsleep 10\n"` → timeout (use 10ms timeout in test)
- `"#!/bin/sh\ncat \"$1\" > /dev/null && exit 0\n"` → verifies JSON arg was passed

## Rule 7: Service Queue Reset

`api/services` uses package-level vars. Between tests that call `InitServices`, call `resetForTest()` (unexported, accessible from `package services` white-box tests):

```go
func TestSomething(t *testing.T) {
    t.Cleanup(resetForTest)
    // ...
    InitServices(cfg)
}
```

Also set `config.TeamRange` directly — never call `config.SetConfig` in tests (it panics on second call):

```go
config.TeamRange = []types.TeamNum{1, 2}
```

Or use `testutil.SetupConfig(cfg)` which does the same thing.

## Rule 8: Route Tests

Use `testutil.NewTestRouter(pool, cfg)` — returns a Gin engine with custom middleware that injects pool and config directly into context. Do NOT use `conn.DBMiddleware()` or `config.ConfigMiddleware()` (both panic without the global state).

HTTP assertions:
```go
func get(router *gin.Engine, path string) *httptest.ResponseRecorder {
    w := httptest.NewRecorder()
    req, _ := http.NewRequest(http.MethodGet, path, nil)
    router.ServeHTTP(w, req)
    return w
}

// For POST with JSON body:
body, _ := json.Marshal(score)
req, _ := http.NewRequest(http.MethodPost, "/api/v1/agent/completed-score", bytes.NewReader(body))
req.Header.Set("Content-Type", "application/json")
```

## Rule 9: DB Query Tests

Call DB functions directly with `testPool`. Verify by querying the DB again or by checking the returned value:

```go
func TestAddScore_insertsRow(t *testing.T) {
    ctx := context.Background()
    score := types.Score{ServiceName: "http", BoxName: "testbox", TeamNum: 1, Status: true, Message: "ok"}
    if err := dbagent.AddScore(ctx, testPool, score); err != nil {
        t.Fatalf("AddScore: %v", err)
    }
    var count int
    testPool.QueryRow(ctx, `SELECT COUNT(*) FROM scores WHERE service = $1`, "http").Scan(&count)
    if count == 0 {
        t.Error("no row in scores")
    }
}
```

Seed data lives in the `final_scores` and `scoring_users` tables after `TestMain`. The `scores` and `recent_scores` tables start empty — seed them per-test if needed via `dbagent.AddScore`.

## What to Read Before Writing Tests

1. **The function/handler under test** — actual signatures, return types, error paths
2. **`src/testutil/`** — what helpers already exist; don't duplicate
3. **`src/api/db/schema.sql`** — table shapes for direct SQL assertions
4. **Existing test files in the same package** — match their style and `TestMain` if one exists

## Running Tests

From `src/`:
```bash
go test ./...                          # all tests
go test -count=1 -race ./...           # with race detector (recommended)
go test -run TestFoo ./api/db/agent/   # single package, single test
```

DB-backed tests require Docker. If using a non-standard socket (e.g. Colima, Podman), set `DOCKER_HOST` before running:
```bash
DOCKER_HOST=unix:///path/to/docker.sock go test ./...
```

## When Tests Are in a Package That Already Has a TestMain

If you're adding tests to a DB-backed package that already has a `TestMain`, just add test functions to the existing file (or a new `_test.go` file in the same package). Do NOT add another `TestMain` — there can only be one per package.

## When to Update testutil/

Only update `testutil/` if:
- A new table is added to the schema and needs seeding
- A new shared router setup pattern is needed for a new route group
- `MinimalConfig` is structurally wrong for an entire category of tests

Do NOT modify `testutil/` just to handle one test's special case — construct what you need inline in that test's `TestMain`.

$ARGUMENTS
