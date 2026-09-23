package pauses_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/HackUCF/Quincy/src/api/config"
	"github.com/HackUCF/Quincy/src/api/sinks/postgres/pauses"
	"github.com/HackUCF/Quincy/src/common/types"
	"github.com/HackUCF/Quincy/src/testutil"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testPool *pgxpool.Pool
var testCfg *config.APIConfigSpec

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

	if err := testutil.SeedPauses(ctx, pool, testCfg); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

// restoreRunning puts the shared database back in the running state so one
// test's pause cannot leak into the next.
func restoreRunning(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		if _, err := pauses.ChangePauseState(context.Background(), testPool, types.Unpaused); err != nil {
			t.Fatalf("restore running state: %v", err)
		}
	})
}

// currentState reads the pause state through a throwaway transaction, which is
// the only way IsPaused can be called.
func currentState(t *testing.T) types.PauseState {
	t.Helper()
	ctx := context.Background()

	tx, err := testPool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	defer tx.Rollback(ctx)

	state, err := pauses.IsPaused(ctx, tx)
	if err != nil {
		t.Fatalf("IsPaused: %v", err)
	}
	return state
}

func rowCount(t *testing.T) int {
	t.Helper()
	var n int
	if err := testPool.QueryRow(context.Background(), `SELECT COUNT(*) FROM pause_states`).Scan(&n); err != nil {
		t.Fatalf("count pause_states: %v", err)
	}
	return n
}

func TestInitPauses_seedsConfiguredState(t *testing.T) {
	for _, tc := range []struct {
		name  string
		start types.PauseState
	}{
		{"starts running", types.Unpaused},
		{"starts paused", types.Paused},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			// a fresh database, because seeding only ever happens once per database
			pool, cleanup, err := testutil.NewTestDB(ctx)
			if err != nil {
				t.Fatalf("NewTestDB: %v", err)
			}
			defer cleanup()

			cfg := testutil.MinimalConfig()
			cfg.StartPaused = tc.start
			testutil.SetupConfig(cfg)

			if err := pauses.InitPauses(ctx, pool, cfg); err != nil {
				t.Fatalf("InitPauses: %v", err)
			}

			record, err := pauses.GetPauseRecord(ctx, pool)
			if err != nil {
				t.Fatalf("GetPauseRecord: %v", err)
			}
			if record.State != tc.start {
				t.Errorf("seeded state = %v, want %v", record.State, tc.start)
			}
		})
	}
}

func TestInitPauses_secondRunIsNoOp(t *testing.T) {
	ctx := context.Background()

	pool, cleanup, err := testutil.NewTestDB(ctx)
	if err != nil {
		t.Fatalf("NewTestDB: %v", err)
	}
	defer cleanup()

	cfg := testutil.MinimalConfig()
	cfg.StartPaused = types.Unpaused
	testutil.SetupConfig(cfg)

	if err := pauses.InitPauses(ctx, pool, cfg); err != nil {
		t.Fatalf("InitPauses (first): %v", err)
	}
	if _, err := pauses.ChangePauseState(ctx, pool, types.Paused); err != nil {
		t.Fatalf("ChangePauseState: %v", err)
	}

	// a restart must not overwrite the state an operator actually left behind
	cfg.StartPaused = types.Unpaused
	if err := pauses.InitPauses(ctx, pool, cfg); err != nil {
		t.Fatalf("InitPauses (second): %v", err)
	}

	record, err := pauses.GetPauseRecord(ctx, pool)
	if err != nil {
		t.Fatalf("GetPauseRecord: %v", err)
	}
	if record.State != types.Paused {
		t.Errorf("state = %v after re-init, want the stored %v", record.State, types.Paused)
	}
}

func TestChangePauseState_reportsOnlyRealTransitions(t *testing.T) {
	ctx := context.Background()
	restoreRunning(t)

	changed, err := pauses.ChangePauseState(ctx, testPool, types.Paused)
	if err != nil {
		t.Fatalf("ChangePauseState: %v", err)
	}
	if !changed {
		t.Error("ChangePauseState(paused) = false on first pause, want true")
	}
	if got := currentState(t); got != types.Paused {
		t.Errorf("state = %v after pausing, want %v", got, types.Paused)
	}

	before := rowCount(t)
	changed, err = pauses.ChangePauseState(ctx, testPool, types.Paused)
	if err != nil {
		t.Fatalf("ChangePauseState (repeat): %v", err)
	}
	if changed {
		t.Error("ChangePauseState(paused) = true when already paused, want false")
	}
	if after := rowCount(t); after != before {
		t.Errorf("pause_states grew from %d to %d on a redundant pause, want no new row", before, after)
	}
}

func TestChangePauseState_appendsHistory(t *testing.T) {
	ctx := context.Background()
	restoreRunning(t)

	before := rowCount(t)

	for _, state := range []types.PauseState{types.Paused, types.Unpaused, types.Paused} {
		if _, err := pauses.ChangePauseState(ctx, testPool, state); err != nil {
			t.Fatalf("ChangePauseState(%v): %v", state, err)
		}
	}

	if after := rowCount(t); after != before+3 {
		t.Errorf("pause_states holds %d rows, want %d — every transition is kept", after, before+3)
	}
}

func TestGetPauseRecord_reportsStateAndMoment(t *testing.T) {
	ctx := context.Background()
	restoreRunning(t)

	start := time.Now()
	if _, err := pauses.ChangePauseState(ctx, testPool, types.Paused); err != nil {
		t.Fatalf("ChangePauseState: %v", err)
	}

	record, err := pauses.GetPauseRecord(ctx, testPool)
	if err != nil {
		t.Fatalf("GetPauseRecord: %v", err)
	}
	if record.State != types.Paused {
		t.Errorf("State = %v, want %v", record.State, types.Paused)
	}
	if record.Since.Before(start.Add(-time.Second)) || record.Since.After(time.Now().Add(time.Second)) {
		t.Errorf("Since = %v, want a timestamp from the moment of the pause (~%v)", record.Since, start)
	}
}

func TestCleanupPauses_collapsesRepeatedStates(t *testing.T) {
	ctx := context.Background()

	pool, cleanup, err := testutil.NewTestDB(ctx)
	if err != nil {
		t.Fatalf("NewTestDB: %v", err)
	}
	defer cleanup()

	// write a run of duplicates directly, since the state changer refuses to create one
	ts := time.Now().UnixMicro()
	states := []types.PauseState{types.Paused, types.Paused, types.Paused, types.Unpaused, types.Unpaused}
	for i, state := range states {
		if _, err := pool.Exec(ctx,
			`INSERT INTO pause_states (timestamp, state) VALUES ($1, $2)`, ts+int64(i), state,
		); err != nil {
			t.Fatalf("insert %d: %v", i, err)
		}
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	if err := pauses.CleanupPauses(ctx, tx); err != nil {
		t.Fatalf("CleanupPauses: %v", err)
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatalf("commit: %v", err)
	}

	rows, err := pool.Query(ctx, `SELECT state FROM pause_states ORDER BY timestamp ASC`)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	defer rows.Close()

	var got []types.PauseState
	for rows.Next() {
		var state types.PauseState
		if err := rows.Scan(&state); err != nil {
			t.Fatalf("scan: %v", err)
		}
		got = append(got, state)
	}

	want := []types.PauseState{types.Paused, types.Unpaused}
	if len(got) != len(want) {
		t.Fatalf("history = %v, want the alternating %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("history = %v, want the alternating %v", got, want)
		}
	}
}

// seedRun writes a run of identical states straight into the table, bypassing
// the state changer, which refuses to create one.
func seedRun(t *testing.T, pool *pgxpool.Pool, states ...types.PauseState) int64 {
	t.Helper()
	ctx := context.Background()

	ts := time.Now().UnixMicro()
	for i, state := range states {
		if _, err := pool.Exec(ctx,
			`INSERT INTO pause_states (timestamp, state) VALUES ($1, $2)`, ts+int64(i), state,
		); err != nil {
			t.Fatalf("insert %d: %v", i, err)
		}
	}
	return ts
}

func countRows(t *testing.T, pool *pgxpool.Pool) int {
	t.Helper()
	var n int
	if err := pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM pause_states`).Scan(&n); err != nil {
		t.Fatalf("count pause_states: %v", err)
	}
	return n
}

func TestGetPauseRecord_collapsesDuplicatesBeforeReading(t *testing.T) {
	ctx := context.Background()

	pool, cleanup, err := testutil.NewTestDB(ctx)
	if err != nil {
		t.Fatalf("NewTestDB: %v", err)
	}
	defer cleanup()

	// three pauses in a row, as two racing requests would leave behind
	first := seedRun(t, pool, types.Paused, types.Paused, types.Paused)

	record, err := pauses.GetPauseRecord(ctx, pool)
	if err != nil {
		t.Fatalf("GetPauseRecord: %v", err)
	}

	if n := countRows(t, pool); n != 1 {
		t.Errorf("pause_states holds %d rows after a read, want 1 — the read must collapse the run", n)
	}
	if record.State != types.Paused {
		t.Errorf("State = %v, want %v", record.State, types.Paused)
	}
	// the surviving entry is the first of the run, so the reported moment is
	// when the pause actually began rather than when the duplicate landed
	if got := record.Since.UnixMicro(); got != first {
		t.Errorf("Since = %d, want %d — the first entry of the run is the one that survives", got, first)
	}
}

func TestChangePauseState_collapsesDuplicatesBeforeDeciding(t *testing.T) {
	ctx := context.Background()

	pool, cleanup, err := testutil.NewTestDB(ctx)
	if err != nil {
		t.Fatalf("NewTestDB: %v", err)
	}
	defer cleanup()

	seedRun(t, pool, types.Paused, types.Paused)

	// the duplicate must not confuse the decision: the state is still paused,
	// so unpausing is a real transition and pausing again is not
	changed, err := pauses.ChangePauseState(ctx, pool, types.Paused)
	if err != nil {
		t.Fatalf("ChangePauseState(paused): %v", err)
	}
	if changed {
		t.Error("ChangePauseState(paused) = true over a duplicated pause run, want false")
	}

	// the no-change path returns before committing, so its cleanup is rolled
	// back with the rest of the transaction and the duplicate is still there
	if n := countRows(t, pool); n != 2 {
		t.Errorf("pause_states holds %d rows after a no-op transition, want the 2 seeded", n)
	}

	changed, err = pauses.ChangePauseState(ctx, pool, types.Unpaused)
	if err != nil {
		t.Fatalf("ChangePauseState(unpaused): %v", err)
	}
	if !changed {
		t.Error("ChangePauseState(unpaused) = false, want true")
	}
	// a committed transition does persist the collapse: the duplicated pause
	// becomes one row, plus the unpause just written
	if n := countRows(t, pool); n != 2 {
		t.Errorf("pause_states holds %d rows, want 2 — one collapsed pause plus one unpause", n)
	}
}
