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

// restoreRunning puts the shared database back in the running state for both
// pause kinds, so one test's pause cannot leak into the next.
func restoreRunning(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		ctx := context.Background()
		for _, pt := range []types.PauseType{types.ScoringPause, types.CheckPause} {
			if _, err := pauses.ChangePauseState(ctx, testPool, types.Unpaused, pt); err != nil {
				t.Fatalf("restore running state for %q: %v", pt, err)
			}
		}
	})
}

func rowCount(t *testing.T, pauseType types.PauseType) int {
	t.Helper()
	var n int
	err := testPool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM pauses WHERE type = $1`, pauseType,
	).Scan(&n)
	if err != nil {
		t.Fatalf("count pauses for %q: %v", pauseType, err)
	}
	return n
}

// InitPauses writes one seed row per pause kind in a single transaction. The
// timestamp column is uniquely constrained, so this is also the regression test
// for both seeds sharing a timestamp and colliding.
func TestInitPauses_seedsBothKinds(t *testing.T) {
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
				testutil.SkipDBTests(err)
			}
			defer cleanup()

			cfg := testutil.MinimalConfig()
			cfg.StartPaused = tc.start

			if err := pauses.InitPauses(ctx, pool, cfg); err != nil {
				t.Fatalf("InitPauses: %v", err)
			}

			for _, pt := range []types.PauseType{types.ScoringPause, types.CheckPause} {
				got, err := pauses.IsPaused(ctx, pool, pt)
				if err != nil {
					t.Fatalf("IsPaused(%q): %v", pt, err)
				}
				if got != tc.start {
					t.Errorf("%q state = %v, want %v", pt, got, tc.start)
				}
			}
		})
	}
}

func TestInitPauses_isNoOpAfterFirstBoot(t *testing.T) {
	ctx := context.Background()

	before := rowCount(t, types.ScoringPause) + rowCount(t, types.CheckPause)

	// TestMain already seeded this database; a second call must write nothing
	if err := pauses.InitPauses(ctx, testPool, testCfg); err != nil {
		t.Fatalf("InitPauses: %v", err)
	}

	if after := rowCount(t, types.ScoringPause) + rowCount(t, types.CheckPause); after != before {
		t.Errorf("row count = %d, want %d; re-initializing seeded a duplicate", after, before)
	}
}

func TestChangePauseState_togglesAndReportsChange(t *testing.T) {
	restoreRunning(t)
	ctx := context.Background()

	changed, err := pauses.ChangePauseState(ctx, testPool, types.Paused, types.ScoringPause)
	if err != nil {
		t.Fatalf("ChangePauseState: %v", err)
	}
	if !changed {
		t.Error("changed = false, want true on a real transition")
	}

	got, err := pauses.IsPaused(ctx, testPool, types.ScoringPause)
	if err != nil {
		t.Fatalf("IsPaused: %v", err)
	}
	if got != types.Paused {
		t.Errorf("state = %v, want paused", got)
	}
}

func TestChangePauseState_redundantCallWritesNothing(t *testing.T) {
	restoreRunning(t)
	ctx := context.Background()

	if _, err := pauses.ChangePauseState(ctx, testPool, types.Paused, types.CheckPause); err != nil {
		t.Fatalf("first pause: %v", err)
	}
	before := rowCount(t, types.CheckPause)

	changed, err := pauses.ChangePauseState(ctx, testPool, types.Paused, types.CheckPause)
	if err != nil {
		t.Fatalf("second pause: %v", err)
	}
	if changed {
		t.Error("changed = true, want false when already in the desired state")
	}
	if after := rowCount(t, types.CheckPause); after != before {
		t.Errorf("row count = %d, want %d; a redundant call appended an entry", after, before)
	}
}

// The two pause kinds share a table and must not read or clobber each other.
func TestChangePauseState_kindsAreIndependent(t *testing.T) {
	restoreRunning(t)
	ctx := context.Background()

	if _, err := pauses.ChangePauseState(ctx, testPool, types.Paused, types.CheckPause); err != nil {
		t.Fatalf("pause checks: %v", err)
	}

	scoring, err := pauses.IsPaused(ctx, testPool, types.ScoringPause)
	if err != nil {
		t.Fatalf("IsPaused(scoring): %v", err)
	}
	if scoring != types.Unpaused {
		t.Error("pausing checks also paused scoring")
	}
}

func TestIsPausedTx_readsInsideATransaction(t *testing.T) {
	restoreRunning(t)
	ctx := context.Background()

	if _, err := pauses.ChangePauseState(ctx, testPool, types.Paused, types.ScoringPause); err != nil {
		t.Fatalf("pause scoring: %v", err)
	}

	tx, err := testPool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	defer tx.Rollback(ctx)

	got, err := pauses.IsPausedTx(ctx, tx, types.ScoringPause)
	if err != nil {
		t.Fatalf("IsPausedTx: %v", err)
	}
	if got != types.Paused {
		t.Errorf("state = %v, want paused", got)
	}
}

// GetPauseRecord reads both kinds in one transaction; the two reads are easy to
// cross-wire, so assert each field against a state only that kind is in.
func TestGetPauseRecord_reportsEachKindSeparately(t *testing.T) {
	restoreRunning(t)
	ctx := context.Background()

	if _, err := pauses.ChangePauseState(ctx, testPool, types.Paused, types.CheckPause); err != nil {
		t.Fatalf("pause checks: %v", err)
	}

	record, err := pauses.GetPauseRecord(ctx, testPool)
	if err != nil {
		t.Fatalf("GetPauseRecord: %v", err)
	}

	if record.CheckState != types.Paused {
		t.Error("CheckState = running, want paused")
	}
	if record.ScoringState != types.Unpaused {
		t.Error("ScoringState = paused, want running; the scoring read may be scanning the check row")
	}
	if record.CheckSince.IsZero() || record.ScoringSince.IsZero() {
		t.Errorf("timestamps not populated: scoring %v, checks %v", record.ScoringSince, record.CheckSince)
	}
}

func TestGetPauseRecord_sinceMovesWithTheTransition(t *testing.T) {
	restoreRunning(t)
	ctx := context.Background()

	first, err := pauses.GetPauseRecord(ctx, testPool)
	if err != nil {
		t.Fatalf("GetPauseRecord: %v", err)
	}

	time.Sleep(2 * time.Millisecond)

	if _, err := pauses.ChangePauseState(ctx, testPool, types.Paused, types.ScoringPause); err != nil {
		t.Fatalf("pause scoring: %v", err)
	}

	second, err := pauses.GetPauseRecord(ctx, testPool)
	if err != nil {
		t.Fatalf("GetPauseRecord: %v", err)
	}

	if !second.ScoringSince.After(first.ScoringSince) {
		t.Errorf("ScoringSince did not advance: %v then %v", first.ScoringSince, second.ScoringSince)
	}
	if !second.CheckSince.Equal(first.CheckSince) {
		t.Errorf("CheckSince moved on a scoring-only transition: %v then %v", first.CheckSince, second.CheckSince)
	}
}
